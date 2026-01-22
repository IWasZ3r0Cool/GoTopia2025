package server

import "github.com/IWasZ3r0Cool/GoTopia2025/internal/game"

// Message Types
const (
	MsgJoinGame  = "JOIN_GAME"
	MsgBuild     = "BUILD"
	MsgGameState = "GAME_STATE"
	MsgError     = "ERROR"
)

// BaseMessage acts as the envelope
type BaseMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

// ClientPayloads
type JoinPayload struct {
	Name string `json:"name"`
}

type BuildPayload struct {
	BuildingType game.BuildingType `json:"buildingType"`
	X            int               `json:"x"`
	Y            int               `json:"y"`
}

// ServerPayloads
type ErrorPayload struct {
	Message string `json:"message"`
}
