package main

import (
	"flag"
	"fmt"
	"github.com/winwisely268/ular"
	"os"
)

const (
	defaultWidth  = 80
	defaultHeight = 40
)

var (
	version  string
	revision string
)

func main() {
	var w, h int
	var verCmd bool

	flag.IntVar(&w, "w", defaultWidth, "width of the game board")
	flag.IntVar(&h, "h", defaultHeight, "height of the game board")
	flag.BoolVar(&verCmd, "v", false, "show version information")
	flag.Parse()
	
	if verCmd {
		_, _ = fmt.Fprintf(os.Stdout, "Version: %s, Revision: %s\n", version, revision)
	} else {
		ular.Run(w, h)
	}
}