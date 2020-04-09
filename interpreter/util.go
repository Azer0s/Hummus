package interpreter

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ordinalize(num int) string {
	ordinalDictionary := map[int]string{
		0: "th",
		1: "st",
		2: "nd",
		3: "rd",
		4: "th",
		5: "th",
		6: "th",
		7: "th",
		8: "th",
		9: "th",
	}

	floatNum := math.Abs(float64(num))
	positiveNum := int(floatNum)

	if ((positiveNum % 100) >= 11) && ((positiveNum % 100) <= 13) {
		return strconv.Itoa(num) + "th"
	}

	return strconv.Itoa(num) + ordinalDictionary[positiveNum]

}

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

// CopyVariableState copy the variable state to another map
func CopyVariableState(variables, ctx *map[string]Node) {
	for k, v := range *variables {
		(*ctx)[k] = v
	}
}

// EnsureTypes ensures the variable type of the parameter for a native function is valid
func EnsureTypes(val *[]Node, nth int, nt []NodeType, who string) {
	valid := false
	for _, nodeType := range nt {
		valid = valid || (*val)[nth].NodeType == nodeType
	}

	if !valid {
		stringTypes := make([]string, 0)

		for _, nodeType := range nt {
			stringTypes = append(stringTypes, nodeType.String())
		}

		panic(who + " expects " + strings.Join(stringTypes, " or ") + " as the " + ordinalize(nth) + " argument!")
	}
}

// EnsureType ensures the variable type of the parameter for a native function is valid
func EnsureType(val *[]Node, nth int, nt NodeType, who string) {
	if (*val)[nth].NodeType != nt {
		panic(who + " expects " + nt.String() + " as the " + ordinalize(nth) + " argument!")
	}
}

// EnsureSingleType ensures the variable type of the parameter for a native function is valid
func EnsureSingleType(val *Node, nth int, nt NodeType, who string) {
	if (*val).NodeType != nt {
		panic(who + " expects " + nt.String() + " as the " + ordinalize(nth) + " argument!")
	}
}
