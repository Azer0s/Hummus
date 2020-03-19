package parser

import "github.com/Azer0s/Hummus/lexer"

// NodeType a parser node-type
type NodeType uint8

// Node a parser node
type Node struct {
	Type      NodeType
	Arguments []Node
	Token     lexer.Token
}

const (
	// ACTION_DEF define
	ACTION_DEF NodeType = 0
	// ACTION_CALL call function
	ACTION_CALL NodeType = 1
	// ACTION_FOR for loop
	ACTION_FOR NodeType = 2
	// ACTION_WHILE for loop
	ACTION_WHILE NodeType = 15
	// ACTION_IF if condition
	ACTION_IF NodeType = 3
	// ACTION_MAP_ACCESS map access
	ACTION_MAP_ACCESS = 4
	// ACTION_MAP map creation
	ACTION_MAP = 14
	// ACTION_BRANCH branch statement
	ACTION_BRANCH = 13
	// IDENTIFIER identifier
	IDENTIFIER NodeType = 5
	// LITERAL_FN function literal
	LITERAL_FN NodeType = 6
	// LITERAL_STRING string literal
	LITERAL_STRING NodeType = 7
	// LITERAL_INT int literal
	LITERAL_INT NodeType = 8
	// LITERAL_FLOAT float literal
	LITERAL_FLOAT NodeType = 9
	// LITERAL_BOOL boolean literal
	LITERAL_BOOL NodeType = 10
	// LITERAL_ATOM atom literal
	LITERAL_ATOM NodeType = 11
	// PARAMETERS a parameter node
	PARAMETERS NodeType = 12
	// STRUCT_DEF struct definition
	STRUCT_DEF NodeType = 16
)
