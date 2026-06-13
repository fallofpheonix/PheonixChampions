// Package lexer_test: Verification for the P-Script scanner.
// Validates that the move_forward() program tokenizes correctly.
package lexer

import (
	"phoenix-game/core/pscript/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `fn main() { move_forward() }`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.FN, "fn"},
		{token.IDENT, "main"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "move_forward"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
