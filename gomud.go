package main

import (
	"github.com/brianseitel/oasis-mud/mud"
)

func main() {
	server := &mud.Server{}
	game := newGame(server)
	game.Start()
}
