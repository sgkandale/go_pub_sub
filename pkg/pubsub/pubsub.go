package pubsub

const (
	PubsubChannelSize = 10_000
)

type PubSub struct {
	// topics is a map of topic name to listener.
	// A dedicated goroutine is spawned for each listener,
	// which will forward the messages to subscribers.
	topics map[string]*TopicListener
}

func New() *PubSub {
	return &PubSub{
		topics: make(map[string]*TopicListener, 5),
	}
}

func (p *PubSub) Publish(topic string, msg *Message) error {
	topicListener, exist := p.topics[topic]
	if !exist {
		topicListener = NewTopicListener(topic)
		p.topics[topic] = topicListener
	}
	topicListener.Send(msg)
	return nil
}

func (p *PubSub) Subscribe(topic string) *Subscriber {
	topicListener, exist := p.topics[topic]
	if !exist {
		topicListener = NewTopicListener(topic)
		p.topics[topic] = topicListener
	}
	return topicListener.AddNewSubscriber()
}

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
