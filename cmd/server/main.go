package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/Eyal-Shalev/bitmap-go"
	"github.com/Eyal-Shalev/game-of-life-go/board"
	"github.com/Eyal-Shalev/game-of-life-go/runner"
	"github.com/Eyal-Shalev/game-of-life-go/www"
	"github.com/dottedmag/must"
)

type Settings struct {
	Rows     int           `json:"rows"`
	Columns  int           `json:"columns"`
	Interval time.Duration `json:"interval"`
	Seed     uint64        `json:"seed,omitempty"`
}

func main() {
	defer slog.Info("Server Shutdown")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/game", gameHandler)
	mux.Handle("GET /", http.FileServerFS(www.FS))

	slog.Info("Listening on http://localhost:7676")
	err := http.ListenAndServe("127.0.0.1:7676", mux)
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("Error closed")
	} else if err != nil {
		slog.Error("Failed to start server", slog.Any("error", err))
	}
}

/*
00000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000101000000000000000000000000000110000001100000000000011000000000000000100010000110000000000001100001100000000100000100011000000000000000000110000000010001011000010100000000000000000000000001000001000000010000000000000000000000000010001000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000101000000000000000000000000000110000001100000000000011000000000000000100010000110000000000001100001100000000100000100011000000000000000000110000000010001011000010100000000000000000000000001000001000000010000000000000000000000000010001000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000000
000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000101000000000000000000000000000110000001100000000000011000000000000000100010000110000000000001100001100000000100000100011000000000000000000110000000010001011000010100000000000000000000000001000001000000010000000000000000000000000010001000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000000

'00000000',

0000000000000000000000001000000000000000
0000000000000000000000101000000000000000
0000000000001100000011000000000000110000
0000000000010001000011000000000000110000
1100000000100000100011000000000000000000
1100000000100010110000101000000000000000
0000000000100000100000001000000000000000
0000000000010001000000000000000000000000
0000000000001100000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000
0000000000000000000000000000000000000000

*/

func gameHandler(w http.ResponseWriter, r *http.Request) {
	initBoard, seed, err := parseInitFunc(r)
	if err != nil {
		slog.Warn("Failed to parse init function", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settings := Settings{Rows: initBoard.Rows(), Columns: initBoard.Columns(), Interval: 128 * time.Millisecond, Seed: seed}

	bChan := make(chan *board.Board)
	defer close(bChan)

	gameRunner := must.OK1(runner.New(
		runner.WithBaseContext(r.Context()),
		runner.WithBoardSize(settings.Rows, settings.Columns),
		runner.WithClockInterval(settings.Interval),
		runner.WithInitBoard(initBoard),
		runner.WithListener(bChan),
	))

	if err := gameRunner.Start(); err != nil {
		slog.Error("failed to start game runner", "error", err, "settings", settings)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer doIgnore(gameRunner.Close)

	settingsBytes := must.OK1(json.Marshal(settings))

	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	_, _ = fmt.Fprintf(w, "event: settings\n")
	_, _ = fmt.Fprintf(w, "data: %s\n\n", string(settingsBytes))

	for {
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		select {
		case b := <-bChan:
			boardText, err := b.Data().MarshalText()
			if err != nil {
				slog.Error("Failed to marshal board boardText", slog.Any("error", err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			_, _ = fmt.Fprintf(w, "data: %s\n\n", string(boardText))
		case <-r.Context().Done():
			return
		}
	}
}

const defaultRows = 10

func parseInitFunc(r *http.Request) (*board.Board, uint64, error) {
	var err error
	rowsStr := r.FormValue("rows")
	initStateStr := r.FormValue("init_state")
	if rowsStr == "" {
		rowsStr = strconv.Itoa(defaultRows)
	}
	rows, err := strconv.Atoi(rowsStr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse rows: %w", err)
	}

	if initStateStr == "" {
		var rng *rand.Rand
		var seed uint64
		rng, seed, err = parseSeed(r)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse init state: %w", err)
		}
		initBoard := must.OK1(board.New(rows, rows))
		for p := range initBoard.Points() {
			err = errors.Join(err, initBoard.Set(p, rng.Float32() > 0.5))
		}
		if err != nil {
			return nil, 0, fmt.Errorf("failed to init board: %w", err)
		}
		return initBoard, seed, nil
	}

	data := new(bitmap.BitMap)
	err = data.UnmarshalText([]byte(initStateStr))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse init state: %w", err)
	}

	initBoard, err := board.NewFromData(rows, data)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse init board: %w", err)
	}
	return initBoard, 0, nil
}

func parseSeed(r *http.Request) (*rand.Rand, uint64, error) {
	var err error
	seed := rand.Uint64()
	seedStr := r.FormValue("seed")
	if seedStr != "" {
		seed, err = strconv.ParseUint(seedStr, 10, 64)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("%q is not a valid base64 unsigned integer: %w", seedStr, err)
	}
	return rand.New(rand.NewPCG(seed, seed)), seed, nil
}

func doIgnore(fn func() error) {
	_ = fn()
}
