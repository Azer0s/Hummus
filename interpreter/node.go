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
	// NODETYPE_MACRO macro type
	NODETYPE_MACRO NodeType = 9
)

func (nt NodeType) String() string {
	t := ""

	switch nt {
	case NODETYPE_INT:
		t = "int"
	case NODETYPE_FLOAT:
		t = "float"
	case NODETYPE_STRING:
		t = "string"
	case NODETYPE_BOOL:
		t = "bool"
	case NODETYPE_ATOM:
		t = "atom"
	case NODETYPE_FN:
		t = "fn"
	case NODETYPE_LIST:
		t = "list"
	case NODETYPE_MAP:
		t = "map"
	case NODETYPE_STRUCT:
		t = "struct"
	}

	return t
}

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

type MacroParameter struct {
	Parameter string
	Literal   bool
}

type MacroDef struct {
	Parameters []MacroParameter
	Body       []parser.Node
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

// OptionNode return an optional node
func OptionNode(val Node, err bool) Node {
	return NodeMap(map[string]Node{
		"value": val,
		"error": BoolNode(err),
	})
}

// Nothing represents an empty return (int 0)
var Nothing = Node{
	Value:    0,
	NodeType: 0,
}

// NodeList returns an interpreter.Node from a []Node
func NodeList(val []Node) Node {
	return Node{
		Value:    ListNode{Values: val},
		NodeType: NODETYPE_LIST,
	}
}

// NodeMap returns an interpreter.Node from a map[string]Node
func NodeMap(val map[string]Node) Node {
	return Node{
		Value:    MapNode{Values: val},
		NodeType: NODETYPE_MAP,
	}
}

// BoolNode returns an interpreter.Node from a bool
func BoolNode(val bool) Node {
	return Node{
		Value:    val,
		NodeType: NODETYPE_BOOL,
	}
}

// StringNode returns an interpreter.Node from a string
func StringNode(val string) Node {
	return Node{
		Value:    val,
		NodeType: NODETYPE_STRING,
	}
}

// AtomNode returns an interpreter.Node from a string
func AtomNode(val string) Node {
	return Node{
		Value:    val,
		NodeType: NODETYPE_ATOM,
	}
}

// IntNode returns an interpreter.Node from an int
func IntNode(val int) Node {
	return Node{
		Value:    val,
		NodeType: NODETYPE_INT,
	}
}

// FloatNode returns an interpreter.Node from a float64
func FloatNode(val float64) Node {
	return Node{
		Value:    val,
		NodeType: NODETYPE_FLOAT,
	}
}
