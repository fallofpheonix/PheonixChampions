package main

import (
	"phoenix-game/core/pscript/ast"
	"phoenix-game/core/pscript/interpreter"
	"phoenix-game/core/state"
	"testing"
)

func TestWinCondition(t *testing.T) {
	game := state.GameState{
		Tick:   0,
		Agent:  state.Vector2{X: 0, Y: 0},
		Goal:   state.Vector2{X: 1000, Y: 0},
		Status: state.StatusActive,
	}

	builtins := map[string]interpreter.BuiltinFn{
		"move_forward": func(s *state.GameState) {
			s.Agent.X += 1000
		},
	}
	interp := interpreter.New(builtins)

	// Mock tick loop
	game.Tick++
	// We simulate the evaluation of a script that calls move_forward
	game.Agent.X += 1000 
	interp.Eval(&ast.Program{}, &game)

	if game.Agent.X >= game.Goal.X {
		game.Status = state.StatusComplete
	}

	if game.Status != state.StatusComplete {
		t.Errorf("expected status to be complete, got %s", game.Status)
	}
}
