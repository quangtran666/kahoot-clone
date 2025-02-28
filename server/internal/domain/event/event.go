package event

import "encoding/json"

type EventType string

const (
	SendMessage      EventType = "send_message"
	UserConnected    EventType = "user_connected"
	UserDisconnected EventType = "user_disconnected"
)

// Event interface for all events
type Event interface {
	Type() EventType
}

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

type UserConnectedPayload struct {
	Username string `json:"username"`
}

type UserDisconnectedPayload struct {
	Username string `json:"username"`
}
