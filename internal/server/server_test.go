package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/IWasZ3r0Cool/GoTopia2025/internal/game"
	"github.com/gorilla/websocket"
)

type wireMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func startTestServer(t *testing.T) (*game.GameEngine, string) {
	t.Helper()
	engine := game.NewGameEngine()
	hub := NewHub(engine)
	ctx, cancel := context.WithCancel(context.Background())
	go hub.Run(ctx)
	httpServer := httptest.NewServer(NewHTTPHandler(hub, t.TempDir()))
	t.Cleanup(func() {
		httpServer.Close()
		cancel()
		select {
		case <-hub.done:
		case <-time.After(time.Second):
			t.Error("hub did not stop")
		}
	})
	return engine, "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
}

func dialTestClient(t *testing.T, address string) *websocket.Conn {
	t.Helper()
	connection, response, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		if response != nil {
			t.Fatalf("dial websocket: %v (HTTP %s)", err, response.Status)
		}
		t.Fatalf("dial websocket: %v", err)
	}
	t.Cleanup(func() { _ = connection.Close() })
	return connection
}

func readMessage(t *testing.T, connection *websocket.Conn) wireMessage {
	t.Helper()
	_ = connection.SetReadDeadline(time.Now().Add(time.Second))
	var message wireMessage
	if err := connection.ReadJSON(&message); err != nil {
		t.Fatalf("read websocket message: %v", err)
	}
	return message
}

func joinClient(t *testing.T, connection *websocket.Conn, name string) (string, game.GameState) {
	t.Helper()
	if err := connection.WriteJSON(outgoingMessage{Type: MsgJoinGame, Payload: JoinPayload{Name: name}}); err != nil {
		t.Fatalf("send join: %v", err)
	}
	var playerID string
	var state game.GameState
	for playerID == "" || state.Players == nil {
		message := readMessage(t, connection)
		switch message.Type {
		case MsgWelcome:
			var payload WelcomePayload
			if err := json.Unmarshal(message.Payload, &payload); err != nil {
				t.Fatal(err)
			}
			playerID = payload.PlayerID
		case MsgGameState:
			if err := json.Unmarshal(message.Payload, &state); err != nil {
				t.Fatal(err)
			}
		default:
			t.Fatalf("unexpected message %q", message.Type)
		}
	}
	return playerID, state
}

func TestOriginAllowedRequiresAnExactOrigin(t *testing.T) {
	tests := []struct {
		name    string
		request *http.Request
		origin  string
		allowed bool
	}{
		{
			name:    "non-browser client without origin",
			request: httptest.NewRequest(http.MethodGet, "http://game.example/ws", nil),
			allowed: true,
		},
		{
			name:    "same HTTP origin",
			request: httptest.NewRequest(http.MethodGet, "http://game.example/ws", nil),
			origin:  "http://game.example",
			allowed: true,
		},
		{
			name:    "same HTTPS origin",
			request: httptest.NewRequest(http.MethodGet, "https://game.example/ws", nil),
			origin:  "https://game.example",
			allowed: true,
		},
		{
			name:    "different port",
			request: httptest.NewRequest(http.MethodGet, "https://game.example/ws", nil),
			origin:  "https://game.example:4444",
			allowed: false,
		},
		{
			name:    "different scheme",
			request: httptest.NewRequest(http.MethodGet, "https://game.example/ws", nil),
			origin:  "http://game.example",
			allowed: false,
		},
		{
			name:    "different host",
			request: httptest.NewRequest(http.MethodGet, "https://game.example/ws", nil),
			origin:  "https://attacker.example",
			allowed: false,
		},
		{
			name:    "origin with path",
			request: httptest.NewRequest(http.MethodGet, "https://game.example/ws", nil),
			origin:  "https://game.example/not-an-origin",
			allowed: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.origin != "" {
				test.request.Header.Set("Origin", test.origin)
			}
			if actual := originAllowed(test.request); actual != test.allowed {
				t.Fatalf("originAllowed() = %t, want %t", actual, test.allowed)
			}
		})
	}
}

func TestOriginAllowedAcceptsConfiguredOrigin(t *testing.T) {
	t.Setenv("GOTOPIA_ALLOWED_ORIGINS", "https://admin.example, https://game.example:4444")
	request := httptest.NewRequest(http.MethodGet, "https://game.example/ws", nil)
	request.Header.Set("Origin", "https://game.example:4444")
	if !originAllowed(request) {
		t.Fatal("explicitly configured origin was rejected")
	}
}

func TestWebSocketJoinAndBuildFlow(t *testing.T) {
	_, address := startTestServer(t)
	connection := dialTestClient(t, address)
	playerID, state := joinClient(t, connection, "Ada")
	if state.Players[playerID].Name != "Ada" {
		t.Fatalf("player was not present in state: %#v", state.Players)
	}

	island := state.Islands[playerID]
	build := BuildPayload{BuildingType: game.BuildingFarm, X: island.X, Y: island.Y}
	if err := connection.WriteJSON(outgoingMessage{Type: MsgBuild, Payload: build}); err != nil {
		t.Fatalf("send build: %v", err)
	}
	message := readMessage(t, connection)
	if message.Type != MsgGameState {
		t.Fatalf("expected state update, got %q", message.Type)
	}
	if err := json.Unmarshal(message.Payload, &state); err != nil {
		t.Fatal(err)
	}
	if _, exists := state.Islands[playerID].Buildings["10,10"]; !exists {
		t.Fatal("built farm missing from state")
	}
	if state.Players[playerID].Gold != game.StartingGold-100 {
		t.Fatalf("unexpected gold: %d", state.Players[playerID].Gold)
	}
}

func TestWebSocketRejectsInvalidAction(t *testing.T) {
	_, address := startTestServer(t)
	connection := dialTestClient(t, address)
	_, _ = joinClient(t, connection, "Grace")

	if err := connection.WriteJSON(outgoingMessage{Type: MsgBuild, Payload: BuildPayload{
		BuildingType: game.BuildingHouse, X: -1, Y: -1,
	}}); err != nil {
		t.Fatal(err)
	}
	message := readMessage(t, connection)
	if message.Type != MsgError {
		t.Fatalf("expected error, got %q", message.Type)
	}
	var payload ErrorPayload
	if err := json.Unmarshal(message.Payload, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Code != "OUTSIDE_ISLAND" {
		t.Fatalf("expected OUTSIDE_ISLAND, got %q", payload.Code)
	}
}

func TestStateBroadcastsToExistingClients(t *testing.T) {
	_, address := startTestServer(t)
	first := dialTestClient(t, address)
	_, _ = joinClient(t, first, "First")
	second := dialTestClient(t, address)
	_, _ = joinClient(t, second, "Second")

	message := readMessage(t, first)
	if message.Type != MsgGameState {
		t.Fatalf("expected state broadcast, got %q", message.Type)
	}
	var state game.GameState
	if err := json.Unmarshal(message.Payload, &state); err != nil {
		t.Fatal(err)
	}
	if len(state.Players) != 2 {
		t.Fatalf("expected two players, got %d", len(state.Players))
	}
}

func TestHTTPHealth(t *testing.T) {
	engine := game.NewGameEngine()
	handler := NewHTTPHandler(NewHub(engine), t.TempDir())
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"status":"ok"`) {
		t.Fatalf("unexpected health response: %d %s", response.Code, response.Body.String())
	}
}
