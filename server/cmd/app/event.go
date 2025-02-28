package main

import (
	"encoding/json"
	"errors"
	"log"
)

const (
	SendMessageEvent      = "send_message"
	UserConnectedEvent    = "user_connected"
	UserDisconnectedEvent = "user_disconnected"
)

var (
	EventNotSupported = errors.New("event not supported")
)

type IncomingEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type OutgoingEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func createOutgoingEvent(eventType string, payload interface{}) ([]byte, error) {
	event := OutgoingEvent{
		Type:    eventType,
		Payload: payload,
	}

	return json.Marshal(event)
}

type SendMessageEventPayload struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

type UserConnectedEventPayload struct {
	Username string `json:"username"`
}

type UserDisconnectedEventPayload struct {
	Username string `json:"username"`
}

type EventHandler func(event IncomingEvent, client *Client) error

func handleSendMessageEvent(event IncomingEvent, client *Client) error {
	var payload SendMessageEventPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Println("Error unmarshalling send_message_event payload:", err)
		return err
	}

	message, err := createOutgoingEvent(SendMessageEvent, payload)
	if err != nil {
		log.Printf("Error creating outgoing send_message_event: %v", err)
		return err
	}

	client.manager.broadcastExceptSelf(message, client)
	return nil
}

func handleUserConnectedEvent(event IncomingEvent, client *Client) error {
	return nil
}

func handleUserDisconnectedEvent(event IncomingEvent, client *Client) error {
	return nil
}
