package board

import (
	"fmt"
)

type Point struct {
	Row, Column int
}

// Neighbours returns the list of possible neighboring points. Some of the might be invalid in the game board.
func (p Point) Neighbours() [8]Point {
	return [8]Point{
		{p.Row - 1, p.Column - 1}, {p.Row - 1, p.Column}, {p.Row - 1, p.Column + 1},
		{p.Row, p.Column - 1}, {p.Row, p.Column + 1},
		{p.Row + 1, p.Column - 1}, {p.Row + 1, p.Column}, {p.Row + 1, p.Column + 1},
	}
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.Row, p.Column)
}

func (p Point) GoString() string {
	return fmt.Sprintf("board.Point{Row:%d, Column:%d}", p.Row, p.Column)
}
