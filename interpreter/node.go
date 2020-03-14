package interpreter

type NodeType uint8

type Node struct {
	Value    interface{}
	NodeType NodeType
}

const (
	NODETYPE_INT    NodeType = 0
	NODETYPE_FLOAT  NodeType = 1
	NODETYPE_STRING NodeType = 2
	NODETYPE_BOOL   NodeType = 3
	NODETYPE_ATOM   NodeType = 4
)
