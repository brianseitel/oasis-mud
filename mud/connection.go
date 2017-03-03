package mud

import (
	"bufio"
	"net"
)

type Connection struct {
	conn   net.Conn
	buffer *bufio.ReadWriter
	player *Player
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
	//  if player != c.player && c.player.Room == player.Room {
	//      SendString(c.player.Name + " says, \"" + text + "\"" + helpers.Newline)
	//  }
	// }
}
