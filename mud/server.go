package mud

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/brianseitel/oasis-mud/helpers"
)

var dbCommands *CommandDatabase
var dbRooms *RoomDatabase
var dbItems *ItemDatabase
var activePlayers []Player
var activeConnections []Connection

type Connection struct {
	conn   net.Conn
	buffer *bufio.ReadWriter
	player *Player
}

type Server struct {
	Port             int
	Log              *log.Logger
	Connections      int
	TotalConnections int
}

func (server *Server) Handle(c *Connection) {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for _ = range ticker.C {
			c.SendString("Tick!" + helpers.Newline)
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
				c.SendString("Seeya!" + helpers.Newline)
				c.conn.Close()
				return
			}

			if command, ok := dbCommands.Lookup(cmd); ok {
				command.Handle(c, line)
			} else {
				c.SendString("I'm sorry. I don't know what you mean." + helpers.Newline)
			}
		}
	}
}

func (server *Server) Logger() *log.Logger {
	return server.Log
}

func (server *Server) Serve() {
	dbItems = NewItemDatabase()
	dbRooms = NewRoomDatabase()
	dbCommands = NewCommandDatabase()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.Port))
	if err != nil {
		server.Log.Printf("cannot start server: %s\n", err)
		os.Exit(1)
	}
	server.Log.Printf("waiting for connections on %s\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			server.Log.Printf("could not accept: %s\n", err)
		} else {
			server.Log.Printf("connected: %s\n", conn.RemoteAddr())
			server.Connections++
			server.TotalConnections++
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
	// for _, connection := range activeConnections {
	// 	if player != c.player && c.player.Room == player.Room {
	// 		SendString(c.player.Name + " says, \"" + text + "\"" + helpers.Newline)
	// 	}
	// }
}
