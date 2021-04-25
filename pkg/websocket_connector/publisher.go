package websocket_connector

import (
	"context"
	"fmt"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/gamelogic"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
	"go.uber.org/zap"
	"net/http"
)

type WebsocketPublisher struct {
	log      *zap.SugaredLogger
	ctx      context.Context
	cancelFn context.CancelFunc
}

func NewPublisher(logger *zap.SugaredLogger, ctx context.Context, cancelFn context.CancelFunc, topic string) (pubsub.Publisher, error) {
	// TODO Something goes here

	broker := pubsub.NewBroker()
	stopGamelogicChannel := make(chan bool)

	var publisher pubsub.Publisher = broker
	httpHandler := NewHTTPHandler(m.log, allowedCorsOrigins, publisher, stopGamelogicChannel)
	http.Handle("/zombie", httpHandler)

	var messageSender pubsub.Publisher = httpHandler
	gameLogic := gamelogic.NewGameLogic(m.log, messageSender, stopGamelogicChannel, nil)

	broker.Subscribe(gameLogic)

	p := &WebsocketPublisher{
		log:      logger,
		ctx:      ctx,
		cancelFn: cancelFn,
	}

	return p, nil
}

func (m *WebsocketPublisher) SendMsg(msg string) error {
	// TODOD Put something here
	if err != nil {
		m.cancelFn()
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (m *WebsocketPublisher) Close() error {
	m.log.Info("Closing websocket publisher")

	return nil
}
