package main

import (
	"path/filepath"

	"github.com/brianseitel/oasis-mud/mud"
)

type game struct {
	server  *mud.Server
	running bool
}

func newGame(server *mud.Server) *game {
	return &game{
		server:  server,
		running: false,
	}
}

func (g *game) Start() {
	server := mud.Server{}
	server.Serve(8099)
	server.BasePath, _ = filepath.Abs("")
}
