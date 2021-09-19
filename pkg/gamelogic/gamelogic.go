package gamelogic

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	zombie2 "github.com/yngvark/gridwalls3/source/zombie-go/pkg/zombie"
	"go.uber.org/zap"

	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/worldmap"
)

// GameLogic knows how to run the game
type GameLogic struct {
	log       *zap.SugaredLogger
	publisher pubsub.Publisher
	ctx       context.Context
	generator *Generator
}

// Run continuously publishes messages with game logic events. It blocks until signalled to stop.
func (l *GameLogic) Run() {
	l.log.Info("Starting to generate...")

	ticker := time.NewTicker(time.Second * 1) //nolint:gomnd
	defer ticker.Stop()

	for {
		select {
		case <-l.ctx.Done():
			l.log.Info("Zombie generator stopped. TODO this doesn't work, and it's okay, but why doesn't the" +
				" ticker stop??")

			return
		case <-ticker.C:
			zombieMove, err := l.generator.Next()
			if err != nil {
				l.log.Info("could not generate next message: %w", err)
				return
			}

			zombieMoveJSON, err := json.Marshal(zombieMove)
			if err != nil {
				l.log.Info("could not marshal zombie move: %w", err)
				return
			}

			err = l.publisher.SendMsg(string(zombieMoveJSON))
			if err != nil {
				l.log.Error("-- WE SHOULD NEVER SEE THIS I THINK, PUBLISHER FAILED AND SHOULD CANCEL THE CONTEXT")
				return
			}
		}
	}
}

// NewGameLogic returns a new GameLogic
func NewGameLogic(ctx context.Context, logger *zap.SugaredLogger, publisher pubsub.Publisher) *GameLogic {
	m := worldmap.New(20, 10)                                                //nolint:gomnd
	zombie := zombie2.NewZombie("1", 10, 5, m, rand.New(rand.NewSource(45))) //nolint:gosec,gomnd

	return &GameLogic{
		log:       logger,
		publisher: publisher,
		ctx:       ctx,
		generator: NewGenerator(zombie),
	}
}
