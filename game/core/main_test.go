// Package main_test: Verification suite for the simulation core.
// Ensures the Go-Godot bridge contract remains stable.
package main

import (
	"encoding/json"
	"phoenix-game/core/state"
	"testing"
)

func TestGameStateSerialization(t *testing.T) {
	game := state.GameState{
		Tick:   1,
		Agent:  state.Vector2{X: 1000, Y: 2000},
		Goal:   state.Vector2{X: 50000, Y: 0},
		Status: state.StatusActive,
	}

	data, err := json.Marshal(game)
	if err != nil {
		t.Fatalf("Failed to marshal state: %v", err)
	}

	expected := `{"tick":1,"agent":{"x":1000,"y":2000},"goal":{"x":50000,"y":0},"status":"active"}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}
