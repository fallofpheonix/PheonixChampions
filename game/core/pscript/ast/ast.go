// Package ast: Defines the Abstract Syntax Tree for P-Script.
// These are concrete structs representing the hierarchical structure of the code.
// Every node is a concrete struct to ensure type safety and determinism.
package ast

// Program: The root of the AST.
type Program struct {
	Statements []Statement
}

// Statement: A container for different statement types.
// Using pointers to concrete structs instead of interfaces to prevent runtime panics.
type Statement struct {
	FunctionDeclaration *FunctionDeclaration
	CallExpression      *CallExpression
}

// FunctionDeclaration: Represents a named function and its body.
type FunctionDeclaration struct {
	Name string
	Body []Statement
}

// CallExpression: Represents a function invocation.
type CallExpression struct {
	Name string
}
