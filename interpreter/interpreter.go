package interpreter

import (
	"fmt"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
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
	// USE include function
	USE string = "use"
	// MAP_ACCESS map access function
	MAP_ACCESS string = "[]"
)

var globalFns = make(map[string]Node, 0)

func getValueFromNode(parserNode parser.Node, variables *map[string]Node) (node Node) {
	node = Node{
		Value:    0,
		NodeType: 0,
	}

	switch parserNode.Type {
	case parser.LITERAL_INT:
		node.NodeType = NODETYPE_INT
		i, _ := strconv.Atoi(parserNode.Token.Value)
		node.Value = i
		break
	case parser.LITERAL_FLOAT:
		node.NodeType = NODETYPE_FLOAT
		f, _ := strconv.ParseFloat(parserNode.Token.Value, 64)
		node.Value = f
		break
	case parser.LITERAL_STRING:
		node.NodeType = NODETYPE_STRING
		node.Value = parserNode.Token.Value
		break
	case parser.LITERAL_BOOL:
		node.NodeType = NODETYPE_BOOL
		b, _ := strconv.ParseBool(parserNode.Token.Value)
		node.Value = b
		break
	case parser.LITERAL_ATOM:
		node.NodeType = NODETYPE_ATOM
		node.Value = parserNode.Token.Value
		break
	case parser.LITERAL_FN:
		node.NodeType = NODETYPE_FN

		parameters := make([]string, 0)
		for _, argument := range parserNode.Arguments[0].Arguments {
			parameters = append(parameters, argument.Token.Value)
		}

		node.Value = FnLiteral{
			Parameters: parameters,
			Body:       parserNode.Arguments[1:],
		}
		break
	case parser.IDENTIFIER:
		node = (*variables)[parserNode.Token.Value]
		break
	case parser.ACTION_CALL:
		node = doCall(parserNode, variables)
		break
	case parser.ACTION_IF:
		node = doIf(parserNode, variables)
		break

	case parser.ACTION_MAP:
		node = createMap(parserNode, variables)
		break

	case parser.ACTION_MAP_ACCESS:
		node = accessMap(parserNode, variables)
		break
		//TODO: Map access, struct
	}

	return
}

func accessMap(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

	if args[0].NodeType != NODETYPE_ATOM {
		panic(fmt.Sprintf("First argument in map access must be an atom! (line %d)", node.Arguments[0].Token.Line))
	}

	if args[1].NodeType != NODETYPE_MAP {
		panic(fmt.Sprintf("Second argument in map access must be a list! (line %d)", node.Arguments[0].Token.Line))
	}

	return args[1].Value.(MapNode).Values[args[0].Value.(string)]
}

func createMap(node parser.Node, variables *map[string]Node) Node {
	mapNode := MapNode{Values: make(map[string]Node, 0)}

	for _, argument := range node.Arguments {
		mapNode.Values[argument.Arguments[0].Token.Value] = getValueFromNode(argument.Arguments[1], variables)
	}

	return Node{
		Value:    mapNode,
		NodeType: NODETYPE_MAP,
	}
}

func getStructDef(node parser.Node) StructDef {
	structDef := StructDef{Parameters: make([]string, 0)}

	for _, argument := range node.Arguments[1].Arguments {
		structDef.Parameters = append(structDef.Parameters, argument.Token.Value)
	}

	return structDef
}

func defineVariable(node parser.Node, variables *map[string]Node) Node {
	name := node.Arguments[0].Token.Value
	variable := Node{
		Value:    0,
		NodeType: 0,
	}

	//I am doing this here, just like I do macros here because you can only define macros and structs in `def`
	if node.Arguments[1].Type == parser.STRUCT_DEF {
		(*variables)[name] = Node{
			Value:    getStructDef(node),
			NodeType: NODETYPE_STRUCT,
		}
	} else {
		(*variables)[name] = getValueFromNode(node.Arguments[1], variables)
	}

	return variable
}

func doUse(node parser.Node) {
	if node.Arguments[0].Type != parser.LITERAL_ATOM {
		panic(fmt.Sprintf("\"use\" only accepts atoms as parameter! (line %d)", node.Token.Line))
	}

	stdlib, _ := regexp.Compile("<([^>]+)>")

	//TODO: Relative path

	file := node.Arguments[0].Token.Value

	if stdlib.MatchString(file) {
		file = path.Join("stdlib", stdlib.FindStringSubmatch(file)[1]+".hummus")

		b, err := ioutil.ReadFile(file)

		if err != nil {
			panic(err)
		}

		vars := make(map[string]Node, 0)
		_ = Run(parser.Parse(lexer.LexString(string(b))), &vars)

		for k, v := range vars {
			globalFns[k] = v
		}
	}
}

func doVariableCall(node parser.Node, val Node, variables *map[string]Node) Node {
	if val.NodeType == NODETYPE_FN {
		fn := val.Value.(FnLiteral)
		args := getArgs(node.Arguments, fn.Parameters, variables, node.Token.Line)

		if fn.Context != nil {
			for k, v := range fn.Context {
				args[k] = v
			}
		}

		ret := Run(fn.Body, &args)

		if ret.NodeType == NODETYPE_FN {
			arg := resolve(node.Arguments, variables, node.Token.Line)

			// set context in case of currying
			ctx := make(map[string]Node, 0)

			if fn.Context != nil {
				for k, v := range fn.Context {
					ctx[k] = v
				}
			}

			for i := range arg {
				ctx[fn.Parameters[i]] = arg[i]
			}

			ret.Value = FnLiteral{
				Parameters: ret.Value.(FnLiteral).Parameters,
				Body:       ret.Value.(FnLiteral).Body,
				Context:    ctx,
			}
		}

		return ret
	} else if val.NodeType == NODETYPE_STRUCT {
		structDef := val.Value.(StructDef)
		arg := resolve(node.Arguments, variables, node.Token.Line)

		if len(structDef.Parameters) != len(arg) {
			panic(fmt.Sprintf("Struct argument mismatch! (line %d)", node.Token.Line))
		}

		mapNode := MapNode{Values: make(map[string]Node, 0)}

		for i := range structDef.Parameters {
			mapNode.Values[structDef.Parameters[i]] = arg[i]
		}

		return Node{
			Value:    mapNode,
			NodeType: NODETYPE_MAP,
		}
	}

	panic(fmt.Sprintf("Variable %s is not callable! (line %d)", node.Token.Value, node.Token.Line))
}

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

	i := make([]int, 0)

	for _, value := range args[1:] {
		i = append(i, value.Value.(int))
	}

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

func doCall(node parser.Node, variables *map[string]Node) Node {
	if node.Token.Value == USE {
		doUse(node)
	} else if val, ok := globalFns[node.Token.Value]; ok {
		return doVariableCall(node, val, variables)
	} else if val, ok := (*variables)[node.Token.Value]; ok {
		return doVariableCall(node, val, variables)
	} else if node.Token.Type == lexer.ANONYMOUS_FN {
		fn := getValueFromNode(node.Arguments[0], variables)

		node.Arguments = node.Arguments[1:]

		return doVariableCall(node, fn, variables)
	}

	switch node.Token.Value {
	case USE:
		doUse(node)
		return Node{
			Value:    0,
			NodeType: NODETYPE_INT,
		}
	case MAP_ACCESS:
		return accessMap(parser.Node{
			Type:      parser.ACTION_MAP_ACCESS,
			Arguments: node.Arguments,
			Token:     lexer.Token{},
		}, variables)
	case SYSTEM_MATH:
		return doSystemCallMath(node, variables)
	case SYSTEM_MAKE:
		return doSystemCallMake(node, variables)
	case SYSTEM_IO:
		return doSystemCallIo(node, variables)
	case SYSTEM_COMPARE:
		return doSystemCallCompare(node, variables)
	case SYSTEM_COMPARE_ARITHMETIC:
		return doSystemCallCompareArithmetic(node, variables)
	case SYSTEM_CONVERT:
		return doSystemCallConvert(node, variables)
	default:
		panic(fmt.Sprintf("Unknown function %s! (line %d)", node.Token.Value, node.Token.Line))
	}
}

func doIf(node parser.Node, variables *map[string]Node) Node {
	cond := getValueFromNode(node.Arguments[0], variables)

	if cond.NodeType == NODETYPE_BOOL && cond.Value.(bool) {
		return getValueFromNode(node.Arguments[1].Arguments[0], variables)
	}

	if len(node.Arguments) == 3 {
		return getValueFromNode(node.Arguments[2].Arguments[0], variables)
	}

	return Node{
		Value:    0,
		NodeType: 0,
	}
}

func doForLoop(node parser.Node, variables *map[string]Node) {
	context := make(map[string]Node, 0)

	for k, v := range *variables {
		context[k] = v
	}

	list := Run([]parser.Node{
		node.Arguments[1],
	}, variables)

	if list.NodeType != NODETYPE_LIST {
		list.Value = ListNode{Values: []Node{list}}
		list.NodeType = NODETYPE_LIST
	}

	for _, value := range list.Value.(ListNode).Values {
		context[node.Arguments[0].Token.Value] = value
		Run(node.Arguments[2:], &context)
	}
}

func resolve(nodes []parser.Node, variables *map[string]Node, line uint) []Node {
	args := make([]Node, 0)

	for _, node := range nodes {
		args = append(args, getValueFromNode(node, variables))
	}

	return args
}

func getArgs(nodes []parser.Node, parameters []string, variables *map[string]Node, line uint) map[string]Node {
	targetMap := make(map[string]Node, 0)

	for k, v := range *variables {
		targetMap[k] = v
	}

	//TODO: If the args are longer than the parameters, press the rest of the args into the last parameter as list
	if len(parameters) == 1 && len(nodes) > 1 {
		arg := ListNode{Values: resolve(nodes, variables, line)}

		targetMap[parameters[0]] = Node{
			Value:    arg,
			NodeType: NODETYPE_LIST,
		}
	} else if len(parameters) == len(nodes) {
		arg := resolve(nodes, variables, line)
		for i := range nodes {
			targetMap[parameters[i]] = arg[i]
		}
	} else {
		panic(fmt.Sprintf("Argument mismatch! (line %d)", line))
	}

	return targetMap
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

// Run run an AST
func Run(nodes []parser.Node, variables *map[string]Node) (returnVal Node) {
	returnVal = Node{
		Value:    0,
		NodeType: 0,
	}

	for _, node := range nodes {
		switch node.Type {
		case parser.ACTION_DEF:
			returnVal = defineVariable(node, variables)
			break
		//TODO: While loop
		case parser.ACTION_FOR:
			doForLoop(node, variables)
			break
		default:
			returnVal = getValueFromNode(node, variables)
			break
		}
	}

	return
}
