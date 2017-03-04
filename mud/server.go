package mud

import (
	"fmt"
	"net"
	"os"
	"strings"
	// "github.com/brianseitel/oasis-mud/helpers"
)

type Server struct {
	connections []Connection
}

func (server *Server) Handle(c *Connection) {
	c.player = Login(c)
	c.player.room = FindRoom(c.player.Room)
	newAction(c.player, c, "look")
	for {
		c.player.ShowStatusBar()
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input = strings.Trim(input, "\r\n")

		if len(input) > 0 {
			err := newActionWithInput(&action{player: c.player, conn: c, args: strings.Split(input, " ")})
			if err != nil {
				return // we're quitting
			}
		}

	}
}

func (server *Server) Serve(port int) {
	server.init()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("cannot start server: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("waiting for connections on %s\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("could not accept: %s\n", err)
		} else {
			fmt.Printf("connected: %s\n", conn.RemoteAddr())
			go server.Handle(NewConnection(conn))
		}
	}
}

func (server *Server) init() {
	Registry.items = NewItemDatabase()
	Registry.rooms = NewRoomDatabase()
}
