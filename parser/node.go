package parser

import "github.com/Azer0s/Hummus/lexer"

type NodeType uint8

type Node struct {
	Type      NodeType
	Arguments []Node
	Token     lexer.Token
}

const (
	ACTION_DEF     NodeType = 0
	ACTION_CALL    NodeType = 1
	ACTION_FOR     NodeType = 2
	ACTION_IF      NodeType = 3
	IDENTIFIER     NodeType = 4
	LITERAL_FN     NodeType = 5
	LITERAL_STRING NodeType = 6
	LITERAL_INT    NodeType = 7
	LITERAL_FLOAT  NodeType = 8
	LITERAL_BOOl   NodeType = 9
	LITERAL_ATOM   NodeType = 10
)
