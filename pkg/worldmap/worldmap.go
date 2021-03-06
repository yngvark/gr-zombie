package worldmap

import (
	"fmt"
)

type AxisType int

var Axis = struct { //nolint:gochecknoglobals
	X AxisType
	Y AxisType
}{
	X: 1, //nolint:gomnd
	Y: 2, //nolint:gomnd
}

type WorldMap struct {
	MinX int
	MaxX int
	MinY int
	MaxY int
}

func New(maxX int, maxY int) *WorldMap {
	return &WorldMap{
		MinX: 0,
		MaxX: maxX,
		MinY: 0,
		MaxY: maxY,
	}
}

func (m *WorldMap) IsInMap(axisValue int, axis AxisType) (bool, error) {
	switch axis {
	case Axis.X:
		return m.xIsInMap(axisValue), nil
	case Axis.Y:
		return m.yIsInMap(axisValue), nil
	default:
		return false, fmt.Errorf("not a valid Axis type: %d", axis)
	}
}

func (m *WorldMap) xIsInMap(x int) bool {
	return x >= m.MinX && x <= m.MaxX
}

func (m *WorldMap) yIsInMap(y int) bool {
	return y >= m.MinY && y <= m.MaxY
}
