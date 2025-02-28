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

	gameService := service.NewGameService(hub)

	hub.RegisterEventHandler(event.SendMessage, gameService.HandleSendMessage)
	hub.RegisterEventHandler(event.UserConnected, gameService.HandleUserConnected)
	hub.RegisterEventHandler(event.UserDisconnected, gameService.HandleUserDisconnected)

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
