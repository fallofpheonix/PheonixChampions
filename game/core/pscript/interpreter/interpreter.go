// Package interpreter: Implements the execution engine for P-Script.
// It walks the AST and applies transformations to the provided GameState.
package interpreter

import (
	"phoenix-game/core/pscript/ast"
	"phoenix-game/core/state"
)

// BuiltinFn: A function defined in Go that can be called from P-Script.
type BuiltinFn func(s *state.GameState)

// Interpreter: Holds the environment of built-in functions and executes ASTs.
type Interpreter struct {
	builtins map[string]BuiltinFn
}

// New: Creates a new Interpreter with the given built-in functions.
func New(builtins map[string]BuiltinFn) *Interpreter {
	return &Interpreter{
		builtins: builtins,
	}
}

// Eval: Walks the program and executes its statements against the state.
func (i *Interpreter) Eval(program *ast.Program, s *state.GameState) {
	for _, stmt := range program.Statements {
		i.evalStatement(stmt, s)
	}
}

func (i *Interpreter) evalStatement(stmt ast.Statement, s *state.GameState) {
	if stmt.FunctionDeclaration != nil {
		i.evalFunctionDeclaration(stmt.FunctionDeclaration, s)
	}
	if stmt.CallExpression != nil {
		i.evalCallExpression(stmt.CallExpression, s)
	}
}

func (i *Interpreter) evalFunctionDeclaration(node *ast.FunctionDeclaration, s *state.GameState) {
	// For v1, we only execute the 'main' function immediately.
	if node.Name == "main" {
		for _, stmt := range node.Body {
			i.evalStatement(stmt, s)
		}
	}
}

func (i *Interpreter) evalCallExpression(node *ast.CallExpression, s *state.GameState) {
	if fn, ok := i.builtins[node.Name]; ok {
		fn(s)
	}
}
