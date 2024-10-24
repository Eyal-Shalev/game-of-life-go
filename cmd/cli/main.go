package main

import (
	"context"
	"fmt"
	"github.com/Eyal-Shalev/game-of-life-go/board"
	"github.com/Eyal-Shalev/game-of-life-go/runner"
	"github.com/dottedmag/must"
	"log/slog"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer slog.Info("Bye")

	ctx, closeCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer closeCtx()

	bChan := make(chan *board.Board)

	r := must.OK1(runner.New(
		runner.WithBaseContext(ctx),
		runner.WithBoardSize(10, 10),
		runner.WithClockInterval(time.Second),
		runner.WithListener(bChan),
		runner.WithInitFn(func(p board.Point) bool {
			return rand.Float32() > 0.5
		}),
	))
	must.OK(r.Start())
	defer doIgnore(r.Close)

	for {
		select {
		case b := <-bChan:
			printBoard(b)

		case <-ctx.Done():
			slog.Info("Shutting down")
			return
		}
	}
}

func printBoard(b *board.Board) {
	for range b.Rows() {
		fmt.Println()
	}

	for row := range b.Rows() {
		for col := range b.Columns() {
			if must.OK1(b.IsAlive(board.Point{Row: row, Column: col})) {
				fmt.Print("●")
			} else {
				fmt.Print("○")
			}
		}
		fmt.Println()
	}
}

func doIgnore(fn func() error) {
	_ = fn()
}