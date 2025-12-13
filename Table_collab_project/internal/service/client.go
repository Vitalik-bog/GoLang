package service

import (
	"log"
	"sync"
	"time"

	"table_collab/internal/domain"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string
	RoomID   string
	Username string
	Color    string
	Conn     *websocket.Conn
	hub      *Hub
	send     chan domain.Event
	mu       sync.Mutex
	closed   bool
}

func NewClient(conn *websocket.Conn, hub *Hub, roomID string) *Client {
	return &Client{
		ID:     generateID(),
		RoomID: roomID,
		Conn:   conn,
		hub:    hub,
		send:   make(chan domain.Event, 256),
		Color:  generateColor(),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.Close()
	}()

	c.Conn.SetReadLimit(c.hub.config.WebSocket.MaxMessageSize)

	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(time.Duration(c.hub.config.WebSocket.PingPeriod) * time.Second))
		return nil
	})

	for {
		var event domain.Event
		err := c.Conn.ReadJSON(&event)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		c.handleEvent(event)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(time.Duration(c.hub.config.WebSocket.PingPeriod) * time.Second)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case event, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(event); err != nil {
				log.Printf("Write error: %v", err)
				return
			}

		case <-ticker.C:
			c.mu.Lock()
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()
		}
	}
}

func (c *Client) handleEvent(event domain.Event) {
	event.UserID = c.ID
	event.RoomID = c.RoomID
	event.Timestamp = time.Now().UnixMilli()

	switch event.Type {
	case domain.EventJoinRoom:
		if payload, ok := event.Payload.(map[string]interface{}); ok {
			if username, ok := payload["username"].(string); ok {
				c.Username = username
			}
		}
		c.hub.register <- c

	case domain.EventCursorMove:
		c.hub.broadcast <- event

	default:
		log.Printf("Unknown event type: %s", event.Type)
	}
}

func (c *Client) write(messageType int, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return c.Conn.WriteMessage(messageType, data)
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	c.closed = true
	close(c.send)
	c.Conn.Close()
}

func generateID() string {
	return "client_" + time.Now().Format("150405")
}

func generateColor() string {
	colors := []string{"#FF6B6B", "#4ECDC4", "#FFD166", "#06D6A0"}
	return colors[int(time.Now().Unix())%len(colors)]
}
