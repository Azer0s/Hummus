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
		return Node{
			Value:    -args[1].Value.(int),
			NodeType: NODETYPE_INT,
		}
	}

	if mode == "-" && args[1].NodeType == NODETYPE_FLOAT {
		return Node{
			Value:    -args[1].Value.(float64),
			NodeType: NODETYPE_FLOAT,
		}
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
		return Node{
			Value:    args[1].Value != args[2].Value,
			NodeType: NODETYPE_BOOL,
		}
	case "=":
		return Node{
			Value:    args[1].Value == args[2].Value,
			NodeType: NODETYPE_BOOL,
		}
	case "min":
		if args[1].NodeType != NODETYPE_LIST {
			panic(BUILTIN_COMPARE + " :min only accepts lists!")
		}

		list := args[1].Value.(ListNode).Values
		min := list[0]

		for i := 1; i < len(list); i++ {
			if list[i].Smaller(min) {
				min = list[i]
			}
		}

		return min
	case "max":
		if args[1].NodeType != NODETYPE_LIST {
			panic(BUILTIN_COMPARE + " :max only accepts lists!")
		}

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
		if args[1].NodeType != NODETYPE_BOOL {
			panic(BUILTIN_BOOL + " only accepts bools!")
		}

		return Node{
			Value:    !args[1].Value.(bool),
			NodeType: NODETYPE_BOOL,
		}
	}

	if args[1].NodeType != NODETYPE_BOOL || args[2].NodeType != NODETYPE_BOOL {
		panic(BUILTIN_BOOL + " only accepts bools!")
	}

	switch mode {
	case "and":
		return Node{
			Value:    args[1].Value.(bool) && args[2].Value.(bool),
			NodeType: NODETYPE_BOOL,
		}
	case "or":
		return Node{
			Value:    args[1].Value.(bool) || args[2].Value.(bool),
			NodeType: NODETYPE_BOOL,
		}
	default:
		panic("Unrecognized mode")
	}
}

func builtInBitwise(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables)

	mode := args[0].Value.(string)

	if mode == "not" {
		if args[1].NodeType != NODETYPE_INT {
			panic(BUILTIN_BITWISE + " only accepts ints!")
		}

		return Node{
			Value:    int(^uint(args[1].Value.(int))),
			NodeType: NODETYPE_INT,
		}
	}

	if args[1].NodeType != NODETYPE_INT || args[2].NodeType != NODETYPE_INT {
		panic(BUILTIN_BITWISE + " only accepts ints!")
	}

	switch mode {
	case "and":
		return Node{
			Value:    int(uint(args[1].Value.(int)) & uint(args[2].Value.(int))),
			NodeType: NODETYPE_INT,
		}
	case "or":
		return Node{
			Value:    int(uint(args[1].Value.(int)) | uint(args[2].Value.(int))),
			NodeType: NODETYPE_INT,
		}
	case "shiftl":
		return Node{
			Value:    int(uint(args[1].Value.(int)) << uint(args[2].Value.(int))),
			NodeType: NODETYPE_INT,
		}
	case "shiftr":
		return Node{
			Value:    int(uint(args[1].Value.(int)) >> uint(args[2].Value.(int))),
			NodeType: NODETYPE_INT,
		}
	default:
		panic("Unrecognized mode")
	}
}

func doFloatCalculation(mode string, vals []float64) (node Node) {
	node = Node{
		Value:    vals[0],
		NodeType: NODETYPE_FLOAT,
	}

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
	node = Node{
		Value:    vals[0],
		NodeType: NODETYPE_INT,
	}

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
		return Node{
			Value:    f[0] < f[1],
			NodeType: NODETYPE_BOOL,
		}
	case ">":
		return Node{
			Value:    f[0] > f[1],
			NodeType: NODETYPE_BOOL,
		}
	case "<=":
		return Node{
			Value:    f[0] <= f[1],
			NodeType: NODETYPE_BOOL,
		}
	case ">=":
		return Node{
			Value:    f[0] >= f[1],
			NodeType: NODETYPE_BOOL,
		}
	default:
		panic("Unrecognized mode")
	}
}

func doIntCompare(mode string, i []int) Node {
	switch mode {
	case "<":
		return Node{
			Value:    i[0] < i[1],
			NodeType: NODETYPE_BOOL,
		}
	case ">":
		return Node{
			Value:    i[0] > i[1],
			NodeType: NODETYPE_BOOL,
		}
	case "<=":
		return Node{
			Value:    i[0] <= i[1],
			NodeType: NODETYPE_BOOL,
		}
	case ">=":
		return Node{
			Value:    i[0] >= i[1],
			NodeType: NODETYPE_BOOL,
		}
	default:
		panic("Unrecognized mode")
	}
}
