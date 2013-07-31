package broadcastd

import (
	"fmt"
	"io"
	"net/http"
)

type Connection struct {
	send chan string
}

func NewConnection() *Connection {
	send := make(chan string, 256)

	return &Connection{send}
}

func (c *Connection) Write(w io.Writer) {
	for msg := range c.send {
		fmt.Fprintf(w, "%s\n", msg)

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func (c *Connection) Close() {
	close(c.send)
}
