package main

import (
	"context"
	"log"
	"time"

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

	publishCtx, cancelPublishCtx := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelPublishCtx()

	resp, err := pubsubClient.Publish(
		publishCtx,
		&pubsubPB.PublishRequest{
			Message: &pubsubPB.Message{
				Body:  []byte("test-body"),
				Topic: "test-topic-1",
				Tags: map[string]string{
					"tag1": "val1",
					"tag2": "val2",
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("[ERROR] publishing message : %s", err.Error())
	}
	log.Printf("[INFO] published message : %s", resp.GetUuid())
}
