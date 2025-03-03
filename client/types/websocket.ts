export enum EventType {
    SendMessage = 'send_message',
    UserConnected = 'user_connected',
    UserDisconnected = 'user_disconnected',
    CreateRoom = 'create_room',
    JoinRoom = 'join_room',
    LeaveRoom = 'leave_room',
    RoomLeft = 'room_left',
}

export interface SendMessagePaylaod {
    message: string,
    username: string
}

export interface CreateRoomPayload {
    room_name: string
}

export interface JoinRoomPayload {
    room_code: string
    username: string
}

export interface RoomJoinedPayload {
    username: string
}

export interface RoomLeftPayload {
    username: string
}

export interface WebSoketMessage<T = any> {
    type: EventType
    payload: T
}