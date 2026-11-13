package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	ws "mapup-geofencing/internal/websocket"
)

type WebSocketHandler struct {
	hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebSocketHandler) HandleAlerts(
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	h.hub.AddClient(conn)

	defer h.hub.RemoveClient(conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}