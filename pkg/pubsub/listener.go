package pubsub

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type TopicListener struct {
	// topic name
	topic string

	// msgChannel is the channel of messages,
	// which will receive messges from pubsub
	// and listener will forward them to each subscriber
	msgChannel chan *Message

	// subscribers is a map of all the subscribers connected to this topic,
	// key is the subscriber's uuid
	subscribers map[string]*Subscriber

	// to lock the listener while adding or removing subscribers
	mutex sync.Mutex
}

func NewTopicListener(topic string) *TopicListener {
	newTopicListener := &TopicListener{
		topic:       topic,
		msgChannel:  make(chan *Message, PubsubChannelSize),
		subscribers: make(map[string]*Subscriber, PubsubChannelSize),
	}
	go newTopicListener.Listen()
	return newTopicListener
}

func (t *TopicListener) AddNewSubscriber() *Subscriber {
	newSubscriber := &Subscriber{
		topic: t.topic,
		uuid:  uuid.NewString(),
		// DO NOT initialise the subscriber msgchannel here
	}
	t.mutex.Lock()
	t.subscribers[newSubscriber.uuid] = newSubscriber
	t.mutex.Unlock()

	return newSubscriber
}

func (t *TopicListener) RemoveSubscriber(subscriber *Subscriber) {
	t.mutex.Lock()
	delete(t.subscribers, subscriber.uuid)
	t.mutex.Unlock()
}

func (t *TopicListener) Send(msg *Message) {
	t.msgChannel <- msg
}

func (t *TopicListener) Listen() {
	// timer to periodically check active subscribers
	subscriberCountCheck := time.NewTimer(time.Second * 5)

	for {
		select {
		case msg := <-t.msgChannel:
			// forward message to each subscriber
			for _, subscriber := range t.subscribers {
				subscriber.Send(msg)
			}
		// exit listener if no subscribers are active
		case <-subscriberCountCheck.C:
			if len(t.subscribers) == 0 {
				return
			}
		}
	}
}
