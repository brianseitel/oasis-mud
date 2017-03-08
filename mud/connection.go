package mud

import (
	"bufio"
	"net"
)

type connection struct {
	conn   net.Conn
	buffer *bufio.ReadWriter
	player *player
}

func newConnection(c net.Conn) *connection {
	return &connection{
		conn:   c,
		buffer: bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c)),
	}
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
