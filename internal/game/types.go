package game

// Coordinate represents a position on the 2D grid
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type BuildingType string

const (
	BuildingHouse    BuildingType = "HOUSE"
	BuildingFactory  BuildingType = "FACTORY"
	BuildingFarm     BuildingType = "FARM"
	BuildingFort     BuildingType = "FORT"
	BuildingSchool   BuildingType = "SCHOOL"
	BuildingHospital BuildingType = "HOSPITAL"
)

type Building struct {
	Type   BuildingType `json:"type"`
	Health int          `json:"health"`
}

// Island represents a player's territory
type Island struct {
	OwnerID   string              `json:"ownerId"`
	Buildings map[string]Building `json:"buildings"` // Key format: "x,y"
	// Island Bounds
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Player struct {
	ID         string `json:"id"`
	Gold       int    `json:"gold"`
	Population int    `json:"population"`
	Mood       int    `json:"mood"` // Happiness 0-100
}

// GameState holds the entire synchronized state of the world
type GameState struct {
	MapWidth  int                `json:"mapWidth"`
	MapHeight int                `json:"mapHeight"`
	Islands   map[string]*Island `json:"islands"` // Keyed by PlayerID
	Players   map[string]*Player `json:"players"` // Keyed by PlayerID
	Turn      int                `json:"turn"`    // Current turn number
	Round     int                `json:"round"`   // Current round (maybe redundant with Turn, depending on logic)
}
