package main

type Game struct {
	server  *Server
	running bool
}

func NewGame(server *Server) *Game {
	return &Game{
		server:  server,
		running: false,
	}
}

func (g *Game) Start() {
	g.server.Serve()
	g.running = true
}
