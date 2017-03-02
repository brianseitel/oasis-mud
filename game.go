package main

import (
	"github.com/brianseitel/oasis-mud/mud"
)

type Game struct {
	server  *mud.Server
	running bool
}

func NewGame(server *mud.Server) *Game {
	return &Game{
		server:  server,
		running: false,
	}
}

func (g *Game) Start() {
	g.server.Serve()
	g.running = true
}
