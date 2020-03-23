package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	// SYSTEM_MATH the built in math system function
	SYSTEM_MATH string = "--system-do-math!"
	// SYSTEM_MAKE the built in value make system function
	SYSTEM_MAKE string = "--system-do-make!"
	// SYSTEM_IO the built in io function
	SYSTEM_IO string = "--system-do-io!"
	// SYSTEM_COMPARE comparison functions
	SYSTEM_COMPARE string = "--system-do-compare!"
	// SYSTEM_COMPARE_ARITHMETIC arithmetic comparison functions
	SYSTEM_COMPARE_ARITHMETIC string = "--system-do-compare-arithmetic!"
	// SYSTEM_CONVERT conversion functions
	SYSTEM_CONVERT string = "--system-do-convert!"
	// SYSTEM_BOOL boolean algebra functions
	SYSTEM_BOOL string = "--system-do-bool!"
	// SYSTEM_BITWISE bitwise functions
	SYSTEM_BITWISE string = "--system-do-bitwise!"
	// SYSTEM_ENUMERATE enumeration functions
	SYSTEM_ENUMERATE string = "--system-do-enumerate!"
	// SYSTEM_STRING string functions
	SYSTEM_STRING string = "--system-do-strings!"
	// SYSTEM_DEBUG debug functions
	SYSTEM_DEBUG string = "--system-do-debug!"

	// SYSTEM_ENUMERATE_VAL variable where mfr values are stored
	SYSTEM_ENUMERATE_VAL string = "--system-do-enumerate-val"
	// SYSTEM_ACCUMULATE_VAL variable where reduce state is stored
	SYSTEM_ACCUMULATE_VAL string = "--system-do-accumulate-val"
)

var reader = bufio.NewReader(os.Stdin)

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
