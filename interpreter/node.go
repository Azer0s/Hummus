package interpreter

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
)
