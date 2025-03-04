package websocket

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/quangtran666/kahoot-clone/internal/domain/event"
)

var (
	// pongWait is the duration we wait before closing a connection due to ping timeouts.
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait. because it will send new ping before the pongWait
	pingInterval = (pongWait * 9) / 10
)

// Client represents a connected Websocket client.
type Client struct {
	connection *websocket.Conn
	Hub        *Hub
	Egress     chan []byte // Egress = Đầu ra
	UserId     string      // Somehow add username to client
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

	// Configure wait time for pong message,
	// This has to be set before the first initial ping
	if err := client.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error setting read deadline: %v", err)
		return
	}

	// Configure how to handle pong response
	client.connection.SetPongHandler(client.pongHandler)

	for {
		_, payload, err := client.connection.ReadMessage()
		if err != nil {
			// If connection is closed, we will RECEIVE AN ERROR here
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			} else if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Printf("Connection timed out: Client %s exceeds pong waits time of %v", client.UserId, pongWait)
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
	// Create a ticker that trigger a ping at given interval
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		client.Hub.Unregister <- client
	}()

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
		case <-ticker.C:
			// Send the ping
			if err := client.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("error writing ping message: %v", err)
				return // break the loop to close connection and clean up
			}
		}
	}
}

// pongHanlder use to handle pong message from client
func (client *Client) pongHandler(string) error {
	// Reset the read deadline
	return client.connection.SetReadDeadline(time.Now().Add(pongWait))
}
