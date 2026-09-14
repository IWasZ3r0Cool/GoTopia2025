package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/IWasZ3r0Cool/GoTopia2025/internal/game"
)

var (
	errAlreadyJoined = errors.New("client has already joined")
	errNameTooLong   = errors.New("player name must be 24 characters or fewer")
	errHubStopped    = errors.New("game server has stopped")
)

type joinRequest struct {
	client *Client
	name   string
	result chan error
}

type Hub struct {
	engine     *game.GameEngine
	clients    map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
	join       chan joinRequest
	done       chan struct{}
	nextID     uint64
}

func NewHub(engine *game.GameEngine) *Hub {
	return &Hub{
		engine: engine, clients: make(map[*Client]struct{}),
		register: make(chan *Client), unregister: make(chan *Client),
		join: make(chan joinRequest), done: make(chan struct{}),
	}
}

// Run owns the client collection and blocks until the context is cancelled.
func (h *Hub) Run(ctx context.Context) {
	h.engine.Start(ctx)
	defer h.engine.Stop()
	defer close(h.done)

	for {
		select {
		case <-ctx.Done():
			for client := range h.clients {
				client.closeConnection()
			}
			return

		case client := <-h.register:
			h.clients[client] = struct{}{}

		case client := <-h.unregister:
			if _, exists := h.clients[client]; !exists {
				continue
			}
			delete(h.clients, client)
			client.closeConnection()
			if client.ID != "" {
				h.engine.RemovePlayer(client.ID)
			}

		case request := <-h.join:
			request.result <- h.handleJoin(request.client, request.name)

		case <-h.engine.Updates():
			h.broadcastState()
		}
	}
}

func (h *Hub) handleJoin(client *Client, name string) error {
	name = strings.TrimSpace(name)
	if client.ID != "" {
		return errAlreadyJoined
	}
	if len([]rune(name)) > 24 {
		return errNameTooLong
	}
	h.nextID++
	playerID := fmt.Sprintf("player-%d", h.nextID)
	if err := h.engine.AddPlayer(playerID, name); err != nil {
		return err
	}
	client.ID = playerID
	client.sendJSON(outgoingMessage{Type: MsgWelcome, Payload: WelcomePayload{PlayerID: playerID}})
	return nil
}

func (h *Hub) joinGame(client *Client, name string) error {
	request := joinRequest{client: client, name: name, result: make(chan error, 1)}
	select {
	case h.join <- request:
	case <-h.done:
		return errHubStopped
	}
	select {
	case err := <-request.result:
		return err
	case <-h.done:
		return errHubStopped
	}
}

func (h *Hub) addClient(client *Client) bool {
	select {
	case h.register <- client:
		return true
	case <-h.done:
		return false
	}
}

func (h *Hub) removeClient(client *Client) {
	select {
	case h.unregister <- client:
	case <-h.done:
	}
}

func (h *Hub) broadcastState() {
	encoded, err := json.Marshal(outgoingMessage{Type: MsgGameState, Payload: h.engine.Snapshot()})
	if err != nil {
		log.Printf("encode game state: %v", err)
		return
	}
	for client := range h.clients {
		select {
		case client.send <- encoded:
		default:
			client.closeConnection()
			delete(h.clients, client)
			if client.ID != "" {
				h.engine.RemovePlayer(client.ID)
			}
		}
	}
}
