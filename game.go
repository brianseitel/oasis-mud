package main

type Game struct {
	server         *Server
	itemDatabase   *ItemDatabase
	playerDatabase *PlayerDatabase
	roomDatabase   *RoomDatabase
	running        bool
}

func NewGame(server *Server) *Game {
	return &Game{
		server:         server,
		running:        false,
		itemDatabase:   NewItemDatabase(),
		playerDatabase: NewPlayerDatabase(),
		roomDatabase:   NewRoomDatabase(),
	}
}

func (g *Game) Start() {
	g.server.Serve()
	g.running = true
}
