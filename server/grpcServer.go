package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"pubsub/config"
	"pubsub/pkg/errors"
	"pubsub/pkg/pubsubPB"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	pubsubPB.UnimplementedPubSubServerServer

	// internals
	pubsubChan chan *pubsubPB.Message
	listener   net.Listener
	port       int
	grpcServer *grpc.Server
}

func NewGrpc(ctx context.Context, cfg *config.ServerConfig) *GrpcServer {
	newServer := grpc.NewServer()

	// create listener
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		log.Printf("[ERROR] creating tcp listener on port %d for grpc server : %s", cfg.Port, err)
	}

	pubsubServer := &GrpcServer{
		pubsubChan: make(chan *pubsubPB.Message, 10_000),
		listener:   lis,
		port:       cfg.Port,
		grpcServer: newServer,
	}
	pubsubPB.RegisterPubSubServerServer(newServer, pubsubServer)

	return pubsubServer
}

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
	if len(in.Message.Body) == 0 {
		return nil, errors.ErrNoMsgBody
	}

	msgUuid := uuid.NewString()
	in.Message.Uuid = msgUuid
	in.Message.Timestamp = timestamppb.New(time.Now().UTC())

	// push message in channel
	s.pubsubChan <- in.Message

	// return response
	return &pubsubPB.PublishResponse{
		Uuid: msgUuid,
	}, nil
}

func (s *GrpcServer) Subscribe(in *pubsubPB.SubscribeRequest, stream pubsubPB.PubSubServer_SubscribeServer) error {
	log.Print("subscribe")
	// validations
	if in == nil {
		return errors.ErrInvalidRequest
	}
	var err error

	// read messages from channel
	for msg := range s.pubsubChan {
		// send message to client
		err = stream.Send(
			&pubsubPB.SubscribeResponse{
				Message: msg,
			},
		)
		if err != nil {
			return status.Errorf(codes.Internal, "writing response to stream")
		}
	}
	return nil
}
