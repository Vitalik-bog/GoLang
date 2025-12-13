package ws

import (
	"log"
	"net/http"

	"table_collab/internal/service"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub *service.Hub
}

func NewHandler(hub *service.Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) ServeWebSocket(roomID string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := service.NewClient(conn, h.hub, roomID)

	go client.WritePump()
	go client.ReadPump()

	log.Printf("Client connected to room %s", roomID)
}
