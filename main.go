package main

import (
	"context"
	"fmt"
	gamelogicPkg "github.com/yngvark/gridwalls3/source/zombie-go/pkg/gamelogic"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/log2"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pulsar_connector"
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

	producer, err = pulsar_connector.NewPublisher(logger, ctx, cancelFn, "zombie")
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	defer producer.Close()

	// Create consumer
	consumerChan := make(chan string)

	var consumer pubsub.Consumer
	consumer, err = pulsar_connector.NewConsumer(logger, ctx, "gameinit", consumerChan)
	if err != nil {
		return fmt.Errorf("could not create consumer: %w", err)
	}

	defer consumer.Close()

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
