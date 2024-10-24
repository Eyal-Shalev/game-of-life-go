package board

import (
	"fmt"
	"iter"

	"github.com/Eyal-Shalev/bitmap-go"
	. "github.com/Eyal-Shalev/game-of-life-go/internal"
	"github.com/dottedmag/must"
)

type Board struct {
	data *bitmap.BitMap
	rows int
}

func (b *Board) Clone() *Board {
	return &Board{
		data: b.data.Clone(),
		rows: b.rows,
	}
}

func (b *Board) Rows() int {
	if b == nil {
		return 0
	}
	return b.rows
}

func (b *Board) Columns() int {
	if b == nil {
		return 0
	}
	return b.data.Length() / b.rows
}

func (b *Board) Data() *bitmap.BitMap {
	return b.data.Clone()
}

func (b *Board) neighbours(p Point) []Point {
	if b == nil {
		return nil
	}
	possibleNeighbours := p.Neighbours()
	neighbours := make([]Point, 0, len(possibleNeighbours))
	for _, possibleNeighbor := range possibleNeighbours {
		if b.isValidPoint(possibleNeighbor) {
			neighbours = append(neighbours, possibleNeighbor)
		}
	}
	return neighbours
}

func (b *Board) isValidPoint(p Point) bool {
	if b == nil {
		return false
	}
	return p.Row >= 0 && p.Row < b.rows && p.Column >= 0 && p.Column < b.Columns()
}

func (b *Board) pointToBitPosition(p Point) int {
	if b == nil {
		return 0
	}
	return p.Row*b.Columns() + p.Column
}

func (b *Board) IsAlive(p Point) (bool, error) {
	if b == nil {
		return false, &NilPointerError{Target: b}
	}
	if !b.isValidPoint(p) {
		return false, &InvalidPointError{Rows: b.Rows(), Columns: b.Columns(), Point: p}
	}
	bitPos := b.pointToBitPosition(p)

	isAlive, err := b.data.IsSet(bitPos)
	if err != nil {
		return false, &InvalidPointError{Rows: b.Rows(), Columns: b.Columns(), Point: p, Origin: err}
	}

	return isAlive, nil
}

func (b *Board) IsDead(p Point) (bool, error) {
	isAlive, err := b.IsAlive(p)
	return !isAlive, err
}

func (b *Board) countLivingNeighbours(p Point) (int, error) {
	result := 0
	for _, neighbour := range b.neighbours(p) {
		isAlive, err := b.IsAlive(neighbour)
		if err != nil {
			return 0, fmt.Errorf("error counting living neighbours: %w", err)
		}
		if isAlive {
			result++
		}
	}
	return result, nil
}

func (b *Board) IsAliveNextCycle(p Point) (bool, error) {
	pIsAlive, err := b.IsAlive(p)
	if err != nil {
		return false, fmt.Errorf("IsAliveNextCycle: %w", err)
	}

	livingNeighboursCount, err := b.countLivingNeighbours(p)
	if err != nil {
		return false, fmt.Errorf("IsAliveNextCycle: %w", err)
	}

	// Any live cell with two or three live neighbours lives on to the next generation.
	survivalRule := pIsAlive && (livingNeighboursCount == 2 || livingNeighboursCount == 3)

	// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
	reproductionRule := !pIsAlive && livingNeighboursCount == 3

	return survivalRule || reproductionRule, nil
}

func (b *Board) IsDeadNextCycle(p Point) (bool, error) {
	isAlive, err := b.IsAliveNextCycle(p)
	return !isAlive, err
}
func (b *Board) Set(p Point, isAlive bool) error {
	bitPosition := b.pointToBitPosition(p)
	err := b.data.SetVal(bitPosition, isAlive)
	if err != nil {
		return &InvalidPointError{Rows: b.rows}
	}
	return nil
}

func (b *Board) NextBoard() (*Board, error) {
	nextBoard := must.OK1(New(b.Rows(), b.Columns()))

	for row := 0; row < b.rows; row++ {
		for column := 0; column < b.Columns(); column++ {
			p := Point{Row: row, Column: column}
			isAliveNext, loopErr := b.IsAliveNextCycle(p)
			if loopErr != nil {
				return nil, loopErr
			}
			if !isAliveNext {
				continue
			}
			loopErr = nextBoard.Set(p, isAliveNext)
			if loopErr != nil {
				return nil, loopErr
			}
		}
	}

	return nextBoard, nil
}

func (b *Board) Points() iter.Seq[Point] {
	return func(yield func(Point) bool) {
		for row := range b.Rows() {
			for column := range b.Columns() {
				if !yield(Point{Row: row, Column: column}) {
					return
				}
			}
		}
	}
}

func New(rows, cols int) (*Board, error) {
	data, err := bitmap.New(rows * cols)
	if err != nil {
		return nil, newBoardError(err)
	}

	return &Board{data: data, rows: rows}, nil
}

func NewFromData(rows int, data *bitmap.BitMap) (*Board, error) {
	if data.Length() < rows {
		return nil, fmt.Errorf("data is too small (%d < %d)", data.Length(), rows)
	}

	return &Board{data: data, rows: rows}, nil
}
