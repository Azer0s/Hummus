package parser

import "github.com/Azer0s/Hummus/lexer"

type NodeType uint8

type Node struct {
	Type      NodeType
	Arguments []Node
	Value     lexer.Token
}

const (
	ACTION_DEF     NodeType = 0
	ACTION_CALL    NodeType = 1
	ACTION_FOR     NodeType = 2
	ACTION_IF      NodeType = 3
	LITERAL_FN     NodeType = 4
	LITERAL_STRING NodeType = 5
)
