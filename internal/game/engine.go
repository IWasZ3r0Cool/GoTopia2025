package game

import (
	"fmt"
	"sync"
	"time"
)

type GameEngine struct {
	State      *GameState
	mutex      sync.RWMutex
	ticker     *time.Ticker
	turnTicker *time.Ticker
	stop       chan bool
}

func NewGameEngine() *GameEngine {
	// Initialize default state
	state := &GameState{
		MapWidth:  MapWidth,
		MapHeight: MapHeight,
		Islands:   make(map[string]*Island),
		Players:   make(map[string]*Player),
		Turn:      1,
		Round:     1,
	}

	return &GameEngine{
		State: state,
		stop:  make(chan bool),
	}
}

func (ge *GameEngine) Start() {
	ge.ticker = time.NewTicker(200 * time.Millisecond) // 5 ticks per second for smooth updates
	ge.turnTicker = time.NewTicker(time.Duration(TurnDurationSeconds) * time.Second)

	go ge.loop()
}

func (ge *GameEngine) loop() {
	for {
		select {
		case <-ge.stop:
			ge.ticker.Stop()
			ge.turnTicker.Stop()
			return
		case <-ge.ticker.C:
			ge.processTick()
		case <-ge.turnTicker.C:
			ge.processTurn()
		}
	}
}

func (ge *GameEngine) processTick() {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	// Logic for real-time events (boat movement, etc.) will go here
}

func (ge *GameEngine) processTurn() {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	ge.State.Turn++

	// Economic calculations
	for playerID, island := range ge.State.Islands {
		player, exists := ge.State.Players[playerID]
		if !exists {
			continue
		}

		income := 0
		for _, building := range island.Buildings {
			switch building.Type {
			case BuildingFactory:
				income += IncomeFactory
			case BuildingFarm:
				income += IncomeFarm
			}
		}

		player.Gold += income
	}
}

func (ge *GameEngine) AddPlayer(id string) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	ge.State.Players[id] = &Player{
		ID:         id,
		Gold:       1000,
		Population: 100,
		Mood:       100,
	}

	// Assign Island Slot based on number of players
	// Slot 0: Top-Left (10, 10)
	// Slot 1: Top-Right (165, 10)
	// Slot 2: Bottom-Left (10, 165)
	// Slot 3: Bottom-Right (165, 165)

	playerCount := len(ge.State.Islands)
	var x, y int

	switch playerCount {
	case 0:
		x, y = 10, 10
	case 1:
		x, y = 165, 10
	case 2:
		x, y = 10, 165
	case 3:
		x, y = 165, 165
	default:
		// Spectator or error handling? For now, just dump them at 0,0 or reject.
		x, y = 0, 0
	}

	ge.State.Islands[id] = &Island{
		OwnerID:   id,
		Buildings: make(map[string]Building),
		X:         x,
		Y:         y,
		Width:     25,
		Height:    25,
	}
}

func (ge *GameEngine) HandleBuild(playerID string, buildingType BuildingType, x, y int) bool {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	player, ok := ge.State.Players[playerID]
	if !ok {
		return false
	}

	island, ok := ge.State.Islands[playerID]
	if !ok {
		return false
	}

	// 1. Check Bounds
	if x < island.X || x >= island.X+island.Width || y < island.Y || y >= island.Y+island.Height {
		return false
	}

	// 2. Check if tile is empty
	coordKey := fmt.Sprintf("%d,%d", x, y)
	if _, exists := island.Buildings[coordKey]; exists {
		return false
	}

	// 3. Check Cost and Deduct Gold
	cost := 0
	switch buildingType {
	case BuildingHouse:
		cost = CostHouse
	case BuildingFactory:
		cost = CostFactory
	case BuildingFarm:
		cost = CostFarm
	case BuildingFort:
		cost = CostFort
	case BuildingSchool:
		cost = CostSchool
	case BuildingHospital:
		cost = CostHospital
	default:
		return false
	}

	if player.Gold < cost {
		return false
	}

	player.Gold -= cost
	island.Buildings[coordKey] = Building{
		Type:   buildingType,
		Health: 100,
	}

	return true
}

func (ge *GameEngine) GetState() GameState {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	// Return a copy or serialization recommended,
	// but for now returning the struct (careful with maps)
	// In a real app we'd deep copy or just serialize to JSON here.
	return *ge.State
}
