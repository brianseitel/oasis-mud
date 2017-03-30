package mud

import (
	"bufio"
	"net"
)

type connection struct {
	conn     net.Conn
	buffer   *bufio.ReadWriter
	mob      *mob
	snoopBy  *connection
	original *mob
}

func newConnection(c net.Conn) *connection {
	return &connection{
		conn:   c,
		buffer: bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c)),
	}
}

func (c *connection) end() {
	server := gameServer
	for j, con := range server.connections {
		if con.conn == c.conn {
			server.connections = append(server.connections[0:j], server.connections[j+1:]...)
			break
		}
	}

	c.conn.Close()
	gameServer = server
}

func (c *connection) SendString(text string) {
	c.conn.Write([]byte(text))
}

func (c *connection) BufferData(text string) {
	c.buffer.Write([]byte(text))
}

func (c *connection) SendBuffer() {
	c.buffer.Flush()
}
