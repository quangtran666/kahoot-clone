package main

import (
	"fmt"
	"github.com/gorilla/websocket"
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
	for {
		_, payload, err := client.connection.ReadMessage()
		if err != nil {
			break // Break the loop to close connection and clean up
		}

		fmt.Println(string(payload))
		client.egress <- payload
	}
}

func (client *Client) writeMessages() {
	for {
		message, ok := <-client.egress
		if !ok {
			break
		}
		for client := range client.manager.Clients {
			err := client.connection.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				break
			}
		}
	}
}
