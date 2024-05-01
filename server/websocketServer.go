package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"pubsub/pkg/pubsub"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	httpServer        *http.Server
	websocketUpgrader websocket.Upgrader
	// internals
	port   int
	pubsub *pubsub.PubSub
}

// NewWebsocket creats a new websocket server instance
func NewWebsocket(ctx context.Context, req *NewServerRequest) *WebsocketServer {
	if req == nil {
		log.Printf("[ERROR] nil request to create websocket server")
		return nil
	}
	if req.Cfg == nil {
		log.Printf("[ERROR] nil config to create websocket server")
		return nil
	}
	if req.PubSub == nil {
		log.Printf("[ERROR] nil pubsub value to create websocket server")
		return nil
	}

	websocketServer := &WebsocketServer{
		httpServer: &http.Server{
			Addr: fmt.Sprintf(":%d", req.Cfg.Port),
		},
		port:   req.Cfg.Port,
		pubsub: req.PubSub,
		websocketUpgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	return websocketServer
}

// Serve starts the websocket server
func (s *WebsocketServer) Serve() error {
	log.Printf("[INFO] starting websocket server on port %d", s.port)
	http.HandleFunc("/ws", s.PubSubHandler)
	return s.httpServer.ListenAndServe()
}

// Stop will attempt to stop the websocket server
func (s *WebsocketServer) Stop() {
	if s.httpServer != nil {
		log.Printf("[WARN] stopping websocket server on port %d", s.port)
		err := s.httpServer.Shutdown(context.Background())
		if err != nil {
			log.Printf("[ERROR] failed to stop websocket server on port %d", s.port)
		} else {
			log.Printf("[WARN] websocket server on port %d is stopped", s.port)
		}
	}
}

type WebsocketInitRequest struct {
	ClientType string `json:"client_type"`
	Topic      string `json:"topic"`
}

func (s *WebsocketServer) PubSubHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] upgrading connection to websocket : %s", err.Error())
		return
	}

	// get initial message
	msgType, msgBody, err := conn.ReadMessage()
	if err != nil {
		log.Printf("[ERROR] reading initial message in websocket connection : %s", err.Error())
		err = conn.Close()
		if err != nil {
			log.Printf("[ERROR] closing websocket connection after error in reading initial message : %s", err.Error())
		}
		return
	}

	var initReq WebsocketInitRequest

	// check msg type
	switch msgType {
	case websocket.TextMessage:
		err = json.Unmarshal(msgBody, &initReq)
		if err != nil {
			log.Printf("[ERROR] json unmarshalling initial message in websocket connection : %s", err.Error())
			err = conn.Close()
			if err != nil {
				log.Printf("[ERROR] closing websocket connection after error in reading initial message : %s", err.Error())
			}
		}
	default:
		log.Printf("[ERROR] unknown message type %d in websocket connection, closing connection", msgType)
		err = conn.Close()
		if err != nil {
			log.Printf("[ERROR] closing websocket connection after error in reading initial message : %s", err.Error())
		}
	}

	if strings.EqualFold(initReq.ClientType, "publisher") {
		// listen for incoming messages
		for {
			newMsg := pubsub.Message{}
			err = conn.ReadJSON(&newMsg)
			if err != nil {
				if err == websocket.ErrCloseSent {
					log.Printf("[WARN] closing websocket connection after error in reading initial message : %s", err.Error())
					break
				}
				log.Printf("[ERROR] reading message in websocket connection : %s", err.Error())
				continue
			}
			if newMsg.Topic == "" {
				log.Printf("[ERROR] publisher client type, but no topic specified in message")
				continue
			}
			if len(newMsg.Body) == 0 {
				log.Printf("[ERROR] publisher client type, but no body specified in message")
				continue
			}
			newMsg.Timestamp = time.Now().UTC()
			s.pubsub.Publish(newMsg.Topic, &newMsg)
		}
	} else if strings.EqualFold(initReq.ClientType, "subscriber") {
		if initReq.Topic == "" {
			log.Printf("[ERROR] subscriber client type, but no topic specified, closing connection")
			err = conn.Close()
			if err != nil {
				log.Printf("[ERROR] closing websocket connection after error in reading initial message : %s", err.Error())
			}
		}
		log.Printf("[INFO] new subscriber connected in websocket for topic : '%s'", initReq.Topic)

		// register subscriber in pubsub
		subscriber := s.pubsub.Subscribe(initReq.Topic)
		defer s.pubsub.Unsubscribe(subscriber)

		// send messages from the topic
		for msg := range subscriber.Listen() {
			err = conn.WriteJSON(msg)
			if err != nil {
				if err == websocket.ErrCloseSent {
					log.Printf("[WARN] closing websocket connection after error in reading initial message : %s", err.Error())
					break
				}
				log.Printf("[ERROR] writing message in websocket connection : %s", err.Error())
				continue
			}
		}
	} else {
		err = conn.Close()
		if err != nil {
			log.Printf("[ERROR] closing websocket connection after error in reading initial message : %s", err.Error())
		}
	}
}
