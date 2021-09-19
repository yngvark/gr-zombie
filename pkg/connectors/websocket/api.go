// Package websocket handles publishing and subscribing with websockets
package websocket

import "github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"

func NewPublisher() pubsub.Publisher {
	return nil
}

func NewConsumer(<-chan string) pubsub.Consumer {
	return nil
}
