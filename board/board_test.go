package board_test

import (
	"testing"

	. "github.com/Eyal-Shalev/game-of-life-go/board"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFoo(t *testing.T) {
	b1, err := New(3, 3)
	require.NoError(t, err)
	// 0 1 0
	// 0 1 0
	// 0 1 0
	require.NoError(t, b1.Set(Point{Row: 0, Column: 1}, true))
	require.NoError(t, b1.Set(Point{Row: 1, Column: 1}, true))
	require.NoError(t, b1.Set(Point{Row: 2, Column: 1}, true))

	b2, err := b1.NextBoard()
	require.NoError(t, err)

	// 0 0 0
	// 1 1 1
	// 0 0 0
	expectedB2LivingPoints := []Point{
		{Row: 1, Column: 0},
		{Row: 1, Column: 1},
		{Row: 1, Column: 2},
	}
	for row := range b2.Rows() {
		for column := range b2.Columns() {
			p := Point{Row: row, Column: column}
			isAlive, err := b2.IsAlive(p)
			require.NoError(t, err)
			if isAlive {
				assert.Contains(t, expectedB2LivingPoints, p)
			} else {
				assert.NotContains(t, expectedB2LivingPoints, p)
			}
		}
	}
}
