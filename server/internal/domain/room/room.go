package room

import (
	"github.com/google/uuid"
	"github.com/quangtran666/kahoot-clone/internal/websocket"
	"math/rand"
	"sync"
	"time"
)

const (
	RoomCodeLength  = 6
	RoomCodeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Room represents a room
type Room struct {
	ID        string
	Code      string
	Name      string
	CreatedAt time.Time
	OwnerID   string
	Clients   map[*websocket.Client]bool
	mutex     sync.RWMutex
}

func NewRoom(name, ownerID string) *Room {
	return &Room{
		ID:        uuid.New().String(),
		Code:      GenerateRoomCode(),
		Name:      name,
		CreatedAt: time.Now(),
		OwnerID:   ownerID,
		Clients:   make(map[*websocket.Client]bool),
	}
}

func GenerateRoomCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, RoomCodeLength)
	charsetLength := len(RoomCodeCharset)

	for i := range code {
		code[i] = RoomCodeCharset[rand.Intn(charsetLength)]
	}

	return string(code)
}

func (r *Room) AddClient(client *websocket.Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.Clients[client] = true
}

func (r *Room) RemoveClient(client *websocket.Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.Clients, client)
}

func (r *Room) HasClient(client *websocket.Client) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.Clients[client]
	return exists
}

func (r *Room) GetClientCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.Clients)
}

func (r *Room) BroadcastToClients(message []byte) {
	r.mutex.RLock()
	clients := make([]*websocket.Client, 0, len(r.Clients))
	for client := range r.Clients {
		clients = append(clients, client)
	}
	r.mutex.RUnlock()

	// Broadcast outside to avoid holding the lock, which can cause deadlocks
	for _, client := range clients {
		client.Egress <- message
	}
}

func (r *Room) BroadcastToClientsExcept(message []byte, except *websocket.Client) {
	r.mutex.RLock()
	clients := make([]*websocket.Client, 0, len(r.Clients))
	for client := range r.Clients {
		if client != except {
			clients = append(clients, client)
		}
	}
	r.mutex.RUnlock()

	// Broadcast outside to avoid holding the lock, which can cause deadlocks
	for _, client := range clients {
		client.Egress <- message
	}
}
