package interpreter

import (
	"fmt"
	"github.com/Azer0s/Hummus/parser"
)

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

// Smaller < operator for Node
func (node *Node) Smaller(compareTo Node) bool {
	if node.NodeType != compareTo.NodeType {
		panic("Can't compare nodes of two different types!")
	}

	switch node.NodeType {
	case NODETYPE_INT:
		return node.Value.(int) < compareTo.Value.(int)
	case NODETYPE_FLOAT:
		return node.Value.(float64) < compareTo.Value.(float64)
	case NODETYPE_STRING:
		return node.Value.(string) < compareTo.Value.(string)
	case NODETYPE_ATOM:
		return node.Value.(string) < compareTo.Value.(string)
	default:
		panic(fmt.Sprintf("Nodetype %d cannot be compared!", node.NodeType))
	}
}

// Bigger > operator for Node
func (node *Node) Bigger(compareTo Node) bool {
	if node.NodeType != compareTo.NodeType {
		panic("Can't compare nodes of two different types!")
	}

	switch node.NodeType {
	case NODETYPE_INT:
		return node.Value.(int) > compareTo.Value.(int)
	case NODETYPE_FLOAT:
		return node.Value.(float64) > compareTo.Value.(float64)
	case NODETYPE_STRING:
		return node.Value.(string) > compareTo.Value.(string)
	case NODETYPE_ATOM:
		return node.Value.(string) > compareTo.Value.(string)
	default:
		panic(fmt.Sprintf("Nodetype %d cannot be compared!", node.NodeType))
	}
}

// OptionalNode return an optional node
func OptionalNode(val interface{}, nodeType NodeType, err bool) Node {
	return Node{
		Value: MapNode{Values: map[string]Node{
			"value": {
				Value:    val,
				NodeType: nodeType,
			},
			"error": {
				Value:    err,
				NodeType: NODETYPE_BOOL,
			},
		}},
		NodeType: NODETYPE_MAP,
	}
}
