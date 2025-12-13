package service

import (
	"log"
	"sync"
	"time"

	"table_collab/cmd/server/config"
	"table_collab/internal/domain"
	"table_collab/internal/storage/memory"
)

type Hub struct {
	rooms      *memory.RoomStore
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan domain.Event
	shutdown   chan struct{}
	mu         sync.RWMutex
	config     *config.Config
}

func NewHub(cfg *config.Config) *Hub {
	return &Hub{
		rooms:      memory.NewRoomStore(),
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan domain.Event, 1000),
		shutdown:   make(chan struct{}),
		config:     cfg,
	}
}

func (h *Hub) Run() {
	log.Println("Hub started")

	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)

		case client := <-h.unregister:
			h.handleUnregister(client)

		case event := <-h.broadcast:
			h.handleBroadcast(event)

		case <-h.shutdown:
			h.handleShutdown()
			return
		}
	}
}

func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	h.clients[client.ID] = client
	h.mu.Unlock()

	room, err := h.rooms.Get(client.RoomID)
	if err != nil {
		room = &domain.Room{
			ID:          client.RoomID,
			Name:        client.RoomID,
			Type:        domain.RoomTypeDocument,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
			MaxClients:  h.config.App.MaxClientsPerRoom,
			ClientCount: 1,
		}
		h.rooms.Save(room)
	} else {
		room.ClientCount++
		h.rooms.Save(room)
	}

	log.Printf("Client %s joined room %s", client.ID, client.RoomID)

	// Уведомляем других
	h.broadcast <- domain.Event{
		Type:      domain.EventJoinRoom,
		RoomID:    client.RoomID,
		UserID:    client.ID,
		Timestamp: time.Now().UnixMilli(),
		Payload: domain.JoinRoomPayload{
			Username: client.Username,
			Color:    client.Color,
		},
	}
}

func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	delete(h.clients, client.ID)
	h.mu.Unlock()

	if room, err := h.rooms.Get(client.RoomID); err == nil {
		room.ClientCount--
		if room.ClientCount < 0 {
			room.ClientCount = 0
		}
		h.rooms.Save(room)
	}

	log.Printf("Client %s left room %s", client.ID, client.RoomID)

	h.broadcast <- domain.Event{
		Type:      domain.EventLeaveRoom,
		RoomID:    client.RoomID,
		UserID:    client.ID,
		Timestamp: time.Now().UnixMilli(),
	}
}

func (h *Hub) handleBroadcast(event domain.Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		if client.RoomID == event.RoomID && client.ID != event.UserID {
			select {
			case client.send <- event:
			default:
				go client.Close()
			}
		}
	}
}

func (h *Hub) handleShutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range h.clients {
		client.Close()
	}

	close(h.register)
	close(h.unregister)
	close(h.broadcast)

	log.Println("Hub stopped")
}

func (h *Hub) Stop() {
	close(h.shutdown)
}
