package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

func (client *Client) readMessages() {
	defer func() {
		// Grateful close the connection and channel once this function done
		client.manager.removeClient(client)
	}()

	for {
		_, payload, err := client.connection.ReadMessage()
		if err != nil {
			// If connection is closed, we will RECEIVE AN ERROR here
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			log.Println("client disconnected")
			break // Break the loop to close connection and clean up
		}

		var event IncomingEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Printf("error unmarshalling event: %v", err)
			continue
		}

		// Routing the event to proper handler
		if err := client.manager.routeEventToHandler(event, client); err != nil {
			log.Printf("error handling event: %v", err)
		}
	}
}

func (client *Client) writeMessages() {
	defer func() {
		client.manager.removeClient(client)
	}()

	for {
		select {
		case message, ok := <-client.egress:
			// Oke will be failed if the channel is closed
			if !ok {
				// Todo: Manager has closed this channel, so communicate that to frontend
				if err := client.connection.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("error writing close message: %v", err)
				}
				// Return to close the goroutine
				return
			}

			if err := client.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("error writing message: %v", err)
				return
			}

			log.Printf("wrote message to client: %v", string(message))
		}
	}
}
