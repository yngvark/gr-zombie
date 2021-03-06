package pulsar_connector

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
	"go.uber.org/zap"
	"time"
)

type PulsarConsumer struct {
	log        *zap.SugaredLogger
	ctx        context.Context
	client     pulsar.Client
	consumer   pulsar.Consumer
	subscriber chan<- string
}

func NewConsumer(
	logger *zap.SugaredLogger,
	ctx context.Context,
	topic string,
	subscriber chan<- string,
) (pubsub.Consumer, error) {
	// Create client
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://localhost:36650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})

	if err != nil {
		return nil, fmt.Errorf("could not instantiate Pulsar client: %w", err)
	}

	// Create consumer
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: "mysub2",
		Type:             pulsar.Exclusive,
	})

	c := &PulsarConsumer{
		log:        logger,
		ctx:        ctx,
		client:     client,
		consumer:   consumer,
		subscriber: subscriber,
	}

	return c, nil
}

// ListenForMessages reads messages from Pulsar. This function blocks until the context provided on creation is done.
func (c *PulsarConsumer) ListenForMessages() {
	select {
	case msg := <-c.consumer.Chan():
		msgString := string(msg.Payload())
		c.subscriber <- msgString
	case <-c.ctx.Done():
		return
	}
}

func (c *PulsarConsumer) Close() {
	c.log.Info("Closing pulsar consumer")

	c.consumer.Close()
	c.client.Close()
}
