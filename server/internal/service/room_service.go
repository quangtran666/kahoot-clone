package service

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/quangtran666/kahoot-clone/internal/domain/event"
	"github.com/quangtran666/kahoot-clone/internal/domain/room"
	"github.com/quangtran666/kahoot-clone/internal/websocket"
)

type RoomService interface {
	CreateRoom(event event.IncomingEvent, client *websocket.Client) error
	JoinRoom(event event.IncomingEvent, client *websocket.Client) error
	LeaveRoom(event event.IncomingEvent, client *websocket.Client) error
	HandleClientDisconnect(client *websocket.Client) error
	GetClientRoom(client *websocket.Client) (*room.Room, bool)
}

type roomService struct {
	hub   *websocket.Hub
	rooms map[string]*room.Room // Key: room code
	mutex sync.RWMutex
}

func NewRoomService(hub *websocket.Hub) RoomService {
	return &roomService{
		hub:   hub,
		rooms: make(map[string]*room.Room),
	}
}

func (r *roomService) CreateRoom(eventIncoming event.IncomingEvent, client *websocket.Client) error {
	var payload event.CreateRoomPayload
	if err := json.Unmarshal(eventIncoming.Payload, &payload); err != nil {
		log.Printf("Error unmarshalling create_room_event payload: %v", err)
		return err
	}

	newRoom := room.NewRoom(payload.RoomName, client.UserId)

	// Store the room
	r.mutex.Lock()
	r.rooms[newRoom.Code] = newRoom
	r.mutex.Unlock()

	// Add Client to the room
	newRoom.AddClient(client)

	log.Printf("Created room %v on %v", payload.RoomName, newRoom.Code)
	return nil
}

func (r *roomService) JoinRoom(eventIncoming event.IncomingEvent, client *websocket.Client) error {
	var payload event.JoinRoomPayload
	if err := json.Unmarshal(eventIncoming.Payload, &payload); err != nil {
		log.Printf("Error unmarshalling join_room_event payload:%v", err)
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// find the room
	room, exists := r.rooms[payload.RoomCode]
	if !exists {
		log.Printf("Room %v does not exist", payload.RoomCode)
		return errors.New("room not found")
	}

	// add client to room
	room.AddClient(client)

	// Notify other clients in the room
	if err := r.notifyRoom(room, event.JoinRoom, event.RoomJoinedPayload{
		Username: payload.Username,
	}); err != nil {
		log.Printf("Error notifying room: %v", err)
		return err
	}

	log.Printf("Client %v joined room %v", client.UserId, payload.RoomCode)
	return nil
}

func (r *roomService) LeaveRoom(eventIncoming event.IncomingEvent, client *websocket.Client) error {
	r.mutex.Lock()

	var clientRoom *room.Room
	var codeRoom string

	// find the room the client is in
	for code, room := range r.rooms {
		if room.HasClient(client) {
			clientRoom = room
			codeRoom = code
			break
		}
	}
	r.mutex.Unlock()

	if clientRoom == nil {
		log.Printf("Room %v does not exist while leaving room", codeRoom)
		return errors.New("client not in any room")
	}

	r.leaveRoom(clientRoom, client)
	return nil
}

func (r *roomService) HandleClientDisconnect(client *websocket.Client) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, room := range r.rooms {
		if room.HasClient(client) {
			log.Printf("Client %v disconnected from room %v", client.UserId, room.Code)
			r.leaveRoom(room, client)
			break
		}
	}

	return nil
}

func (r *roomService) GetClientRoom(client *websocket.Client) (*room.Room, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, room := range r.rooms {
		if room.HasClient(client) {
			return room, true
		}
	}

	return nil, false
}

func (r *roomService) leaveRoom(clientRoom *room.Room, client *websocket.Client) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Remove client from room
	room, exists := r.rooms[clientRoom.Code]
	if !exists {
		log.Printf("Room %v does not exist while leaving room", clientRoom.Code)
		return errors.New("room not found")
	}

	room.RemoveClient(client)

	// Delete room if no clients left
	if room.GetClientCount() == 0 {
		delete(r.rooms, clientRoom.Code)
	} else {
		// Or broadcast to other clients that a client has left
		if err := r.notifyRoom(room, event.RoomLeft, event.RoomLeftPayload{
			Username: client.UserId,
		}); err != nil {
			log.Printf("Error notifying room: %v", err)
			return err
		}
	}

	return nil
}

func (r *roomService) notifyRoom(room *room.Room, eventType event.EventType, payload interface{}) error {
	message, err := event.CreateOutgoingEvent(eventType, payload)
	if err != nil {
		log.Printf("Error creating outgoing event: %v", err)
		return err
	}

	room.BroadcastToClients(message)
	return nil
}
