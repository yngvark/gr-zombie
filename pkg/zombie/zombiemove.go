package zombie

type Move struct {
	Id string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

func NewZombieMove(id string, x int, y int) *Move {
	return &Move{
		Id: id,
		X:  x,
		Y:  y,
	}
}
