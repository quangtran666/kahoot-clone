import { EventType, type CreateRoomPayload, type JoinRoomPayload, type WebSoketMessage } from "~/types/websocket";

export const useRoom = (sendMessage: (msg: WebSoketMessage) => void, currentRoom: Ref<string>, roomName: Ref<string>) => {
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
    
            roomName.value = '';
            currentRoom.value = '';
        } else {
            console.error("No room to leave");
        }
    }

    return {
        roomName,
        createRoom,
        joinRoom,
        leaveRoom
    }
}