// Package websocket provides publish subscribe using websockets
package websocket

import (
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
	"go.uber.org/zap"
)

type websocketPublisher struct {
	logger *zap.SugaredLogger
}

// SendMsg sends messages
func (p websocketPublisher) SendMsg(msg string) error {
	p.logger.Infof("NOT sending: %s", msg)
	return nil
}

// Close closes the publisher
func (p websocketPublisher) Close() error {
	p.logger.Info("NOT closing")
	return nil
}

// NewPublisher returns a new publisher for websockets
func NewPublisher(logger *zap.SugaredLogger) pubsub.Publisher {
	return websocketPublisher{
		logger: logger,
	}
}
