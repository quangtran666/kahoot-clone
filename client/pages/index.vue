<template>
  <div class="">
    <section>
      <UCard>
        <template #header>
          <h1 class="text-2xl font-bold text-center">Current room: {{ room }}</h1>
          <UForm :state="roomState" :schema="roomSchema" @submit="joinRoom" class="space-y-4">
            <UFormGroup label="Enter a room id" name="room_id">
              <UInput placeholder="Enter a room id..." v-model="roomState.room_id" />
            </UFormGroup>
            <UButton type="submit" block>
              Join
            </UButton>
          </UForm>
        </template>
  
        <p class="text-xl font-bold">Conversation</p>
        <UTextarea color="white" disabled v-model="messageText" />
        
        <template #footer>
          <UForm :schema="schema" :state="state" @submit="sendMessage" class="space-y-2">
            <UFormGroup label="Username" name="username">
              <UInput placeholder="Enter a username..." v-model="state.username" icon="solar:user-bold"/>
            </UFormGroup>
              
            <UFormGroup label="Message" name="message">
              <UInput placeholder="Enter a message..." v-model="state.message" icon="solar:chat-round-bold"/>
            </UFormGroup>
            <UButton type="submit" block>
              Send
            </UButton>
          </UForm>
        </template>
      </UCard>
    </section>
  </div>
</template>

<script setup lang="ts">
import type {FormSubmitEvent} from "#ui/types";
import {z} from "zod";

const room = ref(0)

const roomSchema = z.object({
  room_id: z.string().min(1, "Room id is required"),
})

const roomState = ref({
  room_id: ""
})

type RoomSchema = z.output<typeof roomSchema>;

const joinRoom = async (event: FormSubmitEvent<RoomSchema>) => {
  console.log(event);
}

const schema = z.object({
  message: z.string().min(1, "Message is required"),
  username: z.string().min(1, "Username is required"),
})

type Schema = z.output<typeof schema>;

const state = ref({
  message: "", 
  username: "",
})

interface Event {
  type: string;
  payload: ChatEvent;
}

interface ChatEvent {
  username: string;
  message: string;
}

const sendMessage = async (event: FormSubmitEvent<Schema>) => {
  const eventMessage: Event = {
    type: "send_message",
    payload: {
      username: event.data.username,
      message: event.data.message,
    }
  }
  
  webSocket.value?.send(JSON.stringify(eventMessage));
}

const webSocket = ref<WebSocket | null>(null);
const messages = ref<string[]>([]);
const messageText = computed(() => messages.value.join("\n"));

onMounted(() => {
  webSocket.value = new WebSocket("ws://localhost:8080/ws");
  
  webSocket.value.onopen = () => {
    console.log("WebSocket connection established");
  }
  
  webSocket.value.onmessage = (event: MessageEvent) => {
    messages.value.push(event.data);
  }
})

</script>