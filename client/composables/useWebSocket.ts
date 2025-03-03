import { EventType, type WebSoketMessage } from "~/types/websocket";

export const useWebSocket = (url: string) => {
    const socket = ref<WebSocket | null>(null);
    const isConnected = ref(false);
    const messages = ref<string[]>([]);
    const error = ref<string>('');

    const connect = () => {
        socket.value = new WebSocket(url);

        socket.value.onopen = () => {
            isConnected.value = true;

            error.value = ""
        }

        socket.value.onclose = () => {
            isConnected.value = false
        }

        socket.value.onmessage = (event: MessageEvent) => {
            const data: WebSoketMessage = JSON.parse(event.data);

            switch (data.type) {
                case EventType.SendMessage:
                    const msgPayload = data.payload as { username: string, message: string }
                    messages.value.push(`${msgPayload.username}: ${msgPayload.message}`);
                    break

                case EventType.UserConnected:
                    messages.value.push(`User connected: ${data.payload.username}`);
                    break

                case EventType.UserDisconnected:
                    messages.value.push(`User disconnected: ${data.payload.username}`);
                    break
                case EventType.RoomLeft:
                    messages.value.push(`${data.payload.username} has left the room`);
                    break
            }
        }
    }

    const sendMessage = (message: WebSoketMessage) => {
        if (socket.value?.readyState === WebSocket.OPEN) {
            socket.value.send(JSON.stringify(message));
        } else {
            error.value = "WebSocket is not connected";
        }
    }

    return {
        isConnected,
        messages,
        error,
        connect,
        sendMessage
    }
}