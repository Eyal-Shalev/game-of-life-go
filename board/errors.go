package board

import (
	"fmt"

	. "github.com/Eyal-Shalev/game-of-life-go/internal"
)

const errorPrefix = "game_of_life/board."
const baseNewBoardError = StringError("game_of_life/board.NewBoardError")

func newBoardError(origin error) error {
	return JoinErrors(baseNewBoardError, origin)
}

type NewBoardError struct {
	Origin error
}

func (e NewBoardError) Error() string {
	return fmt.Sprintf(errorPrefix+"NewBoardError: %s", e.Origin.Error())
}

type InvalidPointError struct {
	Rows, Columns int
	Point         Point
	Origin        error
}

func (e InvalidPointError) Error() string {
	msg := fmt.Sprintf(errorPrefix+"InvalidPointError: point=%s board.Rows=%d board.Columns=%d", e.Point, e.Rows, e.Columns)
	if e.Origin != nil {
		msg += " origin=" + e.Origin.Error()
	}
	return msg
}

func (e InvalidPointError) Unwrap() error {
	return e.Origin
}
