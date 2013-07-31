package broadcastd

type Pubsub struct {
	connections map[*Connection]bool
	broadcast   chan string
	register    chan *Connection
	unregister  chan *Connection
}

func NewPubsub() *Pubsub {
	return &Pubsub{
		connections: make(map[*Connection]bool),
		broadcast:   make(chan string),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
	}
}

func (p *Pubsub) Subscribe() *Connection {
	c := NewConnection()

	p.register <- c

	return c
}

func (p *Pubsub) Unsubscribe(c *Connection) {
	p.unregister <- c
}

func (p *Pubsub) Publish(msg string) {
	p.broadcast <- msg
}

func (p *Pubsub) Run() {
	for {
		select {
		case c := <-p.register:
			p.connections[c] = true
		case c := <-p.unregister:
			delete(p.connections, c)
			c.Close()
		case msg := <-p.broadcast:
			for c, _ := range p.connections {
				select {
				case c.send <- msg:
				default:
					delete(p.connections, c)
					c.Close()
				}
			}
		}
	}
}
