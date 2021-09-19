// Package broadcast knows how to broadcast messages to subscribers
package broadcast

import "github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"

// Broker is used for sending (broadcasting) messages to a number of subscribers
type Broker struct {
	pubsub.Publisher
	subscribers []chan<- string
}

// AddSubscriber adds a Subscriber to its list of subscribers
func (n *Broker) AddSubscriber(subscriber chan<- string) {
	n.subscribers = append(n.subscribers, subscriber)
}

// SendMsg sends a message to all Subscriber-s
func (n *Broker) SendMsg(msg string) error {
	for _, subscriber := range n.subscribers {
		subscriber <- msg
	}

	return nil
}

// New returns a new Broker
func New() *Broker {
	return &Broker{
		subscribers: make([]chan<- string, 0),
	}
}
