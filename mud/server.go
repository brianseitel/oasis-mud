package mud

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/brianseitel/oasis-mud/helpers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

var gameServer Server

type Server struct {
	connections []connection
}

func (server *Server) handle(c *connection) {
	c.mob = login(c)

	err := registerConnection(c)
	if err != nil {
		return
	}

	newAction(c.mob, c, "look")
	for {
		c.mob.statusBar()
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input = strings.Trim(input, "\r\n")

		if len(input) > 0 {
			err := newActionWithInput(&action{mob: c.mob, conn: c, args: strings.Split(input, " ")})
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
			go server.handle(newConnection(conn))
		}
	}
}

func (server *Server) timing() {
	pulse := time.NewTicker(time.Second / pulsePerSecond)
	tick := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-pulse.C:
			fmt.Printf(".")
			updateHandler()
			break
		case <-tick.C:
			fmt.Printf("o")
			break
		}
	}
}

func registerConnection(c *connection) error {
	for j, oc := range gameServer.connections {
		if c.mob.Name == oc.mob.Name {
			extractChar(c.mob)
			c.SendString(fmt.Sprintf("This user is already playing. Bye! %s", helpers.Newline))
			c.end()

			gameServer.connections = append(gameServer.connections[:j], gameServer.connections[j+1:]...)

			return errors.New("this user is already playing")
		}
	}

	gameServer.connections = append(gameServer.connections, *c)
	return nil
}

var (
	db *gorm.DB
)

func (server *Server) init() {
	bootDB()
}

func getConnections() []connection {
	return gameServer.connections
}
