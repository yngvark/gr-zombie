package gamelogic

import (
	"fmt"
	zombiePkg "github.com/yngvark/gridwalls3/source/zombie-go/pkg/zombie"
)

type Generator struct {
	zombie *zombiePkg.Zombie
}

func NewGenerator(initialZombie *zombiePkg.Zombie) *Generator {
	return &Generator{
		zombie: initialZombie,
	}
}

func (g *Generator) Next() (*zombiePkg.Move, error) {
	z, move, err := g.zombie.Move()
	if err != nil {
		return nil, fmt.Errorf("could not move zombie: %w", err)
	}

	g.zombie = z

	return move, nil
}
