package zombie

import (
	"fmt"
	"math/rand"

	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/worldmap"
)

type Zombie struct {
	Id       string
	X        int
	Y        int
	WorldMap *worldmap.WorldMap
	Rand     *rand.Rand
}

func NewZombie(id string, x int, y int, worldMap *worldmap.WorldMap, rnd *rand.Rand) *Zombie {
	return &Zombie{
		Id:       id,
		X:        x,
		Y:        y,
		WorldMap: worldMap,
		Rand:     rnd,
	}
}

func (z *Zombie) Move() (*Zombie, *Move, error) {
	newX, err := z.getNewCoordPart(z.X, worldmap.Axis.X)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get new x coordinate: %w", err)
	}

	newY, err := z.getNewCoordPart(z.Y, worldmap.Axis.Y)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get new x coordinate: %w", err)
	}

	newZ := NewZombie(z.Id, newX, newY, z.WorldMap, z.Rand)
	move := NewZombieMove(z.Id, newX, newY)

	return newZ, move, nil
}

func (z *Zombie) getNewCoordPart(currentValue int, axisType worldmap.AxisType) (int, error) {
	direction := z.Rand.Intn(3) - 1 //nolint:gomnd,gosec    // [-1, 1]
	suggestion := currentValue + direction

	isInMap, err := z.WorldMap.IsInMap(suggestion, axisType)
	if err != nil {
		return -1, fmt.Errorf("could not detect if value is within map: %w", err)
	}

	if isInMap {
		return suggestion, nil
	}

	return currentValue, nil
}
