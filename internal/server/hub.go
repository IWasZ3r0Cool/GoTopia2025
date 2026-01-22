package server

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/IWasZ3r0Cool/GoTopia2025/internal/game"
)

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client

	GameEngine *game.GameEngine
	mutex      sync.Mutex
}

func NewHub(ge *game.GameEngine) *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		GameEngine: ge,
	}
}

func (h *Hub) Run() {
	// Start the game engine loop concurrently
	h.GameEngine.Start()

	// In a real app, we'd have a ticker here to broadcast state periodically or on change.
	// For now, we'll hook into the GameEngine.
	// The GameEngine creates ticks. We might need to poll it or have it push to us.
	// Simpler: Just rely on incoming messages or a ticker here to broadcast state.

	// Let's broadcast state on every meaningful action or periodically.
	// For 200ms ticks, maybe broadcast every tick?
	// To do this cleanly, let's just loop.

	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.Clients[client] = true
			h.mutex.Unlock()

			// Assign player ID (simple UUID/Timestamp or just use memory address for now... wait, Client needs ID)
			// h.GameEngine.AddPlayer(client.ID) -> Assumes Client checks Join message.
			log.Printf("Client registered: %s", client.conn.RemoteAddr())

		case client := <-h.Unregister:
			h.mutex.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
			h.mutex.Unlock()

		case message := <-h.Broadcast:
			h.mutex.Lock()
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
			h.mutex.Unlock()
		}
	}
}

func (h *Hub) BroadcastState() {
	state := h.GameEngine.GetState()

	// Wrap in message
	msg := struct {
		Type    string         `json:"type"`
		Payload game.GameState `json:"payload"`
	}{
		Type:    "GAME_STATE",
		Payload: state,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling state: %v", err)
		return
	}

	h.Broadcast <- bytes
}
