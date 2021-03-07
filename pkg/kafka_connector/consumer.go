package kafka_connector

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/segmentio/kafka-go"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
	"go.uber.org/zap"
	"time"
)

type KafkaConsumer struct {
	log *zap.SugaredLogger
	ctx *context.Context
	//client     kafka.something
	//consumer   pulsar.Consumer
	subscriber chan<- string
}

func NewConsumer(
	logger *zap.SugaredLogger,
	ctx context.Context,
	topic string,
	subscriber chan<- string,
) (pubsub.Consumer, error) {
	// Create client

	// TODO: https://github.com/segmentio/kafka-go#reader-

	c := &KafkaConsumer{
		log:        logger,
		ctx:        ctx,
		client:     client,
		consumer:   consumer,
		subscriber: subscriber,
	}

	return c, nil
}

// ListenForMessages reads messages from Pulsar. This function blocks until the context provided on creation is done.
func (c *KafkaConsumer) ListenForMessages() {
	select {
	case msg := <-c.consumer.Chan():
		msgString := string(msg.Payload())
		c.subscriber <- msgString
	case <-c.ctx.Done():
		return
	}
}

func (c *KafkaConsumer) Close() {
	c.log.Info("Closing pulsar consumer")

	c.consumer.Close()
	c.client.Close()
}
