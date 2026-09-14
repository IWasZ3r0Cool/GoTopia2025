package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/IWasZ3r0Cool/GoTopia2025/internal/game"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: originAllowed,
}

type Client struct {
	Hub  *Hub
	ID   string
	conn *websocket.Conn
	send chan []byte
	done chan struct{}
	once sync.Once
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.removeClient(c)
		c.closeConnection()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("read websocket: %v", err)
			}
			return
		}
		c.handleMessage(message)
	}
}

func (c *Client) handleMessage(message []byte) {
	var envelope incomingMessage
	if err := decodeStrict(message, &envelope); err != nil {
		c.sendError("INVALID_MESSAGE", "Message must be valid JSON with known fields.")
		return
	}

	switch envelope.Type {
	case MsgJoinGame:
		var payload JoinPayload
		if err := decodeStrict(envelope.Payload, &payload); err != nil {
			c.sendError("INVALID_JOIN", "Join requests require a player name.")
			return
		}
		if err := c.Hub.joinGame(c, payload.Name); err != nil {
			c.sendError(joinErrorCode(err), err.Error())
		}

	case MsgBuild:
		if c.ID == "" {
			c.sendError("NOT_JOINED", "Join the game before building.")
			return
		}
		var payload BuildPayload
		if err := decodeStrict(envelope.Payload, &payload); err != nil {
			c.sendError("INVALID_BUILD", "Build requests require a building type and coordinates.")
			return
		}
		if err := c.Hub.engine.Build(c.ID, payload.BuildingType, payload.X, payload.Y); err != nil {
			c.sendError(buildErrorCode(err), err.Error())
		}

	default:
		c.sendError("UNKNOWN_MESSAGE", "Unknown message type.")
	}
}

func decodeStrict(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("message contains trailing data")
	}
	return nil
}

func joinErrorCode(err error) string {
	switch {
	case errors.Is(err, game.ErrGameFull):
		return "GAME_FULL"
	case errors.Is(err, errAlreadyJoined):
		return "ALREADY_JOINED"
	default:
		return "INVALID_JOIN"
	}
}

func buildErrorCode(err error) string {
	switch {
	case errors.Is(err, game.ErrOutsideIsland):
		return "OUTSIDE_ISLAND"
	case errors.Is(err, game.ErrTileOccupied):
		return "TILE_OCCUPIED"
	case errors.Is(err, game.ErrInsufficientGold):
		return "INSUFFICIENT_GOLD"
	case errors.Is(err, game.ErrInvalidBuilding):
		return "INVALID_BUILDING"
	default:
		return "BUILD_FAILED"
	}
}

func (c *Client) sendJSON(message outgoingMessage) {
	encoded, err := json.Marshal(message)
	if err != nil {
		log.Printf("encode websocket message: %v", err)
		return
	}
	select {
	case c.send <- encoded:
	case <-c.done:
		return
	default:
		c.closeConnection()
	}
}

func (c *Client) sendError(code, message string) {
	c.sendJSON(outgoingMessage{Type: MsgError, Payload: ErrorPayload{Code: code, Message: message}})
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.closeConnection()
	}()
	for {
		select {
		case <-c.done:
			return
		case message := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) closeConnection() {
	c.once.Do(func() {
		close(c.done)
		_ = c.conn.Close()
	})
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{Hub: hub, conn: connection, send: make(chan []byte, 64), done: make(chan struct{})}
	if !hub.addClient(client) {
		_ = connection.Close()
		return
	}
	go client.writePump()
	go client.readPump()
}

func originAllowed(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	for _, allowed := range strings.Split(os.Getenv("GOTOPIA_ALLOWED_ORIGINS"), ",") {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}
	parsed, err := url.Parse(origin)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.User != nil || (parsed.Path != "" && parsed.Path != "/") || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	requestScheme := "http"
	if r.TLS != nil {
		requestScheme = "https"
	}
	return strings.EqualFold(parsed.Scheme, requestScheme) && strings.EqualFold(parsed.Host, r.Host)
}
