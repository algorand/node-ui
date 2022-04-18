package main

import (
	"flag"

	"github.com/algorand/node-ui/tui"
)

var uiPort uint64

func init() {
	flag.Uint64Var(&uiPort, "a", 0, "Port address to host TUI from, set to 0 to run directly")
}

func main() {
	flag.Parse()
	tui.Start(uiPort)
}
