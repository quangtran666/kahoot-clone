package main

import (
	"github.com/joho/godotenv"
	"github.com/quangtran666/kahoot-clone/internal/domain/event"
	"github.com/quangtran666/kahoot-clone/internal/handlers"
	"github.com/quangtran666/kahoot-clone/internal/service"
	"github.com/quangtran666/kahoot-clone/internal/websocket"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	hub := websocket.NewHub()

	go hub.Run()

	roomService := service.NewRoomService(hub)
	gameService := service.NewGameService(hub, roomService)

	hub.RegisterEventHandler(event.SendMessage, gameService.HandleSendMessage)
	hub.RegisterEventHandler(event.CreateRoom, roomService.CreateRoom)
	hub.RegisterEventHandler(event.JoinRoom, roomService.JoinRoom)
	hub.RegisterEventHandler(event.LeaveRoom, roomService.LeaveRoom)

	hub.RegisterService("game", gameService)
	hub.RegisterService("room", roomService)

	wsHandler := handlers.NewWebSocketHandler(hub)

	// Setup routes
	http.Handle("/ws", wsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Kahoot Clone API"))
	})

	if err := http.ListenAndServe(os.Getenv("PORT"), nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
