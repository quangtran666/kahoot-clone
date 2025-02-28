package main

import (
	"context"
	"github.com/joho/godotenv"
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
	context := context.Background()

	manager := NewManager(context)

	http.HandleFunc("/ws", manager.serveWebSocket)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
		return
	})

	err := http.ListenAndServe(os.Getenv("PORT"), nil)

	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
