package main

import (
	"log"
	"net/http"

	"github.com/IWasZ3r0Cool/GoTopia2025/internal/game"
	"github.com/IWasZ3r0Cool/GoTopia2025/internal/server"
)

func main() {
	log.Println("Starting GoTopia2025 Server...")

	// 1. Initialize Game Engine
	gameEngine := game.NewGameEngine()

	// 2. Initialize WebSocket Hub
	hub := server.NewHub(gameEngine)
	go hub.Run()

	// 3. Setup Routes
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(hub, w, r)
	})

	// Optional: Serve frontend static files if we build them into 'dist'
	// fs := http.FileServer(http.Dir("./frontend/dist"))
	// http.Handle("/", fs)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
