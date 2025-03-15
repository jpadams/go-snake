package main

import (
	"context"
	"dagger/go-snake/internal/dagger"
	"fmt"
	"time"
)

func New(
	snakeGameServer *dagger.Service,
) *GoSnake {
	return &GoSnake{
		Driver: dag.Container().
			From("alpine").
			WithExec([]string{"apk", "add", "curl"}).
			WithServiceBinding("snake", snakeGameServer),
	}
}

type GoSnake struct {
	// +internal-use-only
	Driver *dagger.Container
}

// Returns the state of the Snake game
// Full Grid Representation
//
// The board is represented as a 2D array where:
// 0 = empty space
// 1 = snake body
// 2 = food
//
// This helps an LLM "see" the full game layout like an image.
//
// Snake's Length and Exact Position
// Instead of just the head, the full body ([]point) is sent.
func (m *GoSnake) GetGameState(ctx context.Context) (string, error) {
	curlCmd := "curl -s http://snake:8080/state"
	state, err := m.Driver.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"sh", "-c", curlCmd}).Stdout(ctx)
	return state, err
}

// Make a move. Valid moves are "up", "down", "left", or "right"
func (m *GoSnake) MakeMove(ctx context.Context, move string) (string, error) {
	curlCmd := fmt.Sprintf(`curl -s -X POST http://snake:8080/move -d '{"move":"%s"}' -H "Content-Type: application/json"`, move)
	m.Driver.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"sh", "-c", curlCmd}).Sync(ctx)
	//return fmt.Sprintf("Moved %s, better GetGameState and make a move fast!", move)
	return m.GetGameState(ctx)
}

// Start the game!
func (m *GoSnake) StartGame(ctx context.Context) (string, error) {
	curlCmd := "curl -s http://snake:8080/start"
	m.Driver.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"sh", "-c", curlCmd}).Sync(ctx)
	//return "Game on! Watch out! Better GetGameState and make a move fast!"
	return m.GetGameState(ctx)
}
