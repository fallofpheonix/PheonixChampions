// Package parser_test: Verification for the P-Script parser.
// Ensures that the move_forward() program is correctly transformed into an AST.
package parser

import (
	"phoenix-game/core/pscript/lexer"
	"testing"
)

func TestParseProgram(t *testing.T) {
	input := `fn main() { move_forward() }`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt := program.Statements[0]
	if stmt.FunctionDeclaration == nil {
		t.Fatalf("stmt is not a FunctionDeclaration")
	}

	if stmt.FunctionDeclaration.Name != "main" {
		t.Errorf("stmt.Name not 'main'. got=%q", stmt.FunctionDeclaration.Name)
	}

	if len(stmt.FunctionDeclaration.Body) != 1 {
		t.Fatalf("body does not contain 1 statement. got=%d",
			len(stmt.FunctionDeclaration.Body))
	}

	bodyStmt := stmt.FunctionDeclaration.Body[0]
	if bodyStmt.CallExpression == nil {
		t.Fatalf("bodyStmt is not a CallExpression")
	}

	if bodyStmt.CallExpression.Name != "move_forward" {
		t.Errorf("bodyStmt.Name not 'move_forward'. got=%q",
			bodyStmt.CallExpression.Name)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
