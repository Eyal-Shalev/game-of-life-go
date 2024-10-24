package runner

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/Eyal-Shalev/game-of-life-go/board"
)

type option func(o *options)

func WithClockInterval(interval time.Duration) option {
	return func(o *options) {
		o.clockInterval = interval
	}
}

func WithBaseContext(ctx context.Context) option {
	return func(o *options) {
		o.baseCtx = ctx
	}
}

func WithCloseGracePeriod(duration time.Duration) option {
	return func(o *options) {
		o.closeGracePeriod = duration
	}
}

func WithBoardSize(rows, columns int) option {
	return func(o *options) {
		o.boardRows, o.boardColumns = rows, columns
	}
}

func WithListener(listener chan<- *board.Board) option {
	return func(o *options) {
		o.listeners = append(o.listeners, listener)
	}
}

func WithInitFn(fn func(point board.Point) bool) option {
	return func(o *options) {
		o.initFn = fn
	}
}

func WithInitBoard(initBoard *board.Board) option {
	return func(o *options) {
		o.initBoard = initBoard
	}
}

type options struct {
	clockInterval           time.Duration
	baseCtx                 context.Context
	closeGracePeriod        time.Duration
	boardRows, boardColumns int
	listeners               []chan<- *board.Board
	initFn                  func(point board.Point) bool
	initBoard               *board.Board
}

var defaultOptions = options{
	clockInterval:    300 * time.Millisecond,
	baseCtx:          context.Background(),
	closeGracePeriod: time.Minute,
	boardRows:        100,
	boardColumns:     100,
}

//go:generate stringer -type State
type State uint32

const (
	Stopped State = iota
	Starting
	Running
	Stopping
	Errored
)

type Runner struct {
	options *options

	ctx      context.Context
	closeCtx func(error)

	clock *time.Ticker

	currentBoard *board.Board

	// isStoppedChan will be closed when the engine is stopped.
	isStoppedChan chan any

	state State

	loopError error

	mu sync.Mutex
}

func (g *Runner) BoardSize() (int, int) {
	return g.options.boardRows, g.options.boardColumns
}

func (g *Runner) ClockInterval() time.Duration {
	return g.options.clockInterval
}

func (g *Runner) AddListener(l chan<- *board.Board) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.options.listeners = append(g.options.listeners, l)
}

func (g *Runner) RemoveListener(l chan<- *board.Board) {
	g.mu.Lock()
	defer g.mu.Unlock()
	slices.DeleteFunc(g.options.listeners, func(cur chan<- *board.Board) bool {
		return cur == l
	})
}

// Start starts the Runner clock.
// This function is Idempotent so calling it after the engine started does nothing.
func (g *Runner) Start() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.start()
}

func (g *Runner) start() error {
	var err error
	if g.state != Stopped {
		return fmt.Errorf("cannot start Runner in state %s", g.state.String())
	}
	g.state = Starting
	defer func() { g.state = Running }()
	if g.clock == nil {
		g.clock = time.NewTicker(g.options.clockInterval)
	} else {
		g.clock.Reset(g.options.clockInterval)
	}
	g.isStoppedChan = make(chan any)
	g.ctx, g.closeCtx = context.WithCancelCause(g.options.baseCtx)
	if g.options.initBoard != nil {
		g.currentBoard = g.options.initBoard
		g.options.boardRows, g.options.boardColumns = g.currentBoard.Rows(), g.currentBoard.Columns()
	} else {
		g.currentBoard, err = board.New(g.options.boardRows, g.options.boardColumns)
	}
	if g.options.initFn != nil {
		for p := range g.currentBoard.Points() {
			err = errors.Join(err, g.currentBoard.Set(p, g.options.initFn(p)))
		}
	}
	if err != nil {
		return err
	}
	go g.loop()
	return nil
}

// Close closes the
func (g *Runner) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.close()
}

func (g *Runner) close() error {
	if g.state == Errored {

	}
	if g.state != Running {
		return fmt.Errorf("cannot close runner.Runner in state %s", g.state.String())
	}
	g.state = Stopping
	g.clock.Stop()
	g.closeCtx(nil)

	select {
	case <-g.isStoppedChan:
		g.state = Stopped
		return nil

	case <-time.After(g.options.closeGracePeriod):
		g.state = Errored
		return fmt.Errorf("the Runner took longer than %s to stop after calling Runner.Close", g.options.closeGracePeriod)
	}
}

func (g *Runner) loop() {
	defer close(g.isStoppedChan)

	var err error
	for {
		g.broadcastBoard()
		g.currentBoard, err = g.currentBoard.NextBoard()
		if err != nil {
			g.setLoopError(err)
		}

		select {
		case <-g.ctx.Done():
			return
		case <-g.clock.C:
		}
	}
}

func (g *Runner) broadcastBoard() {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, listener := range g.options.listeners {
		bClone := g.currentBoard.Clone()
		listener <- bClone
	}
}

func (g *Runner) setLoopError(err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.loopError = err
	g.state = Errored
}

func (g *Runner) LoopError() error {
	return g.loopError
}

func New(optionFns ...option) (*Runner, error) {
	gameOptions := defaultOptions // Initialize gameOptions to a shallow copy of defaultOptions
	for _, optionFn := range optionFns {
		optionFn(&gameOptions)
	}
	return &Runner{options: &gameOptions}, nil
}
