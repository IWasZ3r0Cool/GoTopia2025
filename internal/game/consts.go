package game

import "time"

const (
	MapWidth            = 200
	MapHeight           = 200
	IslandWidth         = 25
	IslandHeight        = 25
	IslandMargin        = 10
	MinimumIslandGap    = 2
	MaxPlayers          = 4
	StartingGold        = 1000
	StartingPopulation  = 100
	StartingPopCapacity = 100
	StartingMood        = 100
	TurnDurationSeconds = 60

	IncomeFactory = 20
	IncomeFarm    = 5
	PopCapHouse   = 50
	RangeFort     = 5
)

// Config contains the values that may vary in tests or custom games. Zero
// values are replaced with their defaults.
type Config struct {
	MapWidth     int
	MapHeight    int
	IslandWidth  int
	IslandHeight int
	MaxPlayers   int
	TurnDuration time.Duration
}

func DefaultConfig() Config {
	return Config{
		MapWidth: MapWidth, MapHeight: MapHeight,
		IslandWidth: IslandWidth, IslandHeight: IslandHeight,
		MaxPlayers: MaxPlayers, TurnDuration: TurnDurationSeconds * time.Second,
	}
}

type BuildingSpec struct {
	Cost   int
	Income int
}

var buildingSpecs = map[BuildingType]BuildingSpec{
	BuildingHouse:    {Cost: 150},
	BuildingFarm:     {Cost: 100, Income: IncomeFarm},
	BuildingFactory:  {Cost: 200, Income: IncomeFactory},
	BuildingFort:     {Cost: 200},
	BuildingSchool:   {Cost: 300},
	BuildingHospital: {Cost: 300},
}

func SpecForBuilding(buildingType BuildingType) (BuildingSpec, bool) {
	spec, ok := buildingSpecs[buildingType]
	return spec, ok
}
