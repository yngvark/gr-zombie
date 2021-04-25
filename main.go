package main

import (
	"context"
	"fmt"
	gamelogicPkg "github.com/yngvark/gridwalls3/source/zombie-go/pkg/gamelogic"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/log2"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pulsar_connector"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/websocket_connector"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(fmt.Errorf("could not run game: %w\n", err))
	}

	fmt.Println("Main ended.")
}

func run() error {
	logger, err := log2.New()
	if err != nil {
		return fmt.Errorf("could not create logger: %w", err)
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	osInterruptChan := make(chan os.Signal, 1)

	signal.Notify(osInterruptChan, os.Interrupt)

	// Don't listen for interrupts after program quits
	defer func() {
		signal.Stop(osInterruptChan)
		cancelFn()
	}()

	// Listen in the background (i.e. goroutine) if the OS interrupts our program.
	go cancelProgramIfOsInterrupts(osInterruptChan, cancelFn, ctx)

	return runGameLogic(logger, ctx, cancelFn)
}

func cancelProgramIfOsInterrupts(osInterruptChan chan os.Signal, cancelFn context.CancelFunc, ctx context.Context) {
	func() {
		select {
		case <-osInterruptChan:
			cancelFn()
		case <-ctx.Done():
			// Stop listening
		}
	}()
}

func runGameLogic(logger *zap.SugaredLogger, ctx context.Context, cancelFn context.CancelFunc) error {
	// Create producer
	var producer pubsub.Publisher
	var err error

	networkType, networkTypeFound := os.LookupEnv("NETWORK_TYPE")

	if !networkTypeFound || networkType == "pulsar" {
		logger.Info("Using network type: Pulsar")
		producer, err = pulsar_connector.NewPublisher(logger, ctx, cancelFn, "zombie")
	} else if networkType == "websocket" {
		logger.Info("Using network type: Websocket")
		producer, err = websocket_connector.NewPublisher(logger, ctx, cancelFn, "zombie")
	}

	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	defer func() {
		err := producer.Close()
		if err != nil {
			logger.Error("closing producer: %w", err)
		}
	}()

	// Create consumer
	consumerChan := make(chan string)

	var consumer pubsub.Consumer

	if !networkTypeFound || networkType == "pulsar" {
		logger.Info("Using network type: Pulsar")
		consumer, err = pulsar_connector.NewConsumer(logger, ctx, "gameinit", consumerChan)
	} else if networkType == "websocket" {
		logger.Info("Using network type: Websocket")
		consumer, err = websocket_connector.NewConsumer(logger, ctx, "gameinit", consumerChan)
	}

	if err != nil {
		return fmt.Errorf("could not create consumer: %w", err)
	}

	defer func() {
		err := consumer.Close()
		if err != nil {
			logger.Error("closing producer: %w", err)
		}
	}()

	// Create game
	gameLogic := gamelogicPkg.NewGameLogic(logger, producer, ctx)

	// Wait until some external orchestrator sends a "start" message
	go consumer.ListenForMessages()

	logger.Info("Waiting for start message...")
	select {
	case msg := <-consumerChan:
		logger.Info("Waiting for start message... Received: %s", msg)
		if msg == "start" {
			break
		}
	case <-ctx.Done():
		logger.Info("Aborted waiting for game to start")
		return nil
	}

	logger.Info("Running game")
	gameLogic.Run()

	return nil
}
