// Package parser: Implements the recursive-descent parser for P-Script.
// It transforms a stream of tokens into a concrete AST.
package parser

import (
	"fmt"
	"phoenix-game/core/pscript/ast"
	"phoenix-game/core/pscript/lexer"
	"phoenix-game/core/pscript/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, *stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() *ast.Statement {
	switch p.curToken.Type {
	case token.FN:
		return p.parseFunctionDeclaration()
	case token.IDENT:
		if p.peekToken.Type == token.LPAREN {
			return p.parseCallExpression()
		}
	}
	return nil
}

func (p *Parser) parseFunctionDeclaration() *ast.Statement {
	stmt := &ast.Statement{
		FunctionDeclaration: &ast.FunctionDeclaration{},
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.FunctionDeclaration.Name = p.curToken.Literal

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.FunctionDeclaration.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseBlockStatement() []ast.Statement {
	statements := []ast.Statement{}

	p.nextToken()

	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, *stmt)
		}
		p.nextToken()
	}

	return statements
}

func (p *Parser) parseCallExpression() *ast.Statement {
	stmt := &ast.Statement{
		CallExpression: &ast.CallExpression{},
	}

	stmt.CallExpression.Name = p.curToken.Literal

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
