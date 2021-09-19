package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	gamelogicPkg "github.com/yngvark/gridwalls3/source/zombie-go/pkg/gamelogic"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(fmt.Errorf("could not run game: %w", err))
	}

	fmt.Println("Main ended.")
}

func run() error {
	ctx, cancelFn := context.WithCancel(context.Background())
	osInterruptChan := make(chan os.Signal, 1)

	signal.Notify(osInterruptChan, os.Interrupt)

	// Don't listen for interrupts after program quits
	defer func() {
		signal.Stop(osInterruptChan)
		cancelFn()
	}()

	// Listen in the background (i.e. goroutine) if the OS interrupts our program.
	go cancelProgramIfOsInterrupts(ctx, osInterruptChan, cancelFn)

	gameOpts, err := newGameOpts(ctx, cancelFn, os.Getenv)
	if err != nil {
		return fmt.Errorf("creating dependencies: %w", err)
	}

	return runGameLogic(gameOpts)
}

func cancelProgramIfOsInterrupts(ctx context.Context, osInterruptChan chan os.Signal, cancelFn context.CancelFunc) {
	func() {
		select {
		case <-osInterruptChan:
			cancelFn()
		case <-ctx.Done():
			// Stop listening
			return
		}
	}()
}

func runGameLogic(o *GameOpts) error {
	// Create producer
	defer func() {
		err := o.publisher.Close()
		o.logger.Error(err)
	}()

	// Create consumer
	defer func() {
		err := o.consumer.Close()
		o.logger.Error(err)
	}()

	// Create game
	gameLogic := gamelogicPkg.NewGameLogic(o.context, o.logger, o.publisher)

	// Wait until some external orchestrator sends a "start" message
	go o.consumer.ListenForMessages()

	o.logger.Info("Waiting for start message...")

	select {
	case msg := <-o.consumer.SubscriberChannel():
		o.logger.Info("Waiting for start message... Received: %s", msg)

		if msg == "start" {
			break
		}
	case <-o.context.Done():
		o.logger.Info("Aborted waiting for game to start")
		return nil
	}

	o.logger.Info("Running game")
	gameLogic.Run()

	return nil
}
