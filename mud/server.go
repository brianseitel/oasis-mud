package mud

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	// "github.com/brianseitel/oasis-mud/helpers"
)

type Server struct {
	connections []Connection
}

func (server *Server) Handle(c *Connection) {
	c.player = Login(c)
	c.player.isPlayer = true
	c.player.room = FindRoom(c.player.Room)
	pid++
	c.player.pid = pid
	newAction(c.player, c, "look")
	for {
		c.player.ShowStatusBar()
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input = strings.Trim(input, "\r\n")

		if len(input) > 0 {
			err := newActionWithInput(&action{mob: c.player, conn: c, args: strings.Split(input, " ")})
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
	defer listener.Close()

	fmt.Printf("waiting for connections on %s\n", listener.Addr())

	go server.timing()

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

func (server *Server) timing() {
	const (
		tickLen time.Duration = 5
	)

	pulse := time.NewTicker(time.Second)
	tick := time.NewTicker(time.Second * tickLen)

	for {
		select {
		case <-pulse.C:
			fmt.Printf("")
			break
		case <-tick.C:
			for _, r := range Registry.rooms {
				for _, m := range r.Mobs {
					m.wander()
				}
			}
			break
		}
	}
}

func (server *Server) init() {
	Registry.items = NewItemDatabase()
	Registry.mobs = NewMobDatabase()
	Registry.rooms = NewRoomDatabase()
}
