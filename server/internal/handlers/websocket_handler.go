package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	internalWebSocket "github.com/quangtran666/kahoot-clone/internal/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketHandler struct {
	hub *internalWebSocket.Hub
}

func NewWebSocketHandler(hub *internalWebSocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	client := internalWebSocket.NewClient(conn, h.hub)
	h.hub.Register <- client

	go client.ReadMessages()
	go client.WriteMessages()
}
