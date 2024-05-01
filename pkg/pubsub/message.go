package pubsub

import (
	"time"

	"pubsub/pkg/pubsubPB"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	Uuid      string            `json:"uuid"`
	Timestamp time.Time         `json:"timestamp"`
	Body      []byte            `json:"body"`
	Topic     string            `json:"topic"`
	Tags      map[string]string `json:"tags"`
}

func NewMessageFromPB(pb *pubsubPB.Message) *Message {
	if pb == nil {
		return nil
	}
	return &Message{
		Uuid:      pb.Uuid,
		Timestamp: pb.Timestamp.AsTime(),
		Body:      pb.Body,
		Topic:     pb.Topic,
		Tags:      pb.Tags,
	}
}

func MessageToPB(msg *Message) *pubsubPB.Message {
	if msg == nil {
		return nil
	}
	return &pubsubPB.Message{
		Uuid:      msg.Uuid,
		Timestamp: timestamppb.New(msg.Timestamp),
		Body:      msg.Body,
		Topic:     msg.Topic,
		Tags:      msg.Tags,
	}
}
