package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/IWasZ3r0Cool/GoTopia2025/internal/game"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub  *Hub
	ID   string
	conn *websocket.Conn
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.handleMessage(message)
	}
}

func (c *Client) handleMessage(msg []byte) {
	var base BaseMessage
	if err := json.Unmarshal(msg, &base); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return
	}

	switch base.Type {
	case MsgJoinGame:
		// Simple autojoin logic for now. ID should be generated or passed.
		// For prototype, let's assume Client.ID is set on connection or we use payload.
		// NOTE: AddPlayer should be called here.
		if c.ID == "" {
			// Generate valid ID or use payload
			// Assuming payload has Name, use it as ID for now
			var payload JoinPayload
			// Re-marshal to get payload... naive but works.
			// Better: map[string]interface{}
			pBytes, _ := json.Marshal(base.Payload)
			json.Unmarshal(pBytes, &payload)

			if payload.Name != "" {
				c.ID = payload.Name
				c.Hub.GameEngine.AddPlayer(c.ID)

				// Send initial state
				c.sendState()
			}
		}

	case MsgBuild:
		if c.ID == "" {
			return
		}
		var payload BuildPayload
		pBytes, _ := json.Marshal(base.Payload)
		json.Unmarshal(pBytes, &payload)

		success := c.Hub.GameEngine.HandleBuild(c.ID, payload.BuildingType, payload.X, payload.Y)
		if success {
			// Broadcast new state to ALL
			c.Hub.BroadcastState()
		} else {
			// Send Error to Just This Client
			c.sendError("Build Failed: Invalid location or insufficient funds")
		}
	}
}

func (c *Client) sendState() {
	state := c.Hub.GameEngine.GetState()
	bytes, _ := json.Marshal(struct {
		Type    string         `json:"type"`
		Payload game.GameState `json:"payload"`
	}{
		Type:    MsgGameState,
		Payload: state,
	})
	c.send <- bytes
}

func (c *Client) sendError(msg string) {
	bytes, _ := json.Marshal(struct {
		Type    string       `json:"type"`
		Payload ErrorPayload `json:"payload"`
	}{
		Type:    MsgError,
		Payload: ErrorPayload{Message: msg},
	})
	c.send <- bytes
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{Hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
