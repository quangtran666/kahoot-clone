package websocket

type BroadcastMessageData struct {
	Message []byte
	Client  *Client
}
