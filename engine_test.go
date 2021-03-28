package ular

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	g *game
	r *require.Assertions
)

func TestGameStates(t *testing.T) {
	r = require.New(t)
	g = newGame(nil, 100, 100, 100*time.Millisecond)
	t.Run("test game creation", testGameCreation)
	t.Run("test move snake", testMoveSnake)
	t.Run("test snake eats food", testHandleAteFood)
	t.Run("test snake out of bounds", testSnakeOutOfBounds)
	t.Run("test snake eats itself", testHandleAteItself)
}

func testGameCreation(_ *testing.T) {
	r.Equal(g.boardSize.h, 100)
	r.Equal(g.boardSize.w, 100)
	r.Equal(g.state, stateStarted)
}

func testMoveSnake(_ *testing.T) {
	g.move(north)
	r.Equal(-1, g.dir.y)
	r.Equal(0, g.dir.x)

	// should return the same because you can't move snake head on its tail
	g.move(south)
	r.Equal(-1, g.dir.y)
	r.Equal(0, g.dir.x)

	g.move(east)
	r.Equal(1, g.dir.x)
	r.Equal(0, g.dir.y)

	// should return the same because you can't move snake head on its tail
	g.move(west)
	r.Equal(1, g.dir.x)
	r.Equal(0, g.dir.y)

	// should move south
	g.move(south)
	r.Equal(1, g.dir.y)
	r.Equal(0, g.dir.x)

	// should move west
	g.move(west)
	r.Equal(-1, g.dir.x)
	r.Equal(0, g.dir.y)
}

func testHandleAteFood(t *testing.T) {
	f := g.foods[0]
	g.snake.head = f.loc
	g.handleAteFood(f.loc)
	r.Equal(1, g.score)
	r.NotEqual(g.foods[0], f)
	r.Equal(5, len(g.snake.body))
}

func testHandleAteItself(t *testing.T) {
	g.snake.head = g.snake.body[1]
	g.handleAteItself(func() {
	})
	r.Equal(g.state, stateLost)
}

func testSnakeOutOfBounds(t *testing.T) {
	g.snake.head = &point{0, 2}
	g.handleSnakeOutOfBounds()
	r.Equal(&point{98, 2}, g.snake.head)
	
	g.snake.head = &point{2, 0}
	g.handleSnakeOutOfBounds()
	r.Equal(&point{1, 98}, g.snake.head)
	
	g.snake.head = &point{99, 3}
	g.handleSnakeOutOfBounds()
	r.Equal(&point{98, 3}, g.snake.head)
}