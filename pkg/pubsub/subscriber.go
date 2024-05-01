package pubsub

// Subscriber is an instance of the subscriber connected to this pubsub
type Subscriber struct {
	topic      string
	uuid       string
	msgChannel chan *Message
}

// Send pushes the message to the channel of Subscriber
// for it to be read by subscriber instance
func (s *Subscriber) Send(msg *Message) {
	if s.msgChannel != nil {
		s.msgChannel <- msg
	}
}

// Listen returns the channel of Subscriber
func (s *Subscriber) Listen() <-chan *Message {
	if s.msgChannel == nil {
		s.msgChannel = make(chan *Message, PubsubChannelSize)
	}
	return s.msgChannel
}
