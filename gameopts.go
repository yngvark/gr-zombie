package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/connectors/websocket/oslookup"

	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/connectors/pulsar"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/connectors/websocket"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/log2"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
	"go.uber.org/zap"
)

// GameOpts contains various dependencies
type GameOpts struct {
	context   context.Context
	logger    *zap.SugaredLogger
	publisher pubsub.Publisher
	consumer  pubsub.Consumer
}

type getEnv func(key string) string

func newGameOpts(ctx context.Context, cancelFn context.CancelFunc, getEnv getEnv) (*GameOpts, error) {
	logger, err := log2.New()
	if err != nil {
		return nil, fmt.Errorf("could not create logger: %w", err)
	}

	var publisher pubsub.Publisher

	var consumer pubsub.Consumer

	subscriberChan := make(chan string)

	if getEnv("GAME_QUEUE_TYPE") == "websocket" {
		publisher, consumer, err = pubSubForWebsocket(ctx, cancelFn, logger, subscriberChan)
		if err != nil {
			return nil, fmt.Errorf("creating websocket connectors: %w", err)
		}
	} else {
		publisher, consumer, err = pubSubForPulsar(ctx, cancelFn, logger, subscriberChan)
		if err != nil {
			return nil, fmt.Errorf("creating pulsar connectors: %w", err)
		}
	}

	return &GameOpts{
		context:   ctx,
		logger:    logger,
		publisher: publisher,
		consumer:  consumer,
	}, nil
}

const allowedCorsOriginsEnvVarKey = "ALLOWED_CORS_ORIGINS"

func pubSubForWebsocket(
	ctx context.Context,
	cancelFn context.CancelFunc,
	logger *zap.SugaredLogger,
	subscriber chan string,
) (pubsub.Publisher, pubsub.Consumer, error) {
	corsHelper := oslookup.NewCORSHelper(logger)

	allowedCorsOrigins, err := corsHelper.GetAllowedCorsOrigins(os.LookupEnv, allowedCorsOriginsEnvVarKey)
	if err != nil {
		return nil, nil, fmt.Errorf("getting allowed CORS origins: %w", err)
	}

	corsHelper.PrintAllowedCorsOrigins(allowedCorsOrigins)

	publisher, consumer := websocket.New(ctx, cancelFn, logger, subscriber, allowedCorsOrigins)

	return publisher, consumer, nil
}

func pubSubForPulsar(
	ctx context.Context,
	cancelFn context.CancelFunc,
	logger *zap.SugaredLogger,
	subscriberChan chan string,
) (pubsub.Publisher, pubsub.Consumer, error) {
	// Publisher
	publisher, err := pulsar.NewPublisher(ctx, cancelFn, logger, "zombie")
	if err != nil {
		return nil, nil, fmt.Errorf("creating publisher: %w", err)
	}

	// Consumer
	consumer, err := pulsar.NewConsumer(ctx, logger, "gameinit", subscriberChan)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create consumer: %w", err)
	}

	return publisher, consumer, err
}
