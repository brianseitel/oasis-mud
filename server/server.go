package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/brianseitel/mud/players"
)

var dbCommands *CommandDatabase
var dbItems *ItemDatabase
var dbRooms *RoomDatabase
var activePlayers []Player
var activeConnections []Connection

type Connection struct {
	conn   net.Conn
	buffer *bufio.ReadWriter
	player *Player
}

type Server struct {
	port             int
	log              *log.Logger
	connections      int
	totalConnections int
}

func (server *Server) Handle(c *Connection) {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for _ = range ticker.C {
			c.SendString("Tick!" + newline)
		}
	}()

	c.player = Login(c)
	activePlayers = append(activePlayers, *c.player)
	activeConnections = append(activeConnections, *c)
	room := dbRooms.FindRoom(c.player.Room)

	room.Display(*c)
	for {
		c.player.ShowStatusBar()
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input = strings.Trim(input, "\r\n")

		if len(input) > 0 {
			parts := strings.Split(input, " ")
			cmd := parts[0]
			line := strings.Join(parts[1:], " ")

			if cmd == "quit" {
				// Save character first
				command := SaveCommand{}
				command.Handle(c, line)

				// Say goodbye
				c.SendString("Seeya!" + newline)
				c.conn.Close()
				return
			}

			if command, ok := dbCommands.Lookup(cmd); ok {
				command.Handle(c, line)
			} else {
				c.SendString("I'm sorry. I don't know what you mean." + newline)
			}
		}
	}
}

func (server *Server) Logger() *log.Logger {
	return server.log
}

func (server *Server) Serve() {
	dbItems = NewItemDatabase()
	dbRooms = NewRoomDatabase()
	dbCommands = NewCommandDatabase()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.port))
	if err != nil {
		server.log.Printf("cannot start server: %s\n", err)
		os.Exit(1)
	}
	server.log.Printf("waiting for connections on %s\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			server.log.Printf("could not accept: %s\n", err)
		} else {
			server.log.Printf("connected: %s\n", conn.RemoteAddr())
			server.connections++
			server.totalConnections++
			go server.Handle(NewConnection(conn))
		}
	}
}

func NewConnection(connection net.Conn) *Connection {
	return &Connection{
		conn:   connection,
		buffer: bufio.NewReadWriter(bufio.NewReader(connection), bufio.NewWriter(connection)),
	}
}

func (c *Connection) SendString(text string) {
	c.conn.Write([]byte(text))
}

func (c *Connection) BufferData(text string) {
	c.buffer.Write([]byte(text))
}

func (c *Connection) SendBuffer() {
	c.buffer.Flush()
}

func (c *Connection) BroadcastToRoom(text string) {
	for _, connection := range activeConnections {
		if connection.player != c.player && c.player.Room == connection.player.Room {
			connection.SendString(c.player.Name + " says, \"" + text + "\"" + newline)
		}
	}
}
