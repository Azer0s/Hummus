package interpreter

import (
	"fmt"
	"github.com/Azer0s/Hummus/lexer"
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
			panic("Expected an int as parameter for range!")
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
	case "console-out":
		if args[1].NodeType <= NODETYPE_ATOM {
			fmt.Print(args[1].Value)
		} else {
			panic(SYSTEM_IO + " :console-out only accepts int, float, bool, string or atom!")
		}
	case "console-in":
		t, _ := reader.ReadString('\n')
		return Node{
			Value:    t,
			NodeType: NODETYPE_STRING,
		}
	case "console-clear":
		print("\033[H\033[2J")
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
	case "min":
		if args[1].NodeType != NODETYPE_LIST {
			panic(SYSTEM_COMPARE + " :min only accepts lists!")
		}

		list := args[1].Value.(ListNode).Values
		min := list[0]

		for i := 1; i < len(list); i++ {
			if list[i].smaller(min) {
				min = list[i]
			}
		}

		return min
	case "max":
		if args[1].NodeType != NODETYPE_LIST {
			panic(SYSTEM_COMPARE + " :max only accepts lists!")
		}

		list := args[1].Value.(ListNode).Values
		max := list[0]

		for i := 1; i < len(list); i++ {
			if list[i].bigger(max) {
				max = list[i]
			}
		}

		return max
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

	if args[1].NodeType == NODETYPE_LIST {
		panic(SYSTEM_CONVERT + " doesn't accept lists!")
	}

	switch mode {
	case "string":
		return Node{
			Value:    DumpNode(args[1]),
			NodeType: NODETYPE_STRING,
		}
	case "identity":
		return args[1]
	default:
		panic("Unrecognized mode")
	}
}

func doSystemCallBool(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

	mode := args[0].Value.(string)

	if mode == "not" {
		if args[1].NodeType != NODETYPE_BOOL {
			panic(SYSTEM_BOOL + " only accepts bools!")
		}

		return Node{
			Value:    !args[1].Value.(bool),
			NodeType: NODETYPE_BOOL,
		}
	}

	if args[1].NodeType != NODETYPE_BOOL || args[2].NodeType != NODETYPE_BOOL {
		panic(SYSTEM_BOOL + " only accepts bools!")
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

func doSystemCallBitwise(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

	mode := args[0].Value.(string)

	if mode == "not" {
		if args[1].NodeType != NODETYPE_INT {
			panic(SYSTEM_BITWISE + " only accepts ints!")
		}

		return Node{
			Value:    int(^uint(args[1].Value.(int))),
			NodeType: NODETYPE_INT,
		}
	}

	if args[1].NodeType != NODETYPE_INT || args[2].NodeType != NODETYPE_INT {
		panic(SYSTEM_BITWISE + " only accepts ints!")
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

func doSystemCallEnumerate(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

	mode := args[0].Value.(string)

	if args[1].NodeType != NODETYPE_LIST {
		panic(SYSTEM_ENUMERATE + " expects a list as first argument!")
	}

	if args[2].NodeType != NODETYPE_FN {
		panic(SYSTEM_ENUMERATE + " expects a function as second argument!")
	}

	ctx := make(map[string]Node, 0)
	for k, v := range *variables {
		ctx[k] = v
	}

	const SYSTEM_ENUMERATE_VAL = "--system-do-enumerate-val"
	const SYSTEM_ACCUMULATE_VAL = "--system-do-accumulate-val"

	fn := args[2].Value.(FnLiteral)
	list := args[1].Value.(ListNode)

	switch mode {
	case "each":
		if len(fn.Parameters) != 1 {
			panic("Enumerate each should have one parameter in execution function!")
		}

		for _, value := range list.Values {
			ctx[SYSTEM_ENUMERATE_VAL] = value

			doVariableCall(parser.Node{
				Type: 0,
				Arguments: []parser.Node{
					{
						Type:      parser.IDENTIFIER,
						Arguments: nil,
						Token: lexer.Token{
							Value: SYSTEM_ENUMERATE_VAL,
							Type:  0,
							Line:  0,
						},
					},
				},
				Token: lexer.Token{},
			}, args[2], &ctx)
		}

		return Node{
			Value:    0,
			NodeType: 0,
		}
	case "map":
		if len(fn.Parameters) != 1 {
			panic("Enumerate map should have one parameter in execution function!")
		}

		mapResult := ListNode{Values: make([]Node, 0)}

		for _, value := range list.Values {
			ctx[SYSTEM_ENUMERATE_VAL] = value

			mapResult.Values = append(mapResult.Values, doVariableCall(parser.Node{
				Type: 0,
				Arguments: []parser.Node{
					{
						Type:      parser.IDENTIFIER,
						Arguments: nil,
						Token: lexer.Token{
							Value: SYSTEM_ENUMERATE_VAL,
							Type:  0,
							Line:  0,
						},
					},
				},
				Token: lexer.Token{},
			}, args[2], &ctx))
		}

		return Node{
			Value:    mapResult,
			NodeType: NODETYPE_LIST,
		}
	case "filter":
		if len(fn.Parameters) != 1 {
			panic("Enumerate filter should have one parameter in execution function!")
		}

		mapResult := ListNode{Values: make([]Node, 0)}

		for _, value := range list.Values {
			ctx[SYSTEM_ENUMERATE_VAL] = value

			res := doVariableCall(parser.Node{
				Type: 0,
				Arguments: []parser.Node{
					{
						Type:      parser.IDENTIFIER,
						Arguments: nil,
						Token: lexer.Token{
							Value: SYSTEM_ENUMERATE_VAL,
							Type:  0,
							Line:  0,
						},
					},
				},
				Token: lexer.Token{},
			}, args[2], &ctx)

			if res.NodeType != NODETYPE_BOOL {
				panic("Enumerate filter result must be a bool!")
			}

			if res.Value.(bool) {
				mapResult.Values = append(mapResult.Values, value)
			}
		}

		return Node{
			Value:    mapResult,
			NodeType: NODETYPE_LIST,
		}
	case "reduce":
		if len(fn.Parameters) != 2 {
			panic("Enumerate reduce should have two parameters in execution function!")
		}

		ctx[SYSTEM_ACCUMULATE_VAL] = args[3]

		for _, value := range list.Values {
			ctx[SYSTEM_ENUMERATE_VAL] = value

			res := doVariableCall(parser.Node{
				Type: 0,
				Arguments: []parser.Node{
					{
						Type:      parser.IDENTIFIER,
						Arguments: nil,
						Token: lexer.Token{
							Value: SYSTEM_ENUMERATE_VAL,
							Type:  0,
							Line:  0,
						},
					},
					{
						Type:      parser.IDENTIFIER,
						Arguments: nil,
						Token: lexer.Token{
							Value: SYSTEM_ACCUMULATE_VAL,
							Type:  0,
							Line:  0,
						},
					},
				},
				Token: lexer.Token{},
			}, args[2], &ctx)

			if res.NodeType != args[3].NodeType {
				panic("Enumerate reduce result type must be the same as the init type!")
			}

			ctx[SYSTEM_ACCUMULATE_VAL] = res
		}

		return ctx[SYSTEM_ACCUMULATE_VAL]
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
