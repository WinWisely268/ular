package ular

import (
	"github.com/gdamore/tcell/v2"
	"math/rand"
	"runtime"
	"time"
)

// ================================================================================================
// gameState
// ================================================================================================

const (
	statePaused gameState = iota
	stateLost
	stateStarted

	defaultSnakeLength = 4
)

// gameState
// there are 4 global "states" in this game:
// - paused
// - game over / lost
// - started
type gameState int

// gameState satisfies Stringer
func (gs gameState) String() string {
	switch gs {
	case statePaused:
		return "Game Paused"
	case stateLost:
		return "You Lost!"
	case stateStarted:
		return "Game Started"
	default:
		return "unknown"
	}
}

// is it paused?
func (gs gameState) isPaused() bool {
	return gs == statePaused
}

// is the game over?
func (gs gameState) isGameOver() bool {
	return gs == stateLost
}

// is it started / resumed?
func (gs gameState) isStarted() bool {
	return gs == stateStarted
}

// point represents coordinates
type point struct {
	x int
	y int
}

// ================================================================================================
// snake (on a plane)
// ================================================================================================

// snake contains parts created from array of signed integers.
type snake struct {
	head *point
	body []*point
}

func newSnake(x, y int) *snake {
	h := &point{x, y}
	var body []*point
	for i := 1; i <= defaultSnakeLength; i++ {
		if body == nil {
			body = []*point{}
		}
		body = append(body, &point{x - i, y})
	}
	return &snake{h, body}
}

// length of the snake
func (s *snake) length() int {
	return len(s.body)
}

// ================================================================================================
// snake food
// ================================================================================================

type food struct {
	loc *point
}

func newFood(w, h int) *food {
	return &food{
		loc: &point{
			x: rand.Intn(w - 2),
			y: rand.Intn(h - 4),
		},
	}
}

// ================================================================================================
// cardinal direction
// ================================================================================================
type direction int

const (
	north direction = iota
	east
	south
	west
)

func (d direction) getPoint() *point {
	switch d {
	case north:
		return &point{0, -1}
	case south:
		return &point{0, 1}
	case east:
		return &point{1, 0}
	case west:
		return &point{-1, 0}
	default:
		return nil
	}
}

// ================================================================================================
// board
// ================================================================================================
type board struct {
	w int
	h int
}

// ================================================================================================
// actual game data
// ================================================================================================
// game contains entities (snake, and foods)
type game struct {
	snake     *snake
	dir       *point
	foods     []*food
	screen    tcell.Screen
	boardSize board
	tick      time.Duration
	state     gameState
	chQuit    chan struct{}
	score     int
}

// toggle pause status
func (g *game) togglePause() {
	if g.state.isPaused() {
		g.state = stateStarted
	} else {
		g.state = statePaused
	}
}

// onResize: callback on resize screen
// it refreshes board size, regenerate foods,
func (g *game) onResize(s tcell.Screen, w, h int) {
	g.screen = s
	g.boardSize.w = w
	g.boardSize.h = h
	g.foods = []*food{
		newFood(g.boardSize.w, g.boardSize.h),
		newFood(g.boardSize.w, g.boardSize.h),
		newFood(g.boardSize.w, g.boardSize.h),
	}
	g.render()
}

// newGame creates a new game
func newGame(screen tcell.Screen, w, h int, tick time.Duration) *game {
	halfWidth := w / 2
	maxWidth := halfWidth + (halfWidth % 2) - 1

	g := &game{
		snake: newSnake(maxWidth, h/2),
		dir:   east.getPoint(), // starts facing east
		foods: []*food{
			newFood(w, h),
			newFood(w, h),
			newFood(w, h),
		},
		boardSize: board{
			w, h,
		},
		tick:   tick,
		state:  stateStarted,
		score:  0,
		screen: screen,
	}
	return g
}

// move: points the snake head to the chosen cardinal direction
func (g *game) move(d direction) {
	switch d {
	case north:
		if g.dir.y != 1 {
			g.dir = north.getPoint()
		}
	case south:
		if g.dir.y != -1 {
			g.dir = south.getPoint()
		}
	case east:
		if g.dir.x != -1 {
			g.dir = east.getPoint()
		}
	case west:
		if g.dir.x != 1 {
			g.dir = west.getPoint()
		}
	}
}

// onUpdateCallback is the function that is called on every tick.
func (g *game) onUpdateCallback() {
	if g.state.isPaused() || g.state.isGameOver() {
		return
	}
	head := g.snake.head
	var lastPart *point
	// update snake's body
	for i := g.snake.length() - 1; i >= 0; i-- {
		if lastPart == nil {
			lastPart = &point{g.snake.body[i].x, g.snake.body[i].y}
		}
		if i > 0 {
			g.snake.body[i].x = g.snake.body[i-1].x
			g.snake.body[i].y = g.snake.body[i-1].y
		} else {
			g.snake.body[i].x = head.x
			g.snake.body[i].y = head.y
		}
	}
	g.handleAteFood(lastPart)
	g.handleSnakeOutOfBounds()
	g.handleAteItself(func() {
		g.render()
		g.stop()
	})
	g.render()
}

// handleAteFood: increment game score, adds a new body part to the snake
// and generate more random food for the snake.
func (g *game) handleAteFood(lastPart *point) {
	head := g.snake.head
	for i, f := range g.foods {
		if head.x == f.loc.x && head.y == f.loc.y {
			g.score++
			g.snake.body = append(g.snake.body, lastPart)
			g.foods[i] = newFood(g.boardSize.w, g.boardSize.h)
		}
	}
}

// handleAteItself: Ouroboros, a snake cannot eat itself.
// it's the only condition where we will end the game (game over)
func (g *game) handleAteItself(cb func()) {
	head := g.snake.head
	for i := 1; i < g.snake.length()-1; i++ {
		if head.x == g.snake.body[i].x && head.y == g.snake.body[i].y {
			g.state = stateLost
		}
	}
}

// handleSnakeOutOfBounds: on hitting the maximum width and height of the current board,
// the snake will need to update its head (and body parts) to the opposite end of the board.
func (g *game) handleSnakeOutOfBounds() {
	newHead := g.snake.head

	newHead.x += g.dir.x
	newHead.y += g.dir.y
	w := g.boardSize.w
	h := g.boardSize.h

	switch {
	case newHead.x < 1:
		newHead.x = w - 2
	case newHead.x >= w-1:
		newHead.x = 1
	case newHead.y < 1:
		newHead.y = h - 2
	case newHead.y >= h-1:
		newHead.y = 1
	}
	// create new snake
	g.snake.head = newHead
}

// start: starts the main loop of the game
func (g *game) start() {
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		ticker := time.NewTicker(g.tick)
		for {
			select {
			case <-ticker.C:
				g.onUpdateCallback()
			case <-g.chQuit:
				ticker.Stop()
			}
		}
	}()
}

// stop: stops the game.
func (g *game) stop() {
	g.chQuit <- struct{}{}
}
