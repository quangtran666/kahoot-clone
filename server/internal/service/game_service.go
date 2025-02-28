package service

import (
	"encoding/json"
	eventDomain "github.com/quangtran666/kahoot-clone/internal/domain/event"
	"github.com/quangtran666/kahoot-clone/internal/websocket"
	"log"
)

// GameService handle game-related business logic
type GameService interface {
	HandleSendMessage(event eventDomain.IncomingEvent, client *websocket.Client) error
	HandleUserConnected(event eventDomain.IncomingEvent, client *websocket.Client) error
	HandleUserDisconnected(event eventDomain.IncomingEvent, client *websocket.Client) error
}

type gameService struct {
	hub *websocket.Hub
}

func NewGameService(hub *websocket.Hub) GameService {
	return &gameService{
		hub: hub,
	}
}

func (s *gameService) HandleSendMessage(event eventDomain.IncomingEvent, client *websocket.Client) error {
	var payload eventDomain.SendMessagePayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Fatal("Error unmarshalling send_message_event payload:", err)
		return err
	}

	message, err := eventDomain.CreateOutgoingEvent(eventDomain.SendMessage, payload)
	if err != nil {
		log.Fatal("Error creating outgoing event:", err)
		return err
	}

	s.hub.BroadcastMessageToAllExcept(message, client)
	return nil
}

func (s *gameService) HandleUserConnected(event eventDomain.IncomingEvent, client *websocket.Client) error {
	return nil
}

func (s *gameService) HandleUserDisconnected(event eventDomain.IncomingEvent, client *websocket.Client) error {
	return nil
}
