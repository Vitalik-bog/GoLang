package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}
}

func (c *Client) WriteMessage(message []byte) {
	select {
	case c.Send <- message:
	default:
		log.Println("Client send buffer full")
		c.Close()
	}
}

func (c *Client) Close() {
	c.Conn.Close()
	close(c.Send)
}
