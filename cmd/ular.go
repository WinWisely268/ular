package main

import (
	"flag"
	"fmt"
	"github.com/winwisely268/ular"
	"os"
)

var (
	version  string
	revision string
)

func main() {
	var w, h int
	var verCmd bool
	var resizable bool

	flag.IntVar(&w, "w", 0, "width of the game board")
	flag.IntVar(&h, "h", 0, "height of the game board")
	flag.BoolVar(&verCmd, "v", false, "show version information")
	flag.BoolVar(&resizable, "r", true, "resize arena follow screen size")
	flag.Parse()

	if verCmd {
		_, _ = fmt.Fprintf(os.Stdout, "Version: %s, Revision: %s\n", version, revision)
	} else {
		ular.Run(w, h, resizable)
	}
}
