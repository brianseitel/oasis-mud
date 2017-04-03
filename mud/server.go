package mud

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var gameServer Server

// Server object
type Server struct {
	connections []connection
	Up          bool
	Wizlock     bool
}

func (server *Server) handle(c *connection) {
	c.mob = login(c)

	err := registerConnection(c)
	if err != nil {
		return
	}

	interpret(c.mob, "look")
	for {
		if c.mob == nil || c.mob.client == nil {
			break
		}
		c.mob.statusBar()
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input = strings.Trim(input, "\r\n")

		if len(input) > 0 {
			interpret(c.mob, input)
		}

	}
}

// Serve boots up the server
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
			extractMob(c.mob, true)
			c.SendString(fmt.Sprintf("This user is already playing. Bye! %s", newline))
			c.end()

			gameServer.connections = append(gameServer.connections[:j], gameServer.connections[j+1:]...)

			return errors.New("this user is already playing")
		}
	}

	gameServer.connections = append(gameServer.connections, *c)
	return nil
}

func (server *Server) init() {
	bootDB()
}

func getConnections() []connection {
	return gameServer.connections
}
