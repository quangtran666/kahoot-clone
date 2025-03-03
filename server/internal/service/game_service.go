package service

import (
	"encoding/json"
	"log"

	"github.com/quangtran666/kahoot-clone/internal/domain/event"

	"github.com/quangtran666/kahoot-clone/internal/websocket"
)

// GameService handle game-related business logic
type GameService interface {
	HandleSendMessage(event event.IncomingEvent, client *websocket.Client) error
}

type gameService struct {
	hub         *websocket.Hub
	roomService RoomService
}

func NewGameService(hub *websocket.Hub, roomService RoomService) GameService {
	return &gameService{
		hub:         hub,
		roomService: roomService,
	}
}

func (s *gameService) HandleSendMessage(eventIncoming event.IncomingEvent, client *websocket.Client) error {
	var payload event.SendMessagePayload
	if err := json.Unmarshal(eventIncoming.Payload, &payload); err != nil {
		log.Printf("Error unmarshalling send_message_event payload: %v", err)
		return err
	}

	message, err := event.CreateOutgoingEvent(event.SendMessage, payload)
	if err != nil {
		log.Printf("Error creating outgoing event: %v", err)
		return err
	}

	// Only send messages to room participants - this simplifies the logic
	room, inRoom := s.roomService.GetClientRoom(client)
	if inRoom {
		room.BroadcastToClientsExcept(message, client)
	} else {
		// Only send back to this client if not in a room
		log.Println("User not in a room")
	}

	return nil
}
