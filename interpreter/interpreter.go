package interpreter

import (
	"github.com/Azer0s/Hummus/parser"
	"strconv"
)

func defineVariable(node parser.Node, variables *map[string]Node) Node {
	name := node.Arguments[0].Token.Value
	variable := Node{
		Value:    nil,
		NodeType: 0,
	}

	value := node.Arguments[1]

	switch value.Type {
	case parser.LITERAL_INT:
		variable.NodeType = NODETYPE_INT
		i, _ := strconv.Atoi(value.Token.Value)
		variable.Value = i
		break
	case parser.LITERAL_FLOAT:
		variable.NodeType = NODETYPE_FLOAT
		f, _ := strconv.ParseFloat(value.Token.Value, 64)
		variable.Value = f
		break
	case parser.LITERAL_STRING:
		variable.NodeType = NODETYPE_STRING
		variable.Value = value.Token.Value
		break
	case parser.LITERAL_BOOL:
		variable.NodeType = NODETYPE_BOOL
		b, _ := strconv.ParseBool(value.Token.Value)
		variable.Value = b
		break
	case parser.LITERAL_ATOM:
		variable.NodeType = NODETYPE_ATOM
		variable.Value = value.Token.Value
		break
	}

	(*variables)[name] = variable
	return variable
}

// Run run an AST
func Run(nodes []parser.Node, variables map[string]Node) (returnVal Node) {
	for _, node := range nodes {
		switch node.Type {
		case parser.ACTION_DEF:
			returnVal = defineVariable(node, &variables)
			break
		}
	}

	return
}
