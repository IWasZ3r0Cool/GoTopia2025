package server

import (
	"encoding/json"

	"github.com/IWasZ3r0Cool/GoTopia2025/internal/game"
)

const (
	MsgJoinGame  = "JOIN_GAME"
	MsgBuild     = "BUILD"
	MsgWelcome   = "WELCOME"
	MsgGameState = "GAME_STATE"
	MsgError     = "ERROR"
)

type incomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type JoinPayload struct {
	Name string `json:"name"`
}

type BuildPayload struct {
	BuildingType game.BuildingType `json:"buildingType"`
	X            int               `json:"x"`
	Y            int               `json:"y"`
}

type WelcomePayload struct {
	PlayerID string `json:"playerId"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type outgoingMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}
