﻿package event

import "encoding/json"

type EventType string

const (
	SendMessage      EventType = "send_message"
	UserConnected    EventType = "user_connected"
	UserDisconnected EventType = "user_disconnected"

	CreateRoom EventType = "create_room"
	JoinRoom   EventType = "join_room"
	LeaveRoom  EventType = "leave_room"
	RoomLeft   EventType = "room_left"

	Error EventType = "error"
)

// IncomingEvent represents an event coming from clients
type IncomingEvent struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// OutgoingEvent represents an event going to clients
type OutgoingEvent struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

func CreateOutgoingEvent(eventType EventType, payload interface{}) ([]byte, error) {
	event := OutgoingEvent{
		Type:    eventType,
		Payload: payload,
	}

	return json.Marshal(event)
}

type SendMessagePayload struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

type CreateRoomPayload struct {
	RoomName string `json:"room_name"`
}

type JoinRoomPayload struct {
	RoomCode string `json:"room_code"`
	Username string `json:"username"`
}

type RoomJoinedPayload struct {
	Username string `json:"username"`
}

type RoomLeftPayload struct {
	Username string `json:"username"`
}

type RoomClosedPayload struct {
	RoomCode string `json:"room_code"`
}
