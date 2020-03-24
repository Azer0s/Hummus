package interpreter

import (
	"fmt"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
)

const (
	// USE include function
	USE string = "use"
	// MAP_ACCESS map access function
	MAP_ACCESS string = "[]"
	// EXEC_FILE current file
	EXEC_FILE string = "EXEC-FILE"
	// SELF current process id
	SELF string = "self"
)

var globalFns = sync.Map{}
var imports = make([]string, 0)

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

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
		if val, ok := (*variables)[parserNode.Token.Value]; ok {
			node = val
		} else {
			panic(fmt.Sprintf("Unknown variable %s! (line %d)", parserNode.Token.Value, parserNode.Token.Line))
		}
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

func doUse(node parser.Node, currentFile string, pid int) Node {
	if node.Arguments[0].Type != parser.LITERAL_ATOM {
		panic(fmt.Sprintf("\"use\" only accepts atoms as parameter! (line %d)", node.Token.Line))
	}

	stdlib, _ := regexp.Compile("<([^>]+)>")

	file := node.Arguments[0].Token.Value

	if stdlib.MatchString(file) {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			panic(err)
		}

		file = path.Join(dir, "./stdlib", stdlib.FindStringSubmatch(file)[1]+".hummus")
	} else {
		//Let's try a relative path
		dir, err := filepath.Abs(filepath.Dir(currentFile))

		if err != nil {
			panic(err)
		}

		file = path.Join(dir, file+".hummus")
	}

	if contains(imports, file) {
		return Node{
			Value:    0,
			NodeType: 0,
		}
	}

	b, err := ioutil.ReadFile(file)

	if err != nil {
		panic(err)
	}

	vars := make(map[string]Node, 0)
	vars[EXEC_FILE] = Node{
		Value:    file,
		NodeType: NODETYPE_STRING,
	}
	vars[SELF] = Node{
		Value:    pid,
		NodeType: NODETYPE_INT,
	}
	_ = Run(parser.Parse(lexer.LexString(string(b))), &vars)

	for k, v := range vars {
		globalFns.Store(k, v)
	}

	imports = append(imports, file)

	return Node{
		Value:    0,
		NodeType: 0,
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

func doCall(node parser.Node, variables *map[string]Node) Node {
	if node.Token.Value == USE {
		return doUse(node, (*variables)[EXEC_FILE].Value.(string), (*variables)[SELF].Value.(int))
	} else if val, ok := globalFns.Load(node.Token.Value); ok {
		return doVariableCall(node, val.(Node), variables)
	} else if val, ok := (*variables)[node.Token.Value]; ok {
		return doVariableCall(node, val, variables)
	} else if node.Token.Type == lexer.ANONYMOUS_FN {
		fn := getValueFromNode(node.Arguments[0], variables)

		node.Arguments = node.Arguments[1:]

		return doVariableCall(node, fn, variables)
	}

	switch node.Token.Value {
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
	case SYSTEM_BOOL:
		return doSystemCallBool(node, variables)
	case SYSTEM_BITWISE:
		return doSystemCallBitwise(node, variables)
	case SYSTEM_ENUMERATE:
		return doSystemCallEnumerate(node, variables)
	case SYSTEM_STRING:
		return doSystemCallStrings(node, variables)
	case SYSTEM_DEBUG:
		return doSystemCallDebug(node, variables)
	case SYSTEM_SYNC:
		return doSystemCallSync(node, variables)

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

func doWhileLoop(node parser.Node, variables *map[string]Node) {
	context := make(map[string]Node, 0)

	for k, v := range *variables {
		context[k] = v
	}

	for {
		val := Run([]parser.Node{
			node.Arguments[0],
		}, variables)

		if val.NodeType != NODETYPE_BOOL {
			panic(fmt.Sprintf("Expected a bool for loop! (line %d)", node.Arguments[0].Token.Line))
		}

		if !val.Value.(bool) {
			break
		}

		Run(node.Arguments[1:], &context)
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
		case parser.ACTION_WHILE:
			doWhileLoop(node, variables)
			break
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
