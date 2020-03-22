package interpreter

import "github.com/Azer0s/Hummus/parser"

// NodeType a variable type
type NodeType uint8

// Node a variable node
type Node struct {
	Value    interface{}
	NodeType NodeType
}

const (
	// NODETYPE_INT int variable type
	NODETYPE_INT NodeType = 0
	// NODETYPE_FLOAT float variable type
	NODETYPE_FLOAT NodeType = 1
	// NODETYPE_STRING string variable type
	NODETYPE_STRING NodeType = 2
	// NODETYPE_BOOL bool variable type
	NODETYPE_BOOL NodeType = 3
	// NODETYPE_ATOM atom variable type
	NODETYPE_ATOM NodeType = 4
	// NODETYPE_FN function literal
	NODETYPE_FN NodeType = 5
	// NODETYPE_LIST list type
	NODETYPE_LIST NodeType = 6
	// NODETYPE_MAP map type
	NODETYPE_MAP NodeType = 7
	// NODETYPE_STRUCT struct type
	NODETYPE_STRUCT NodeType = 8
)

// FnLiteral a function literal (block)
type FnLiteral struct {
	Parameters []string
	Body       []parser.Node
	Context    map[string]Node
}

// ListNode a list value
type ListNode struct {
	Values []Node
}

// MapNode a map node
type MapNode struct {
	Values map[string]Node
}

// StructDef struct definition
type StructDef struct {
	Parameters []string
}
