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
			ClientType: "subscriber",
			Topic:      "test-topic-1",
		},
	)
	if err != nil {
		log.Fatalf("[ERROR] writing initialisation message to websocket : %s", err)
	}

	// listen for incoming messages
	log.Print("[INFO] listening for messages")
	for {
		var msg pubsub.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Fatalf("[ERROR] reading message from websocket : %s", err)
		}
		log.Printf("[INFO] received message : %+v", msg)
	}
}
