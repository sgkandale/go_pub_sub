package pubsub

type Subscriber struct {
	topic      string
	uuid       string
	msgChannel chan *Message
}

func (s *Subscriber) Send(msg *Message) {
	if s.msgChannel != nil {
		s.msgChannel <- msg
	}
}

func (s *Subscriber) Listen() <-chan *Message {
	if s.msgChannel == nil {
		s.msgChannel = make(chan *Message, PubsubChannelSize)
	}
	return s.msgChannel
}
