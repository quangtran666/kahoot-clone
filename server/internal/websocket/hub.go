package websocket

import (
	"errors"
	"github.com/quangtran666/kahoot-clone/internal/domain/event"
	"log"
	"sync"
)

// Hub only manages active Websocket connections (not game/room state)
type Hub struct {
	Clients             map[*Client]bool
	Register            chan *Client
	Unregister          chan *Client
	Broadcast           chan []byte
	BroadcastExceptSelf chan BroadcastMessageData
	sync.RWMutex
	EventHandlers map[event.EventType]EventHandler
	Services      map[string]interface{}
}

// EventHandler defines a signature for event handlers.
type EventHandler func(event event.IncomingEvent, client *Client) error

func NewHub() *Hub {
	return &Hub{
		Clients:             make(map[*Client]bool),
		Register:            make(chan *Client, 100),
		Unregister:          make(chan *Client, 100),
		Broadcast:           make(chan []byte, 100),
		BroadcastExceptSelf: make(chan BroadcastMessageData, 100),
		EventHandlers:       make(map[event.EventType]EventHandler),
		Services:            make(map[string]interface{}),
	}
}

func (hub *Hub) RegisterEventHandler(eventType event.EventType, handler EventHandler) {
	hub.EventHandlers[eventType] = handler
}

func (hub *Hub) RegisterService(name string, service interface{}) {
	hub.Services[name] = service
}

func (hub *Hub) routeEventToHandler(event event.IncomingEvent, client *Client) error {
	hub.RLock()
	defer hub.RUnlock()

	if handler, ok := hub.EventHandlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			log.Printf("error handling event: %v", err)
			return err
		}
		return nil
	} else {
		log.Printf("no handler for event type: %s", event.Type)
		return errors.New("no handler for event")
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.RegisterClient(client)
		case client := <-hub.Unregister:
			hub.UnregisterClient(client)
		case message := <-hub.Broadcast:
			hub.BroadcastMessage(message)
		case message := <-hub.BroadcastExceptSelf:
			hub.BroadcastMessageExceptSelf(message.Message, message.Client)
		}
	}
}

func (hub *Hub) RegisterClient(client *Client) {
	hub.Lock()
	defer hub.Unlock()

	hub.Clients[client] = true
	log.Printf("client %v registered", client.UserId)
}

func (hub *Hub) UnregisterClient(client *Client) {
	hub.Lock()
	defer hub.Unlock()

	if _, ok := hub.Clients[client]; ok {
		// Let room service know that client is disconnecting
		if roomService, ok := hub.Services["room"]; ok {
			if rs, ok := roomService.(interface{ HandleClientDisconnect(*Client) error }); ok {
				rs.HandleClientDisconnect(client)
			}
		}

		delete(hub.Clients, client)
		close(client.Egress)
		client.connection.Close()
	} else {
		log.Printf("Client not found in hub")
	}
}

func (hub *Hub) BroadcastMessage(message []byte) {
	hub.RLock()
	defer hub.RUnlock()

	for client := range hub.Clients {
		client.Egress <- message
	}
}

func (hub *Hub) BroadcastMessageExceptSelf(message []byte, client *Client) {
	hub.RLock()
	defer hub.RUnlock()

	for c := range hub.Clients {
		if c != client {
			c.Egress <- message
		}
	}
}

func (hub *Hub) BroadcastMessageToAllExcept(message []byte, client *Client) {
	hub.BroadcastExceptSelf <- BroadcastMessageData{
		Message: message,
		Client:  client,
	}
}
