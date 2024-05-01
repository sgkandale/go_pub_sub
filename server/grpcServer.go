package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"pubsub/pkg/errors"
	"pubsub/pkg/pubsub"
	"pubsub/pkg/pubsubPB"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	pubsubPB.UnimplementedPubSubServerServer

	// internals
	listener   net.Listener
	port       int
	grpcServer *grpc.Server
	pubsub     *pubsub.PubSub
}

// NewGrpc creats a new grpc server instance
func NewGrpc(ctx context.Context, req *NewServerRequest) *GrpcServer {
	if req == nil {
		log.Printf("[ERROR] nil request to create grpc server")
		return nil
	}
	if req.Cfg == nil {
		log.Printf("[ERROR] nil config to create grpc server")
		return nil
	}
	if req.PubSub == nil {
		log.Printf("[ERROR] nil pubsub value to create grpc server")
		return nil
	}

	newServer := grpc.NewServer()

	// create listener
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", req.Cfg.Port))
	if err != nil {
		log.Printf("[ERROR] creating tcp listener on port %d for grpc server : %s", req.Cfg.Port, err)
	}

	pubsubServer := &GrpcServer{
		listener:   lis,
		port:       req.Cfg.Port,
		grpcServer: newServer,
		pubsub:     req.PubSub,
	}
	pubsubPB.RegisterPubSubServerServer(newServer, pubsubServer)

	return pubsubServer
}

// Serve starts the grpc server
func (s *GrpcServer) Serve() error {
	if s.listener == nil {
		return fmt.Errorf("grpc listener not registered")
	}
	log.Printf("[INFO] starting grpc server on port %d", s.port)
	err := s.grpcServer.Serve(s.listener)
	if err != nil {
		return fmt.Errorf("starting grpc server on port %d : %s", s.port, err)
	}
	return nil
}

// Stop will attempt to stop the grpc server
func (s *GrpcServer) Stop() {
	if s.grpcServer != nil {
		log.Printf("[WARN] stopping grpc server on port %d", s.port)
		s.grpcServer.Stop()
		log.Printf("[WARN] grpc server on port %d is stopped", s.port)
	}
}

func (s *GrpcServer) Publish(ctx context.Context, in *pubsubPB.PublishRequest) (*pubsubPB.PublishResponse, error) {
	// validations
	if in == nil {
		return nil, errors.ErrInvalidRequest
	}
	if in.Message == nil {
		return nil, errors.ErrInvalidRequest
	}
	if in.Message.Topic == "" {
		return nil, errors.ErrNoTopic
	}
	if len(in.Message.Body) == 0 {
		return nil, errors.ErrNoMsgBody
	}

	msgUuid := uuid.NewString()
	in.Message.Uuid = msgUuid
	in.Message.Timestamp = timestamppb.New(time.Now().UTC())

	// push message in pubsub
	s.pubsub.Publish(in.Message.Topic, pubsub.NewMessageFromPB(in.Message))

	// return response
	return &pubsubPB.PublishResponse{
		Uuid: msgUuid,
	}, nil
}

func (s *GrpcServer) Subscribe(in *pubsubPB.SubscribeRequest, stream pubsubPB.PubSubServer_SubscribeServer) error {
	// validations
	if in == nil {
		return errors.ErrInvalidRequest
	}
	if in.Topic == "" {
		return errors.ErrNoTopic
	}

	// register subscriber in pubsub
	subscriber := s.pubsub.Subscribe(in.Topic)
	defer s.pubsub.Unsubscribe(subscriber)
	log.Printf("[INFO] new subscriber connected for topic : '%s'", in.Topic)

	go func() {
		var err error
		// read messages from subscriber channel
		for msg := range subscriber.Listen() {
			// exit if stream in cancelled
			if stream.Context().Err() != nil {
				log.Printf("[ERROR] client stream in grpc : %s", err.Error())
				return
			}
			// send message to client
			err = stream.Send(
				&pubsubPB.SubscribeResponse{
					Message: pubsub.MessageToPB(msg),
				},
			)
			if err != nil {
				log.Printf("[ERROR] sending message to subscribed client : %s", err.Error())
				if err == context.Canceled {
					return
				}
			}
		}
	}()

	<-stream.Context().Done()
	err := stream.Context().Err()
	if err != nil && err != context.Canceled {
		return err
	}
	return nil
}
