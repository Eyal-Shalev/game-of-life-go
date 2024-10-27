# Game of Life Go

## Overview

The Game of Life Go is a Go implementation of Conway's Game of Life, a cellular automaton devised by mathematician John Conway. This implementation includes a server that handles HTTP requests to initialize and run the game, utilizing various Go packages for functionality.

## Repository Structure

- `cmd/server/main.go`: The entry point for the server application.
- `board`: Contains logic and data structures for the game board.
- `runner`: Handles the execution of the game logic.
- `www`: Static files served by the server.
- `bitmap-go`: A dependency that helps with bitmap operations.

## Getting Started

### Prerequisites

- Go 1.19 or later

### Installation

1. Clone the repository:

```sh
git clone https://github.com/Eyal-Shalev/game-of-life-go.git
cd game-of-life-go
```

2. Install dependencies:

```sh
go mod tidy
```

### Running the Server

To start the server, run:

```sh
go run cmd/server/main.go
```

The server listens on `http://localhost:7676` by default.

## API Endpoints

### `/api/v1/game`

Handles starting a new game with optional parameters:

- `rows`: Number of rows for the game board.
- `init_state`: Initial state of the board in a specific format.
- `seed`: Seed for random number generation.

### `/`

Serves static files from the `www` directory.

## Code Structure

### `main.go`

The main file initializes the server, sets up routes, and starts listening for HTTP requests.

Key functions include:

- `main()`: Sets up the server and routes.
- `gameHandler()`: Handles game initialization and streaming game state updates.
- `parseInitFunc()`: Parses initialization parameters from the request.
- `parseSeed()`: Parses or generates a seed for random number generation.
- `doIgnore()`: Helper function to ignore errors from function calls.

### Board Initialization

The game board can be initialized in two ways:

1. Randomly, using a seed for reproducibility.
2. From a specified initial state.

### Game Runner

The game logic is managed by the `runner` package, which handles the game's execution and state updates.

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add some feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Open a pull request.

## License

This project is licensed under the MIT License.

## Contact

For questions or support, please open an issue on the [GitHub repository](https://github.com/Eyal-Shalev/game-of-life-go/issues).