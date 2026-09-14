package game

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

type Island struct {
	OwnerID   string              `json:"ownerId"`
	Buildings map[string]Building `json:"buildings"`
	X         int                 `json:"x"`
	Y         int                 `json:"y"`
	Width     int                 `json:"width"`
	Height    int                 `json:"height"`
}

type Player struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Gold               int    `json:"gold"`
	Population         int    `json:"population"`
	PopulationCapacity int    `json:"populationCapacity"`
	Mood               int    `json:"mood"`
}

type GameState struct {
	MapWidth  int                `json:"mapWidth"`
	MapHeight int                `json:"mapHeight"`
	Islands   map[string]*Island `json:"islands"`
	Players   map[string]*Player `json:"players"`
	Turn      int                `json:"turn"`
}
