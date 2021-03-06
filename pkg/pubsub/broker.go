package pubsub

type Broker struct {
	Publisher
	subscribers []Subscriber
}

func (n *Broker) Subscribe(l Subscriber) {
	n.subscribers = append(n.subscribers, l)
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make([]Subscriber, 0),
	}
}

func (n *Broker) SendMsg(msg string) error {
	for _, l := range n.subscribers {
		l.MsgReceived(msg)
	}

	return nil
}
