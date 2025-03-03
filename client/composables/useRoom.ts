import { EventType, type CreateRoomPayload, type JoinRoomPayload, type WebSoketMessage } from "~/types/websocket";

export const useRoom = (sendMessage: (msg: WebSoketMessage) => void) => {
    const currentRoom = ref<string>('');
    const roomName = ref<string>('');

    const createRoom = (name: string) => {
        const payload: CreateRoomPayload = {
            room_name: name
        }

        sendMessage({
            type: EventType.CreateRoom,
            payload
        });

        roomName.value = name;
    }

    const joinRoom = (roomCode: string, username: string) => {
        const payload: JoinRoomPayload = {
            room_code: roomCode,
            username
        }

        sendMessage({
            type: EventType.JoinRoom,
            payload
        });

        currentRoom.value = roomCode;
    }

    const leaveRoom = () => {
        if (currentRoom.value) {
            sendMessage({
                type: EventType.LeaveRoom,
                payload: {}
            })

            currentRoom.value = '';
            roomName.value = '';
        }
    }

    return {
        currentRoom,
        roomName,
        createRoom,
        joinRoom,
        leaveRoom
    }
}