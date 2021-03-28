package ular

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"log"
	"os"
	"time"
)

func Run(w, h int, isResizeable bool) {
	encoding.Register()
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err = s.Init(); err != nil {
		log.Fatal(err)
	}

	if w == 0 && h == 0 {
		w, h = s.Size()
	}

	w, h = adjustWidthHeight(s, w, h)

	s.SetStyle(tcell.StyleDefault)
	s.Clear()
	tickerDuration := 80 * time.Millisecond

	g := newGame(s, w, h, tickerDuration)
	g.start()

	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			if isResizeable {
				w, h = s.Size()
			}
			w, h = adjustWidthHeight(s, w, h)
			g.onResize(s, w, h)

		case *tcell.EventKey:
			if g.state.isPaused() {
				if ev.Key() == tcell.KeyEscape {
					s.Fini()
					os.Exit(0)
				}
				if ev.Rune() == 's' || ev.Rune() == 'S' {
					g.togglePause()
				}
			}
			if g.state.isGameOver() {
				switch ev.Rune() {
				case 'r', 'R':
					g = newGame(s, w, h, tickerDuration)
					g.start()
				}
			}
			if g.state.isStarted() {
				switch ev.Key() {
				case tcell.KeyEscape:
					s.Fini()
					os.Exit(0)
				case tcell.KeyUp:
					g.move(north)
				case tcell.KeyDown:
					g.move(south)
				case tcell.KeyLeft:
					g.move(west)
				case tcell.KeyRight:
					g.move(east)
				}
				switch ev.Rune() {
				case 'p', 'P':
					g.togglePause()
				case 'q', 'Q':
					s.Fini()
					os.Exit(0)
				case 'h', 'H':
					g.move(west)
				case 'j', 'J':
					g.move(south)
				case 'k', 'K':
					g.move(north)
				case 'l', 'L':
					g.move(east)
				case '?':
					g.togglePause()
				}
			}
			g.render()
		}
	}
}

func adjustWidthHeight(s tcell.Screen, w, h int) (int, int) {
	screenW, screenH := s.Size()
	if w > screenW {
		w = screenW
	}
	if h > screenH {
		h = screenH
	}
	return w, h
}
