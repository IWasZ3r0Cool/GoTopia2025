package game

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	ErrInvalidPlayer    = errors.New("player name and ID are required")
	ErrInvalidConfig    = errors.New("invalid game configuration")
	ErrDuplicatePlayer  = errors.New("player ID is already in use")
	ErrGameFull         = errors.New("game is full")
	ErrPlayerNotFound   = errors.New("player not found")
	ErrInvalidBuilding  = errors.New("unknown building type")
	ErrOutsideIsland    = errors.New("tile is outside the player's island")
	ErrTileOccupied     = errors.New("tile is already occupied")
	ErrInsufficientGold = errors.New("not enough gold")
)

// GameEngine owns all mutable game state. Callers only receive deep snapshots,
// so encoding a state update can never race with a game mutation.
type GameEngine struct {
	mu      sync.RWMutex
	state   GameState
	config  Config
	updates chan struct{}

	lifecycleMu sync.Mutex
	running     bool
	cancel      context.CancelFunc
	done        chan struct{}
}

func NewGameEngine() *GameEngine {
	engine, err := NewGameEngineWithConfig(DefaultConfig())
	if err != nil {
		panic(fmt.Sprintf("default game configuration is invalid: %v", err))
	}
	return engine
}

func NewGameEngineWithConfig(config Config) (*GameEngine, error) {
	defaults := DefaultConfig()
	if config.MapWidth == 0 {
		config.MapWidth = defaults.MapWidth
	}
	if config.MapHeight == 0 {
		config.MapHeight = defaults.MapHeight
	}
	if config.IslandWidth == 0 {
		config.IslandWidth = defaults.IslandWidth
	}
	if config.IslandHeight == 0 {
		config.IslandHeight = defaults.IslandHeight
	}
	if config.MaxPlayers == 0 {
		config.MaxPlayers = defaults.MaxPlayers
	}
	if config.TurnDuration == 0 {
		config.TurnDuration = defaults.TurnDuration
	}
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return &GameEngine{
		state: GameState{
			MapWidth: config.MapWidth, MapHeight: config.MapHeight,
			Islands: make(map[string]*Island), Players: make(map[string]*Player), Turn: 1,
		},
		config:  config,
		updates: make(chan struct{}, 1),
	}, nil
}

func validateConfig(config Config) error {
	if config.MapWidth < 1 || config.MapHeight < 1 || config.IslandWidth < 1 || config.IslandHeight < 1 {
		return fmt.Errorf("%w: map and island dimensions must be positive", ErrInvalidConfig)
	}
	if config.MaxPlayers < 1 || config.MaxPlayers > MaxPlayers {
		return fmt.Errorf("%w: max players must be between 1 and %d", ErrInvalidConfig, MaxPlayers)
	}
	if config.TurnDuration <= 0 {
		return fmt.Errorf("%w: turn duration must be positive", ErrInvalidConfig)
	}

	requiredWidth := config.IslandWidth + 2*IslandMargin
	requiredHeight := config.IslandHeight + 2*IslandMargin
	if config.MaxPlayers >= 2 {
		requiredWidth = 2*config.IslandWidth + 2*IslandMargin + MinimumIslandGap
	}
	if config.MaxPlayers >= 3 {
		requiredHeight = 2*config.IslandHeight + 2*IslandMargin + MinimumIslandGap
	}
	if config.MapWidth < requiredWidth || config.MapHeight < requiredHeight {
		return fmt.Errorf(
			"%w: a %dx%d map needs to be at least %dx%d for %d %dx%d islands",
			ErrInvalidConfig, config.MapWidth, config.MapHeight, requiredWidth, requiredHeight,
			config.MaxPlayers, config.IslandWidth, config.IslandHeight,
		)
	}
	return nil
}

// Start begins the turn clock. It is safe to call Start more than once.
func (e *GameEngine) Start(parent context.Context) {
	e.lifecycleMu.Lock()
	defer e.lifecycleMu.Unlock()
	if e.running {
		return
	}
	ctx, cancel := context.WithCancel(parent)
	e.cancel = cancel
	e.done = make(chan struct{})
	e.running = true
	go e.run(ctx, e.done)
}

func (e *GameEngine) run(ctx context.Context, done chan struct{}) {
	defer close(done)
	ticker := time.NewTicker(e.config.TurnDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.AdvanceTurn()
		}
	}
}

// Stop ends the turn clock and waits for its goroutine to exit.
func (e *GameEngine) Stop() {
	e.lifecycleMu.Lock()
	if !e.running {
		e.lifecycleMu.Unlock()
		return
	}
	cancel, done := e.cancel, e.done
	e.running = false
	e.lifecycleMu.Unlock()
	cancel()
	<-done
}

// Updates emits a coalesced signal after every state mutation.
func (e *GameEngine) Updates() <-chan struct{} { return e.updates }

func (e *GameEngine) notify() {
	select {
	case e.updates <- struct{}{}:
	default:
	}
}

func (e *GameEngine) AddPlayer(id, name string) error {
	id, name = strings.TrimSpace(id), strings.TrimSpace(name)
	if id == "" || name == "" {
		return ErrInvalidPlayer
	}

	e.mu.Lock()
	if _, exists := e.state.Players[id]; exists {
		e.mu.Unlock()
		return ErrDuplicatePlayer
	}
	if len(e.state.Players) >= e.config.MaxPlayers {
		e.mu.Unlock()
		return ErrGameFull
	}

	slot := e.availableSlot()
	x, y := e.islandOrigin(slot)
	e.state.Players[id] = &Player{
		ID: id, Name: name, Gold: StartingGold,
		Population: StartingPopulation, PopulationCapacity: StartingPopCapacity, Mood: StartingMood,
	}
	e.state.Islands[id] = &Island{
		OwnerID: id, Buildings: make(map[string]Building), X: x, Y: y,
		Width: e.config.IslandWidth, Height: e.config.IslandHeight,
	}
	e.mu.Unlock()
	e.notify()
	return nil
}

func (e *GameEngine) availableSlot() int {
	for slot := 0; slot < e.config.MaxPlayers; slot++ {
		x, y := e.islandOrigin(slot)
		occupied := false
		for _, island := range e.state.Islands {
			if island.X == x && island.Y == y {
				occupied = true
				break
			}
		}
		if !occupied {
			return slot
		}
	}
	return 0
}

func (e *GameEngine) islandOrigin(slot int) (int, int) {
	right := e.config.MapWidth - e.config.IslandWidth - IslandMargin
	bottom := e.config.MapHeight - e.config.IslandHeight - IslandMargin
	positions := [MaxPlayers][2]int{{IslandMargin, IslandMargin}, {right, IslandMargin}, {IslandMargin, bottom}, {right, bottom}}
	return positions[slot][0], positions[slot][1]
}

func (e *GameEngine) RemovePlayer(id string) bool {
	e.mu.Lock()
	if _, exists := e.state.Players[id]; !exists {
		e.mu.Unlock()
		return false
	}
	delete(e.state.Players, id)
	delete(e.state.Islands, id)
	e.mu.Unlock()
	e.notify()
	return true
}

func (e *GameEngine) Build(playerID string, buildingType BuildingType, x, y int) error {
	spec, valid := SpecForBuilding(buildingType)
	if !valid {
		return ErrInvalidBuilding
	}

	e.mu.Lock()
	player, ok := e.state.Players[playerID]
	if !ok {
		e.mu.Unlock()
		return ErrPlayerNotFound
	}
	island := e.state.Islands[playerID]
	if x < island.X || x >= island.X+island.Width || y < island.Y || y >= island.Y+island.Height {
		e.mu.Unlock()
		return ErrOutsideIsland
	}
	key := coordinateKey(x, y)
	if _, occupied := island.Buildings[key]; occupied {
		e.mu.Unlock()
		return ErrTileOccupied
	}
	if player.Gold < spec.Cost {
		e.mu.Unlock()
		return ErrInsufficientGold
	}
	player.Gold -= spec.Cost
	island.Buildings[key] = Building{Type: buildingType, Health: 100}
	if buildingType == BuildingHouse {
		player.PopulationCapacity += PopCapHouse
	}
	e.mu.Unlock()
	e.notify()
	return nil
}

func coordinateKey(x, y int) string { return fmt.Sprintf("%d,%d", x, y) }

func (e *GameEngine) AdvanceTurn() {
	e.mu.Lock()
	e.state.Turn++
	for playerID, island := range e.state.Islands {
		player := e.state.Players[playerID]
		for _, building := range island.Buildings {
			if spec, ok := SpecForBuilding(building.Type); ok {
				player.Gold += spec.Income
			}
		}
		if player.Population < player.PopulationCapacity {
			growth := max(1, player.Population/20)
			player.Population = min(player.Population+growth, player.PopulationCapacity)
		}
	}
	e.mu.Unlock()
	e.notify()
}

func (e *GameEngine) Snapshot() GameState {
	e.mu.RLock()
	defer e.mu.RUnlock()

	snapshot := GameState{
		MapWidth: e.state.MapWidth, MapHeight: e.state.MapHeight, Turn: e.state.Turn,
		Players: make(map[string]*Player, len(e.state.Players)),
		Islands: make(map[string]*Island, len(e.state.Islands)),
	}
	for id, player := range e.state.Players {
		copy := *player
		snapshot.Players[id] = &copy
	}
	for id, island := range e.state.Islands {
		copy := *island
		copy.Buildings = make(map[string]Building, len(island.Buildings))
		for key, building := range island.Buildings {
			copy.Buildings[key] = building
		}
		snapshot.Islands[id] = &copy
	}
	return snapshot
}
