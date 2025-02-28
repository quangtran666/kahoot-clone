package main

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	websocketUpgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	Clients map[*Client]bool
	sync.RWMutex
}

func NewManager(ctx context.Context) *Manager {
	return &Manager{
		Clients: make(map[*Client]bool),
	}
}

func (manager *Manager) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, manager)
	manager.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (manager *Manager) addClient(client *Client) {
	manager.Lock()
	defer manager.Unlock()

	manager.Clients[client] = true
}
