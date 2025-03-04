<template>
  <div class="container">
    <UCard>
      <template #header>
        <div v-if="!currentRoom" class="">
          <h2 class="text-xl font-bold mb-4">Create or Join Room</h2>
          <UForm :state="roomState" @submit="handleRoomAction">
            <UFormGroup label="Room Name" name="roomName" v-if="isCreating">
              <UInput v-model="roomState.roomName" placeholder="Enter room name ..." />
            </UFormGroup>
            <UFormGroup label="Room Code" name="roomCode" v-else>
              <UInput v-model="roomState.roomCode" placeholder="Enter room code ..." />
            </UFormGroup>
            <UFormGroup label="Username" name="username">
              <UInput v-model="roomState.username" placeholder="Enter username ..." />
            </UFormGroup>
            <div class="flex gap-2 mt-2">
              <UButton type="submit" @click="isCreating = true" :disabled="!roomState.username">
                Create Room
              </UButton>
              <UButton type="submit" @click="isCreating = false" :disabled="!roomState.username">
                Join Room
              </UButton>
            </div>
          </UForm>
        </div>
        <div v-else>
          <div class="flex justify-between items-center">
            <h2 class="text-xl font-bold">Room: {{ currentRoom }}</h2>
            <UButton @click="handleLeaveRoom" color="red">Leave Room</UButton>
          </div>
        </div>
      </template>

      <div v-if="error" class="text-red-500 mb-4">{{ error }}</div>

      <div v-if="currentRoom" class="space-y-4">
        <div class="h-64 overflow-y-auto p-4 border rounded">
          <div v-for="(message, i) in messages" :key="i">
            {{  message }}
          </div>
        </div>
      </div>

      <UForm :state="messageState" @submit="handleSendMessage" class="space-y-2">
          <UFormGroup label="Message">
            <UInput v-model="messageState.message" placeholder="Type a message..." />
          </UFormGroup>
          <UButton type="submit" block>Send</UButton>
        </UForm>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import { EventType } from '~/types/websocket';

const currentRoom = ref('')
const { isConnected, messages, error, connect, sendMessage } = useWebSocket("ws://localhost:8080/ws", currentRoom)
const { roomName, createRoom, joinRoom, leaveRoom } = useRoom(sendMessage, currentRoom)

const isCreating = ref(false)
const roomState = ref({
  roomName: '',
  roomCode: '',
  username: '',
})

const messageState = ref({
  message: '',
})

onMounted(() => {
  connect()
})

const handleRoomAction = () => {
  if (isCreating.value) {
    createRoom(roomState.value.roomName)
  } else {
    joinRoom(roomState.value.roomCode, roomState.value.username)
  }
}

const handleLeaveRoom = () => {
  leaveRoom()
}

const handleSendMessage = () => {
  sendMessage({
    type: EventType.SendMessage,
    payload: {
      message: messageState.value.message,
      username: roomState.value.username,
    }
  })
  messageState.value.message = ''
}

</script>