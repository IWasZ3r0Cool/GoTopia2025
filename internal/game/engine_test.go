package game

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestAddPlayerAssignsUniqueIslandsAndReusesOpenSlot(t *testing.T) {
	engine := NewGameEngine()
	for index, id := range []string{"p1", "p2", "p3", "p4"} {
		if err := engine.AddPlayer(id, "Player"); err != nil {
			t.Fatalf("add player %d: %v", index, err)
		}
	}
	if err := engine.AddPlayer("p5", "Player"); !errors.Is(err, ErrGameFull) {
		t.Fatalf("expected game full, got %v", err)
	}

	before := engine.Snapshot().Islands["p2"]
	if !engine.RemovePlayer("p2") {
		t.Fatal("expected player removal to succeed")
	}
	if err := engine.AddPlayer("p5", "Replacement"); err != nil {
		t.Fatalf("replace player: %v", err)
	}
	after := engine.Snapshot().Islands["p5"]
	if after.X != before.X || after.Y != before.Y {
		t.Fatalf("expected open island slot %d,%d; got %d,%d", before.X, before.Y, after.X, after.Y)
	}
}

func TestAddPlayerValidation(t *testing.T) {
	engine := NewGameEngine()
	if err := engine.AddPlayer("", "Player"); !errors.Is(err, ErrInvalidPlayer) {
		t.Fatalf("expected invalid player, got %v", err)
	}
	if err := engine.AddPlayer("p1", "Player"); err != nil {
		t.Fatal(err)
	}
	if err := engine.AddPlayer("p1", "Other"); !errors.Is(err, ErrDuplicatePlayer) {
		t.Fatalf("expected duplicate player, got %v", err)
	}
}

func TestBuildValidationAndCost(t *testing.T) {
	engine := NewGameEngine()
	if err := engine.AddPlayer("p1", "Player"); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		kind BuildingType
		x, y int
		want error
	}{
		{name: "unknown player", kind: BuildingHouse, x: 10, y: 10, want: ErrPlayerNotFound},
		{name: "unknown building", kind: "CASTLE", x: 10, y: 10, want: ErrInvalidBuilding},
		{name: "outside island", kind: BuildingHouse, x: 9, y: 10, want: ErrOutsideIsland},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			playerID := "p1"
			if test.name == "unknown player" {
				playerID = "missing"
			}
			if err := engine.Build(playerID, test.kind, test.x, test.y); !errors.Is(err, test.want) {
				t.Fatalf("expected %v, got %v", test.want, err)
			}
		})
	}

	if err := engine.Build("p1", BuildingHouse, 10, 10); err != nil {
		t.Fatalf("valid build: %v", err)
	}
	if err := engine.Build("p1", BuildingFarm, 10, 10); !errors.Is(err, ErrTileOccupied) {
		t.Fatalf("expected occupied tile, got %v", err)
	}
	state := engine.Snapshot()
	if got := state.Players["p1"].Gold; got != StartingGold-150 {
		t.Fatalf("expected %d gold, got %d", StartingGold-150, got)
	}
}

func TestTurnAddsBuildingIncome(t *testing.T) {
	engine := NewGameEngine()
	if err := engine.AddPlayer("p1", "Player"); err != nil {
		t.Fatal(err)
	}
	if err := engine.Build("p1", BuildingFactory, 10, 10); err != nil {
		t.Fatal(err)
	}
	if err := engine.Build("p1", BuildingFarm, 11, 10); err != nil {
		t.Fatal(err)
	}
	if err := engine.Build("p1", BuildingHouse, 12, 10); err != nil {
		t.Fatal(err)
	}
	engine.AdvanceTurn()

	state := engine.Snapshot()
	if state.Turn != 2 {
		t.Fatalf("expected turn 2, got %d", state.Turn)
	}
	wantGold := StartingGold - 200 - 100 - 150 + IncomeFactory + IncomeFarm
	if state.Players["p1"].Gold != wantGold {
		t.Fatalf("expected %d gold, got %d", wantGold, state.Players["p1"].Gold)
	}
	if state.Players["p1"].PopulationCapacity != StartingPopCapacity+PopCapHouse {
		t.Fatalf("house did not increase population capacity: %d", state.Players["p1"].PopulationCapacity)
	}
	if state.Players["p1"].Population != 105 {
		t.Fatalf("expected population growth to 105, got %d", state.Players["p1"].Population)
	}
}

func TestSnapshotIsDeepCopy(t *testing.T) {
	engine := NewGameEngine()
	if err := engine.AddPlayer("p1", "Player"); err != nil {
		t.Fatal(err)
	}
	snapshot := engine.Snapshot()
	snapshot.Players["p1"].Gold = 0
	snapshot.Islands["p1"].Buildings["10,10"] = Building{Type: BuildingFort}

	actual := engine.Snapshot()
	if actual.Players["p1"].Gold != StartingGold || len(actual.Islands["p1"].Buildings) != 0 {
		t.Fatal("mutating a snapshot changed engine state")
	}
}

func TestStartIsIdempotentAndStopIsSafe(t *testing.T) {
	config := DefaultConfig()
	config.TurnDuration = 10 * time.Millisecond
	engine, err := NewGameEngineWithConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	engine.Start(ctx)
	engine.Start(ctx)
	t.Cleanup(engine.Stop)

	deadline := time.After(250 * time.Millisecond)
	for engine.Snapshot().Turn < 2 {
		select {
		case <-deadline:
			t.Fatal("turn clock did not advance")
		case <-time.After(time.Millisecond):
		}
	}
	engine.Stop()
	engine.Stop()
}

func TestGameEngineRejectsImpossibleGeometry(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Config)
	}{
		{name: "negative dimension", change: func(config *Config) { config.MapWidth = -1 }},
		{name: "too narrow for two islands", change: func(config *Config) {
			config.MaxPlayers = 2
			config.MapWidth = 2*config.IslandWidth + 2*IslandMargin + MinimumIslandGap - 1
		}},
		{name: "too short for three islands", change: func(config *Config) {
			config.MaxPlayers = 3
			config.MapHeight = 2*config.IslandHeight + 2*IslandMargin + MinimumIslandGap - 1
		}},
		{name: "too many players", change: func(config *Config) { config.MaxPlayers = MaxPlayers + 1 }},
		{name: "non-positive turn duration", change: func(config *Config) { config.TurnDuration = -time.Second }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := DefaultConfig()
			test.change(&config)
			if _, err := NewGameEngineWithConfig(config); !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("expected ErrInvalidConfig, got %v", err)
			}
		})
	}
}

func TestConcurrentSnapshotsAndBuilds(t *testing.T) {
	engine := NewGameEngine()
	if err := engine.AddPlayer("p1", "Player"); err != nil {
		t.Fatal(err)
	}
	var group sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		group.Add(1)
		go func(offset int) {
			defer group.Done()
			for index := 0; index < 50; index++ {
				_ = engine.Build("p1", BuildingFarm, 10+(index+offset)%25, 10+(index*3)%25)
				_ = engine.Snapshot()
			}
		}(worker)
	}
	group.Wait()
}
