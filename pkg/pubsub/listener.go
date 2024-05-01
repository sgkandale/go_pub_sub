package pubsub

import (
	"log"
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

// NewTopicListener creates new topic listener with all fields initialised.
// This will also call listen method which listens for incoming messages.
func NewTopicListener(topic string) *TopicListener {
	newTopicListener := &TopicListener{
		topic:       topic,
		msgChannel:  make(chan *Message, PubsubChannelSize),
		subscribers: make(map[string]*Subscriber, PubsubChannelSize),
	}
	go newTopicListener.listen()
	return newTopicListener
}

// AddNewSubscriber add new subscriber to the topic listener's map
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

// RemoveSubscriber remove subscriber from the topic listener's map
func (t *TopicListener) RemoveSubscriber(subscriber *Subscriber) {
	t.mutex.Lock()
	delete(t.subscribers, subscriber.uuid)
	t.mutex.Unlock()
}

// Send sends message to the topic listener's msgChannel
// for it to be received by all the receivers
func (t *TopicListener) Send(msg *Message) {
	t.msgChannel <- msg
}

// listen listens for incoming messages on the msgChannel.
// All the incoming messages will be forwarded to each subscriber
// in the subscribers map.
// This method will exit when there are no subscribers left in map.
func (t *TopicListener) listen() {
	// timer to periodically check active subscribers
	subscriberCountCheck := time.NewTimer(time.Second * 5)
	defer log.Print("exiting listen")

	for {
		select {
		case msg := <-t.msgChannel:
			// forward message to each subscriber
			for _, subscriber := range t.subscribers {
				subscriber.Send(msg)
			}
		// exit listener if no subscribers are active
		case <-subscriberCountCheck.C:
			if t.SubscribersCount() == 0 {
				return
			}
		}
	}
}

// SubscribersCount returns the number of subscribers connected to this topic
func (t *TopicListener) SubscribersCount() int {
	return len(t.subscribers)
}
