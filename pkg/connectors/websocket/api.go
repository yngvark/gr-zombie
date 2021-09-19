// Package websocket handles publishing and subscribing with websockets
package websocket

import "github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"

// NewPublisher returns a new publisher for websockets
func NewPublisher() pubsub.Publisher {
	return nil
}

// NewConsumer returns a new consumer for websockets
func NewConsumer(<-chan string) pubsub.Consumer {
	return nil
}
