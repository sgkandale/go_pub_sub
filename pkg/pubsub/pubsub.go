package pubsub

import (
	"sync"
	"time"
)

const (
	PubsubChannelSize = 10_000
)

type PubSub struct {
	// topics is a map of topic name to listener.
	// A dedicated goroutine is spawned for each listener,
	// which will forward the messages to subscribers.
	topics map[string]*TopicListener

	mutex sync.Mutex
}

// New creates a new PubSub instance.
// This will spawn a goroutine to remove unused topics periodically.
func New() *PubSub {
	newPubSub := &PubSub{
		topics: make(map[string]*TopicListener, 5),
		mutex:  sync.Mutex{},
	}

	go newPubSub.RemoveUnusedTopicsLoop()

	return newPubSub
}

// Publish sends the message to the topic and all of its subscribers.
func (p *PubSub) Publish(topic string, msg *Message) {
	topicListener, exist := p.topics[topic]
	if !exist {
		topicListener = NewTopicListener(topic)
		p.mutex.Lock()
		p.topics[topic] = topicListener
		p.mutex.Unlock()
	}
	topicListener.Send(msg)
}

// Subscribe creates a new subscriber for the topic.
func (p *PubSub) Subscribe(topic string) *Subscriber {
	topicListener, exist := p.topics[topic]
	if !exist {
		topicListener = NewTopicListener(topic)
		p.mutex.Lock()
		p.topics[topic] = topicListener
		p.mutex.Unlock()
	}
	return topicListener.AddNewSubscriber()
}

// Unsubscribe removes the subscriber from the topic.
func (p *PubSub) Unsubscribe(subscriber *Subscriber) {
	if subscriber == nil {
		return
	}

	topicListener, exist := p.topics[subscriber.topic]
	if !exist {
		return
	}
	topicListener.RemoveSubscriber(subscriber)
}

// RemoveUnusedTopicsLoop periodically checks for all the topics in pubsub
// and removes those with no active subscribers.
func (p *PubSub) RemoveUnusedTopicsLoop() {
	for {
		time.Sleep(time.Second)
		p.mutex.Lock()
		for topic, topicListener := range p.topics {
			if topicListener.SubscribersCount() == 0 {
				delete(p.topics, topic)
			}
		}
		p.mutex.Unlock()
	}
}
