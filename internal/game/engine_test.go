package game

import (
	"testing"
	"time"
)

func TestGameEngine_AddPlayer(t *testing.T) {
	ge := NewGameEngine()

	ge.AddPlayer("P1")
	ge.AddPlayer("P2")

	if len(ge.State.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(ge.State.Players))
	}

	p1Island := ge.State.Islands["P1"]
	if p1Island.X != 10 || p1Island.Y != 10 {
		t.Errorf("P1 Island at wrong location: %d,%d", p1Island.X, p1Island.Y)
	}

	p2Island := ge.State.Islands["P2"]
	if p2Island.X != 165 || p2Island.Y != 10 {
		t.Errorf("P2 Island at wrong location: %d,%d", p2Island.X, p2Island.Y)
	}
}

func TestGameEngine_HandleBuild(t *testing.T) {
	ge := NewGameEngine()
	ge.AddPlayer("P1")

	// Test 1: Valid Build (House at 10,10 - top left of island)
	success := ge.HandleBuild("P1", BuildingHouse, 10, 10)
	if !success {
		t.Error("Expected successful build for valid coordinates and funds")
	}

	// Test 2: Dimensions Check (Build at 9,9 - outside island)
	success = ge.HandleBuild("P1", BuildingHouse, 9, 9)
	if success {
		t.Error("Expected failed build for start coordinates outside bounds")
	}

	success = ge.HandleBuild("P1", BuildingHouse, 35, 35) // 10+25 = 35, so 35 is out of bounds (0-indexed logic usually, but here 10-34 is valid range if width 25)
	if success {
		t.Error("Expected failed build for end coordinates outside bounds")
	}

	// Test 3: Insufficient Funds
	// CostHouse is 150. Initial Gold is 1000.
	// Buy 6 more houses (Total 7 * 150 = 1050 > 1000)
	for i := 0; i < 5; i++ {
		ge.HandleBuild("P1", BuildingHouse, 11+i, 10)
	}
	// Gold should be 1000 - 150 - 5*150 = 1000 - 900 = 100
	if ge.State.Players["P1"].Gold != 100 {
		t.Errorf("Expected 100 Gold, got %d", ge.State.Players["P1"].Gold)
	}

	// Try to build one more (Cost 150) -> Fail
	success = ge.HandleBuild("P1", BuildingHouse, 20, 20)
	if success {
		t.Error("Expected failed build due to insufficient funds")
	}
}

func TestGameEngine_TurnLogic(t *testing.T) {
	ge := NewGameEngine()
	ge.AddPlayer("P1")

	// Build a Factory (Cost 200, Income 20)
	// P1 Gold: 1000 -> 800
	ge.HandleBuild("P1", BuildingFactory, 10, 10)

	// Manually trigger turn
	ge.processTurn()

	// Gold should be 800 + 20 = 820
	// Wait, processTurn adds income.

	if ge.State.Players["P1"].Gold != 820 {
		t.Errorf("Expected 820 Gold after turn, got %d", ge.State.Players["P1"].Gold)
	}

	if ge.State.Turn != 2 {
		t.Errorf("Expected Turn 2, got %d", ge.State.Turn)
	}
}

func TestGameEngine_Concurrency(t *testing.T) {
	// Simple test to ensure no panic during concurrent access
	ge := NewGameEngine()
	ge.AddPlayer("P1")
	ge.Start()
	defer func() { ge.stop <- true }()

	go func() {
		for i := 0; i < 100; i++ {
			ge.HandleBuild("P1", BuildingHouse, 10, 10) // Will fail mostly but hits mutex
		}
	}()

	time.Sleep(100 * time.Millisecond)
	// If no panic, we good.
}
