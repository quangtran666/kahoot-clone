package websocket

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/quangtran666/kahoot-clone/internal/domain/event"
	"log"
)

// Client represents a connected Websocket client.
type Client struct {
	connection *websocket.Conn
	Hub        *Hub
	Egress     chan []byte
	UserId     string
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		connection: conn,
		Hub:        hub,
		Egress:     make(chan []byte),
		UserId:     uuid.New().String(),
	}
}

func (client *Client) ReadMessages() {
	defer func() {
		client.Hub.Unregister <- client
	}()

	for {
		_, payload, err := client.connection.ReadMessage()
		if err != nil {
			// If connection is closed, we will RECEIVE AN ERROR here
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close connection and clean up
		}

		var event event.IncomingEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Printf("error unmarshalling event: %v", err)
			continue
		}

		if err := client.Hub.routeEventToHandler(event, client); err != nil {
			log.Printf("error routing event: %v", err)
		}
	}
}

func (client *Client) WriteMessages() {
	for {
		select {
		case message, ok := <-client.Egress:
			// Oke will be failed if the channel is closed
			if !ok {
				// break the loop to close connection and clean up
				return
			}

			if err := client.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("error writing message: %v", err)
				return
			}
		}
	}
}
