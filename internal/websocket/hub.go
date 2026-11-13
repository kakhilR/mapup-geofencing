package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) AddClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn] = true
}

func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, conn)
	conn.Close()
}

func (h *Hub) Broadcast(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if err := client.WriteMessage(
			websocket.TextMessage,
			message,
		); err != nil {
			// Ignore failed clients for now.
			// We'll handle cleanup more robustly later if needed.
		}
	}
}