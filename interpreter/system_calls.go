package interpreter

import (
	"fmt"
	"github.com/Azer0s/Hummus/parser"
)

func doSystemCallMath(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

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

func doSystemCallMake(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

	mode := args[0].Value.(string)

	switch mode {
	case "list":
		if args[1].NodeType == NODETYPE_LIST {
			return Node{
				Value:    args[1].Value,
				NodeType: NODETYPE_LIST,
			}
		}

		return Node{
			Value:    ListNode{Values: args[1:]},
			NodeType: NODETYPE_LIST,
		}
	case "range":
		from := args[1]
		to := args[2]

		if from.NodeType != NODETYPE_INT || to.NodeType != NODETYPE_INT {
			panic("Expected an int as parameter for ")
		}

		f := from.Value.(int)
		t := to.Value.(int)

		list := ListNode{Values: make([]Node, 0)}

		if f > t {
			for i := t; i >= t; i-- {
				list.Values = append(list.Values, Node{
					Value:    i,
					NodeType: NODETYPE_INT,
				})
			}
		} else {
			for i := f; i <= t; i++ {
				list.Values = append(list.Values, Node{
					Value:    i,
					NodeType: NODETYPE_INT,
				})
			}
		}

		return Node{
			Value:    list,
			NodeType: NODETYPE_LIST,
		}

	default:
		panic("Unrecognized mode")
	}
}

func doSystemCallIo(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

	mode := args[0].Value.(string)

	switch mode {
	case "console-print":
		if args[1].NodeType == NODETYPE_LIST {
			panic(SYSTEM_IO + " doesn't accept lists!")
		} else {
			fmt.Print(args[1].Value)
		}
	default:
		panic("Unrecognized mode")
	}

	return Node{
		Value:    0,
		NodeType: 0,
	}
}

func doSystemCallCompare(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

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
	default:
		panic("Unrecognized mode")
	}
}

func doSystemCallCompareArithmetic(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

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

func doSystemCallConvert(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

	mode := args[0].Value.(string)

	switch mode {
	case "string":
		if args[1].NodeType == NODETYPE_LIST {
			panic(SYSTEM_CONVERT + " doesn't accept lists!")
		}

		return Node{
			Value:    fmt.Sprintf("%v", args[1].Value),
			NodeType: NODETYPE_STRING,
		}
	case "identity":
		return args[1]
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
