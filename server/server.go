package server

import (
	"pubsub/config"
	"pubsub/pkg/pubsub"
)

type NewServerRequest struct {
	Cfg    *config.ServerConfig
	PubSub *pubsub.PubSub
}
