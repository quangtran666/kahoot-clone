package service

import (
	"encoding/json"
	"log"

	eventDomain "github.com/quangtran666/kahoot-clone/internal/domain/event"
	"github.com/quangtran666/kahoot-clone/internal/websocket"
)

// GameService handle game-related business logic
type GameService interface {
	HandleSendMessage(event eventDomain.IncomingEvent, client *websocket.Client) error
	HandleUserConnected(client *websocket.Client) error
	HandleUserDisconnected(client *websocket.Client) error
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
		log.Printf("Error unmarshalling send_message_event payload: %v", err)
		return err
	}

	message, err := eventDomain.CreateOutgoingEvent(eventDomain.SendMessage, payload)
	if err != nil {
		log.Printf("Error creating outgoing event: %v", err)
		return err
	}

	s.hub.BroadcastMessageToAllExcept(message, client)
	return nil
}

func (s *gameService) HandleUserConnected(client *websocket.Client) error {
	payload := eventDomain.UserConnectedPayload{
		Username: client.UserId,
	}

	message, err := eventDomain.CreateOutgoingEvent(eventDomain.UserConnected, payload)
	if err != nil {
		log.Printf("Error creating outgoing event: %v", err)
		return err
	}

	s.hub.BroadcastMessageToAllExcept(message, client)
	return nil
}

func (s *gameService) HandleUserDisconnected(client *websocket.Client) error {
	payload := eventDomain.UserDisconnectedPayload{
		Username: client.UserId,
	}

	message, err := eventDomain.CreateOutgoingEvent(eventDomain.UserDisconnected, payload)
	if err != nil {
		log.Printf("Error creating outgoing event: %v", err)
		return err
	}

	s.hub.BroadcastMessageToAllExcept(message, client)
	return nil
}
