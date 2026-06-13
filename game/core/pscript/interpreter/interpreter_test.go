// Package interpreter_test: Verification for the P-Script execution engine.
// Ensures that AST traversal correctly mutates the provided GameState.
package interpreter

import (
	"phoenix-game/core/pscript/lexer"
	"phoenix-game/core/pscript/parser"
	"phoenix-game/core/state"
	"testing"
)

func TestEval(t *testing.T) {
	input := `fn main() { move_forward() }`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	game := &state.GameState{
		Tick:  0,
		Agent: state.Vector2{X: 0, Y: 0},
	}

	builtins := map[string]BuiltinFn{
		"move_forward": func(s *state.GameState) {
			s.Agent.X += 1000
		},
	}

	interp := New(builtins)
	interp.Eval(program, game)

	if game.Agent.X != 1000 {
		t.Errorf("expected Agent.X to be 1000, got %d", game.Agent.X)
	}
}
