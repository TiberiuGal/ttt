package main

import (
	"tibi/lorem/evdispatch/ttt/server"
	"tibi/lorem/evdispatch/ttt"
)

func main() {

	c := make(chan ttt.Point)
	n := make(chan struct{})
	go server.Serve(c, n)
	ttt.GameLoop(c, n)
}
