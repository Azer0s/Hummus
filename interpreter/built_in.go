package interpreter

import (
	"fmt"
	"github.com/Azer0s/Hummus/parser"
)

const (
	// BUILTIN_MATH the built in math system function
	BUILTIN_MATH string = "--builtin-do-math!"
	// BUILTIN_COMPARE comparison functions
	BUILTIN_COMPARE string = "--builtin-do-compare!"
	// BUILTIN_COMPARE_ARITHMETIC arithmetic comparison functions
	BUILTIN_COMPARE_ARITHMETIC string = "--builtin-do-compare-arithmetic!"
	// BUILTIN_BOOL boolean algebra functions
	BUILTIN_BOOL string = "--builtin-do-bool!"
	// BUILTIN_BITWISE bitwise functions
	BUILTIN_BITWISE string = "--builtin-do-bitwise!"
)

func builtInMath(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables)

	mode := args[0].Value.(string)

	if mode == "-" && args[1].NodeType == NODETYPE_INT {
		return IntNode(-args[1].Value.(int))
	}

	if mode == "-" && args[1].NodeType == NODETYPE_FLOAT {
		return FloatNode(-args[1].Value.(float64))
	}

	vals := args[1].Value.(ListNode)

	if len(vals.Values) < 2 {
		panic(fmt.Sprintf("Arithmetic operations expect at least 2 arguments! (line %d)", node.Token.Line))
	}

	//If any arg is a float, we switch to float mode
	//Only ints and floats allowed

	floatMode := false

	for _, value := range vals.Values {
		if value.NodeType != NODETYPE_FLOAT && value.NodeType != NODETYPE_INT {
			panic(fmt.Sprintf("Only float or int allowed for arithmetic operations! (line %d)", node.Token.Line))
		}

		if value.NodeType == NODETYPE_FLOAT {
			floatMode = true
			break
		}
	}

	if floatMode {
		f := make([]float64, 0)

		for _, value := range vals.Values {
			if value.NodeType == NODETYPE_INT {
				f = append(f, float64(value.Value.(int)))
			} else {
				f = append(f, value.Value.(float64))
			}
		}

		return doFloatCalculation(mode, f)
	}

	i := make([]int, 0)

	for _, value := range vals.Values {
		i = append(i, value.Value.(int))
	}

	return doIntCalculation(mode, i)
}

func builtInCompare(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables)

	mode := args[0].Value.(string)

	switch mode {
	case "/=":
		return BoolNode(args[1].Value != args[2].Value)
	case "=":
		return BoolNode(args[1].Value == args[2].Value)
	case "min":
		EnsureType(&args, 1, NODETYPE_LIST, BUILTIN_COMPARE+" :min")
		list := args[1].Value.(ListNode).Values

		min := list[0]

		for i := 1; i < len(list); i++ {
			if list[i].Smaller(min) {
				min = list[i]
			}
		}

		return min
	case "max":
		EnsureType(&args, 1, NODETYPE_LIST, BUILTIN_COMPARE+" :max")
		list := args[1].Value.(ListNode).Values

		max := list[0]

		for i := 1; i < len(list); i++ {
			if list[i].Bigger(max) {
				max = list[i]
			}
		}

		return max
	default:
		panic("Unrecognized mode")
	}
}

func builtInCompareArithmetic(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables)

	mode := args[0].Value.(string)
	//If any arg is a float, we switch to float mode
	//Only ints and floats allowed

	floatMode := false

	for _, value := range args[1:] {
		if value.NodeType != NODETYPE_FLOAT && value.NodeType != NODETYPE_INT {
			panic(fmt.Sprintf("Only float or int allowed for arithmetic operations! (line %d)", node.Token.Line))
		}

		if value.NodeType == NODETYPE_FLOAT {
			floatMode = true
			break
		}
	}

	if floatMode {
		f := make([]float64, 0)

		for _, value := range args[1:] {
			if value.NodeType == NODETYPE_INT {
				f = append(f, float64(value.Value.(int)))
			} else {
				f = append(f, value.Value.(float64))
			}
		}

		return doFloatCompare(mode, f)
	}

	i := make([]int, 0)

	for _, value := range args[1:] {
		i = append(i, value.Value.(int))
	}

	return doIntCompare(mode, i)
}

func builtInBool(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables)

	mode := args[0].Value.(string)

	if mode == "not" {
		EnsureType(&args, 1, NODETYPE_BOOL, BUILTIN_BOOL)

		return BoolNode(!args[1].Value.(bool))
	}

	EnsureType(&args, 1, NODETYPE_BOOL, BUILTIN_BOOL)
	EnsureType(&args, 2, NODETYPE_BOOL, BUILTIN_BOOL)

	switch mode {
	case "and":
		return BoolNode(args[1].Value.(bool) && args[2].Value.(bool))
	case "or":
		return BoolNode(args[1].Value.(bool) || args[2].Value.(bool))
	default:
		panic("Unrecognized mode")
	}
}

func builtInBitwise(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables)

	mode := args[0].Value.(string)

	if mode == "not" {
		EnsureType(&args, 1, NODETYPE_INT, BUILTIN_BITWISE)

		return IntNode(int(^uint(args[1].Value.(int))))
	}

	EnsureType(&args, 1, NODETYPE_INT, BUILTIN_BITWISE)
	EnsureType(&args, 2, NODETYPE_INT, BUILTIN_BITWISE)

	switch mode {
	case "and":
		return IntNode(int(uint(args[1].Value.(int)) & uint(args[2].Value.(int))))
	case "or":
		return IntNode(int(uint(args[1].Value.(int)) | uint(args[2].Value.(int))))
	case "shiftl":
		return IntNode(int(uint(args[1].Value.(int)) << uint(args[2].Value.(int))))
	case "shiftr":
		return IntNode(int(uint(args[1].Value.(int)) >> uint(args[2].Value.(int))))
	default:
		panic("Unrecognized mode")
	}
}

func doFloatCalculation(mode string, vals []float64) (node Node) {
	node = FloatNode(vals[0])

	switch mode {
	case "*":
		for i := 1; i < len(vals); i++ {
			node.Value = node.Value.(float64) * vals[i]
		}
		break
	case "/":
		for i := 1; i < len(vals); i++ {
			node.Value = node.Value.(float64) / vals[i]
		}
		break
	case "+":
		for i := 1; i < len(vals); i++ {
			node.Value = node.Value.(float64) + vals[i]
		}
		break
	case "-":
		for i := 1; i < len(vals); i++ {
			node.Value = node.Value.(float64) - vals[i]
		}
		break
	}

	return
}

func doIntCalculation(mode string, vals []int) (node Node) {
	node = IntNode(vals[0])

	switch mode {
	case "*":
		for i := 1; i < len(vals); i++ {
			node.Value = node.Value.(int) * vals[i]
		}
		break
	case "/":
		for i := 1; i < len(vals); i++ {
			node.Value = node.Value.(int) / vals[i]
		}
		break
	case "+":
		for i := 1; i < len(vals); i++ {
			node.Value = node.Value.(int) + vals[i]
		}
		break
	case "-":
		for i := 1; i < len(vals); i++ {
			node.Value = node.Value.(int) - vals[i]
		}
		break
	}

	return
}

func doFloatCompare(mode string, f []float64) Node {
	switch mode {
	case "<":
		return BoolNode(f[0] < f[1])
	case ">":
		return BoolNode(f[0] > f[1])
	case "<=":
		return BoolNode(f[0] <= f[1])
	case ">=":
		return BoolNode(f[0] >= f[1])
	default:
		panic("Unrecognized mode")
	}
}

func doIntCompare(mode string, i []int) Node {
	switch mode {
	case "<":
		return BoolNode(i[0] < i[1])
	case ">":
		return BoolNode(i[0] > i[1])
	case "<=":
		return BoolNode(i[0] <= i[1])
	case ">=":
		return BoolNode(i[0] >= i[1])
	default:
		panic("Unrecognized mode")
	}
}
