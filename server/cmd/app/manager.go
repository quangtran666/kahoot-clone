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
	eventHandlers map[string]EventHandler
}

func NewManager(ctx context.Context) *Manager {
	return &Manager{
		Clients: make(map[*Client]bool),
		eventHandlers: map[string]EventHandler{
			SendMessageEvent:      handleSendMessageEvent,
			UserConnectedEvent:    handleUserConnectedEvent,
			UserDisconnectedEvent: handleUserDisconnectedEvent,
		},
	}
}

func (manager *Manager) broadcast(message []byte) {
	// Prevent concurrent writes or delete to the Clients map
	manager.RLock()
	defer manager.RUnlock()

	for client := range manager.Clients {
		client.egress <- message
	}
}

func (manager *Manager) broadcastExceptSelf(message []byte, client *Client) {
	manager.RLock()
	defer manager.RUnlock()

	for c := range manager.Clients {
		if c != client {
			c.egress <- message
		}
	}
}

func (manager *Manager) routeEventToHandler(event IncomingEvent, client *Client) error {
	if handler, ok := manager.eventHandlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			return err
		}
		return nil
	} else {
		return EventNotSupported
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

func (manager *Manager) removeClient(client *Client) {
	manager.Lock()
	defer manager.Unlock()

	if _, ok := manager.Clients[client]; ok {
		client.connection.Close()
		delete(manager.Clients, client)
		close(client.egress)
	}
}
