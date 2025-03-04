import { EventType, type RoomCreatedPayload, type RoomJoinedPayload, type WebSoketMessage } from "~/types/websocket";

export const useWebSocket = (url: string, currentRoomCode: Ref<string>, currentRoomName: Ref<string>) => {
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
                case EventType.RoomCreated:
                    const roomCreatedPayload = data.payload as RoomCreatedPayload
                    currentRoomCode.value = roomCreatedPayload.room_code;
                    messages.value.push(`Room created successfully! Room Code: ${roomCreatedPayload.room_code}`)
                    break
                case EventType.RoomJoin:
                    const roomJoinedPayload = data.payload as RoomJoinedPayload
                    console.log(roomJoinedPayload);
                    currentRoomName.value = roomJoinedPayload.room_name;
                    messages.value.push(`${roomJoinedPayload.username} has joined the room ${roomJoinedPayload.room_name}`);
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