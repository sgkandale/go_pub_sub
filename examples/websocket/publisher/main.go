package main

import (
	"log"

	"pubsub/pkg/pubsub"
	"pubsub/server"

	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8082/ws", nil)
	if err != nil {
		log.Fatalf("[ERROR] dialing websocket : %s", err)
	}
	defer conn.Close()

	// send initialisation message to websocket
	err = conn.WriteJSON(
		server.WebsocketInitRequest{
			ClientType: "publisher",
			Topic:      "test-topic-1",
		},
	)
	if err != nil {
		log.Fatalf("[ERROR] writing initialisation message to websocket : %s", err)
	}

	err = conn.WriteJSON(
		pubsub.Message{
			Body:  []byte("test-boidy"),
			Topic: "test-topic-1",
			Tags: map[string]string{
				"tag1": "val1",
				"tag2": "val2",
			},
		},
	)
	if err != nil {
		log.Fatalf("[ERROR] writing pubsub message to websocket : %s", err)
	}
}
