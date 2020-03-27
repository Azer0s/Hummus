package interpreter

import (
	"fmt"
	"strings"
)

// DumpNode returns the string representation of a node
func DumpNode(node Node) string {
	ret := ""

	if node.NodeType == NODETYPE_LIST {
		ret += "("

		for _, value := range node.Value.(ListNode).Values {
			ret += DumpNode(value) + " "
		}

		ret = strings.TrimSuffix(ret, " ")
		ret += ")"
	} else if node.NodeType == NODETYPE_MAP {
		ret += "("

		for k, v := range node.Value.(MapNode).Values {
			ret += fmt.Sprintf("%s => %s ", k, DumpNode(v))
		}

		ret = strings.TrimSuffix(ret, " ")
		ret += ")"
	} else if node.NodeType == NODETYPE_FN {
		ret += "[fn "

		for _, parameter := range node.Value.(FnLiteral).Parameters {
			ret += parameter + " "
		}

		ret = strings.TrimSuffix(ret, " ")
		ret += "]"
	} else if node.NodeType == NODETYPE_STRUCT {
		ret += "[struct "

		for _, parameter := range node.Value.(StructDef).Parameters {
			ret += parameter + " "
		}

		ret = strings.TrimSuffix(ret, " ")
		ret += "]"
	} else {
		ret = fmt.Sprintf("%v", node.Value)
	}

	return ret
}
