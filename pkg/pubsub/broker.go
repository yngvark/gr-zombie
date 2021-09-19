package pubsub

// Subscriber receives messages from a Broker. A Broker listens for incomming messages and forwards (broadcasts) them
// to all Subscriber-s.
type Subscriber interface {
	MsgReceived(msg string)
}

// Broker is used for sending (broadcasting) messages to a number of subscribers
type Broker struct {
	Publisher
	subscribers []Subscriber
}

// AddSubscriber adds a Subscriber to its list of subscribers
func (n *Broker) AddSubscriber(l Subscriber) {
	n.subscribers = append(n.subscribers, l)
}

// SendMsg sends a message to all Subscriber-s
func (n *Broker) SendMsg(msg string) error {
	for _, l := range n.subscribers {
		l.MsgReceived(msg)
	}

	return nil
}

// NewBroker returns a new Broker
func NewBroker() *Broker {
	return &Broker{
		subscribers: make([]Subscriber, 0),
	}
}
