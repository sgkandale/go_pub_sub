package main

import (
	"context"
	"log"

	"pubsub/pkg/pubsubPB"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// create grpc client
	grpcConn, err := grpc.Dial(
		"0.0.0.0:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("[ERROR] creating grpc conn : %s", err.Error())
	}
	defer grpcConn.Close()

	pubsubClient := pubsubPB.NewPubSubServerClient(grpcConn)

	stream, err := pubsubClient.Subscribe(
		context.Background(),
		&pubsubPB.SubscribeRequest{
			Topic: "test-topic-1",
		},
	)
	if err != nil {
		log.Fatalf("[ERROR] subscribing : %s", err.Error())
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("[ERROR] receiving message : %s", err.Error())
		}
		log.Printf("[INFO] message received : %s", msg.Message)
	}
}
