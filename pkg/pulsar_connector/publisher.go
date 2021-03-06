package pulsar_connector

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
	"go.uber.org/zap"
	"time"
)

type PulsarPublisher struct {
	log      *zap.SugaredLogger
	ctx      context.Context
	cancelFn context.CancelFunc
	client   pulsar.Client
	producer pulsar.Producer
}

func NewPublisher(logger *zap.SugaredLogger, ctx context.Context, cancelFn context.CancelFunc, topic string) (pubsub.Publisher, error) {
	// Create client
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://localhost:36650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})

	if err != nil {
		return nil, fmt.Errorf("could not instantiate Pulsar client: %w", err)
	}

	// Create producer
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})

	if err != nil {
		return nil, fmt.Errorf("could not create producer: %w", err)
	}

	p := &PulsarPublisher{
		log:      logger,
		ctx:      ctx,
		cancelFn: cancelFn,
		client:   client,
		producer: producer,
	}

	return p, nil
}

func (m *PulsarPublisher) SendMsg(msg string) error {
	_, err := m.producer.Send(m.ctx, &pulsar.ProducerMessage{
		Payload: []byte(msg),
	})

	if err != nil {
		m.cancelFn()
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (m *PulsarPublisher) Close() {
	m.log.Info("Closing pulsar publisher")
	m.producer.Close()
	m.client.Close()
}
