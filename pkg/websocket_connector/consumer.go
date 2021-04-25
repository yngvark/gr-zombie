package websocket_connector

import (
	"context"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
	"go.uber.org/zap"
)

type WebsocketConsumer struct {
	log        *zap.SugaredLogger
	ctx        context.Context
	subscriber chan<- string
}

func NewConsumer(
	logger *zap.SugaredLogger,
	ctx context.Context,
	topic string,
	subscriber chan<- string,
) (pubsub.Consumer, error) {
	// TODO Add stuff here

	// Create client
	c := &WebsocketConsumer{
		log:        logger,
		ctx:        ctx,
		subscriber: subscriber,
	}

	return c, nil
}

// ListenForMessages reads messages from the connected websocket. This function blocks until the context provided on
// creation is done.
func (c *WebsocketConsumer) ListenForMessages() {
	select {
	case msg := <-c.consumer.Chan():
		msgString := string(msg.Payload())
		c.subscriber <- msgString
	case <-c.ctx.Done():
		return
	}
}

func (c *WebsocketConsumer) Close() error {
	c.log.Info("Closing websocket consumer")

	c.consumer.Close()
	c.client.Close()

	return nil
}
