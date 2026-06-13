// Package main: The authoritative simulation core for Phoenix Champions.
// This process acts as the "Brains" of the game, managing state and logic.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"phoenix-game/core/pscript/interpreter"
	"phoenix-game/core/pscript/lexer"
	"phoenix-game/core/pscript/parser"
	"phoenix-game/core/state"
	"time"
)

func main() {
	// 1. Initialize State
	game := state.GameState{
		Tick:   0,
		Agent:  state.Vector2{X: 0, Y: 0},
		Goal:   state.Vector2{X: 50000, Y: 0}, // 500 pixels at 100u/px
		Status: state.StatusActive,
	}

	// 2. Load and Parse P-Script
	scriptPath := os.Getenv("PHX_SCRIPT_PATH")
	if scriptPath == "" {
		scriptPath = "scripts/agent.ps" // Fallback
	}
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading script: %v\n", err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Fprintf(os.Stderr, "Parser errors in %s:\n", scriptPath)
		for _, msg := range p.Errors() {
			fmt.Fprintf(os.Stderr, "  - %s\n", msg)
		}
		os.Exit(1)
	}

	// 3. Setup Interpreter and Builtins
	builtins := map[string]interpreter.BuiltinFn{
		"move_forward": func(s *state.GameState) {
			// Using fixed-point math: 1 pixel = 100 units
			// Moving 10 pixels = 1000 units
			s.Agent.X += 1000
		},
	}
	interp := interpreter.New(builtins)

	// 4. Main Simulation Loop
	ticker := time.NewTicker(100 * time.Millisecond) // Faster tick for smoother feeling
	defer ticker.Stop()

	fmt.Println("--- Phoenix Game Core Started (P-Script Driven) ---")

	var lastMod time.Time

	for {
		select {
		case <-ticker.C:
			// 4.1 Check for Script Reload
			info, err := os.Stat(scriptPath)
			if err == nil && info.ModTime().After(lastMod) {
				fmt.Fprintf(os.Stderr, "Reloading script...\n")
				content, err := os.ReadFile(scriptPath)
				if err == nil {
					l := lexer.New(string(content))
					p := parser.New(l)
					newProgram := p.ParseProgram()
					if len(p.Errors()) == 0 {
						program = newProgram
						lastMod = info.ModTime()
						// Reset simulation
						game.Tick = 0
						game.Agent = state.Vector2{X: 0, Y: 0}
						game.Status = state.StatusActive
					}
				}
			}

			if game.Status == state.StatusComplete {
				continue // Wait for reload
			}

			game.Tick++

			// Authoritative logic execution
			interp.Eval(program, &game)

			// 5. Win Condition Check
			if game.Agent.X >= game.Goal.X && game.Agent.Y >= game.Goal.Y {
				game.Status = state.StatusComplete
			}

			// Emit state to bridge
			data, _ := json.Marshal(game)
			fmt.Println(string(data))
		}
	}
}
