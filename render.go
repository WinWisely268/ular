// i know that renderer should not be mixed with the game data stats
package ular

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"math/rand"
)

const (
	snakeHeadChar = ' '
	snakeChar     = ' '
	foodChar      = '⬤'
	obstacleVChar = '█'
	obstacleHChar = '▄'
)

// render everything
func (g *game) render() {
	s := g.screen
	s.Clear()

	g.renderSnake()
	g.renderFoods()
	g.renderBorder()
	g.renderScore()
	g.renderKeyStatusBar()

	if g.state.isPaused() {
		g.renderPause()
	}
	if g.state.isGameOver() {
		g.renderGameOver()
	}
	s.Show()
}

// renderSnake: renders the snake
func (g *game) renderSnake() {
	s := g.screen
	headStyle := tcell.StyleDefault.Background(tcell.ColorLightCyan)
	style := tcell.StyleDefault.Background(tcell.ColorBlue)

	s.SetContent(g.snake.head.x, g.snake.head.y, snakeHeadChar, nil, headStyle)
	for _, b := range g.snake.body {
		s.SetContent(b.x, b.y, snakeChar, nil, style)
	}
}

// renderFoods: renders all the food
func (g *game) renderFoods() {
	for _, f := range g.foods {
		g.screen.SetContent(f.loc.x, f.loc.y, foodChar, nil, blinkRandFg())
	}
}

// blinkRandFg: creates a blinking colors like it's christmas
func blinkRandFg() tcell.Style {
	fgs := []tcell.Color{
		tcell.ColorRed,
		tcell.ColorYellow,
		tcell.ColorGreen,
		tcell.ColorLemonChiffon,
	}
	return tcell.StyleDefault.Foreground(fgs[rand.Intn(len(fgs)-1)])
}

// renderScore: renders the score
func (g *game) renderScore() {
	helpStr := fmt.Sprintf("SCORE: %d", g.score)
	w := g.boardSize.w
	for i, c := range helpStr {
		g.screen.SetContent(
			w-len([]rune(helpStr))+i-2,
			0,
			c,
			nil,
			tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true),
		)
	}
}

func (g *game) renderKeyStatusBar() {
	keyStr := "?/p/P: show help/pause"
	for i, c := range keyStr {
		g.screen.SetContent(
			2+i,
			0,
			c,
			nil,
			tcell.StyleDefault.Foreground(tcell.ColorLightYellow),
		)
	}
}

// renderBorder: renders the border around the board
func (g *game) renderBorder() {
	w := g.boardSize.w
	h := g.boardSize.h
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue)
	renderBox(g.screen, 0, 1, w, h-1, style)
}

// renderDialog: renders the dialog popup
func (g *game) renderDialog(messages []string, messageStyle, boxStyle tcell.Style) {
	w := g.boardSize.w
	h := g.boardSize.h
	var maxLength int
	for _, m := range messages {
		length := len([]rune(m))
		if length > maxLength {
			maxLength = length
		}
	}
	boxLeft := w/2 - maxLength/2 - 1
	boxTop := h/2 - 1 - len(messages)
	renderBox(g.screen, boxLeft, boxTop, maxLength+4, len(messages)+2, boxStyle)
	for i, message := range messages {
		for j, c := range message {
			g.screen.SetContent(boxLeft+2+j, boxTop+1+i, c, nil, messageStyle)
		}
	}
}

// renderPause: renders the pause dialog
func (g *game) renderPause() {
	g.renderDialog([]string{
		statePaused.String(),
		"h/←  - left",
		"j/↓  - down",
		"k/↑  - up",
		"l/→  - right",
		"p/P    - pause",
		"r/R    - new game",
		"s/S    - resume",
		"q/Q    - quit",
		"esc    - quit",
	}, tcell.StyleDefault.Background(tcell.ColorBlack), tcell.StyleDefault.Foreground(tcell.ColorLightCyan).Background(tcell.ColorBlack))
}

// renderGameOver: renders the game over dialog
func (g *game) renderGameOver() {
	g.renderDialog([]string{
		stateLost.String(),
		"r/R to restart",
	}, tcell.StyleDefault.Background(tcell.ColorBlack), tcell.StyleDefault.Foreground(tcell.ColorRed))
}

// renderBox: common utility function to create any box we need (dialogs, boards, etc)
func renderBox(s tcell.Screen, startX, startY, width, height int, style tcell.Style) {
	s.SetContent(startX, startY, '┏', nil, style)
	s.SetContent(startX, startY+height-1, '┗', nil, style)
	s.SetContent(startX+width-1, startY, '┓', nil, style)
	s.SetContent(startX+width-1, startY+height-1, '┛', nil, style)
	for y := startY; y <= startY+height-1; y += height - 1 {
		for x := startX + 1; x < startX+width-1; x++ {
			s.SetContent(x, y, '━', nil, style)
		}
	}
	for x := startX; x <= startX+width-1; x += width - 1 {
		for y := startY + 1; y < startY+height-1; y++ {
			s.SetContent(x, y, '┃', nil, style)
		}
	}
}
