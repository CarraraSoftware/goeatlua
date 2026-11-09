package main

import (
	"goeatlua/interp"
)


func main() {
	interp.Init()
	defer interp.Deinit()

	interp.REPL()
}
