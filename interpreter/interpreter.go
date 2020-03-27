package interpreter

import (
	"fmt"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"regexp"
	"strconv"
	"sync"
)

const (
	// USE include function
	USE string = "use"
	// TYPE type function
	TYPE string = "type"
	// MAP_ACCESS map access function
	MAP_ACCESS string = "[]"
	// EXEC_FILE current file
	EXEC_FILE string = "EXEC-FILE"
	// SELF current process id
	SELF string = "self"
)

var globalFnsMu = &sync.RWMutex{}
var globalFns = make(map[string]Node, 0)

var nativeFnsMu = &sync.RWMutex{}
var nativeFns = make(map[string]func([]Node, *map[string]Node) Node, 0)

var importsMu = &sync.RWMutex{}
var imports = make([]string, 0)

func importsHas(str string) bool {
	has := false

	importsMu.RLock()
	for _, a := range imports {
		if a == str {
			has = true
			break
		}
	}
	importsMu.RUnlock()

	return has
}

func loadGlobalFns(key string) (Node, bool) {
	globalFnsMu.RLock()
	val, ok := globalFns[key]
	globalFnsMu.RUnlock()

	return val, ok
}

func storeGlobalFns(key string, val Node) {
	globalFnsMu.Lock()
	globalFns[key] = val
	globalFnsMu.Unlock()
}

func loadNativeFns(key string) (func([]Node, *map[string]Node) Node, bool) {
	nativeFnsMu.RLock()
	val, ok := nativeFns[key]
	nativeFnsMu.RUnlock()

	return val, ok
}

func storeNativeFns(key string, val func([]Node, *map[string]Node) Node) {
	nativeFnsMu.Lock()
	nativeFns[key] = val
	nativeFnsMu.Unlock()
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
		} else if val, ok := loadGlobalFns(parserNode.Token.Value); ok {
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
	args := resolve(node.Arguments, variables)

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

func doNativeUse(name, currentFile string, variables *map[string]Node) Node {
	dir, err := filepath.Abs(filepath.Dir(currentFile))

	if err != nil {
		panic(err)
	}

	file := path.Join(dir, name)

	if importsHas(file) {
		return Node{
			Value:    0,
			NodeType: 0,
		}
	}

	p, err := plugin.Open(file)
	if err != nil {
		panic(err)
	}

	n, err := p.Lookup("CALL")
	if err != nil {
		panic(err)
	}

	i, err := p.Lookup("Init")
	if err != nil {
		panic(err)
	}

	i.(func(*map[string]Node))(variables)

	fn, err := p.Lookup("DoSystemCall")
	if err != nil {
		panic(err)
	}

	dName := *n.(*string)
	dFn := fn.(func([]Node, *map[string]Node) Node)

	storeNativeFns(dName, dFn)

	importsMu.Lock()
	imports = append(imports, file)
	importsMu.Unlock()

	return Node{
		Value:    0,
		NodeType: 0,
	}
}

func doUse(node parser.Node, currentFile string, pid int, variables *map[string]Node) Node {
	if node.Arguments[0].Type != parser.LITERAL_ATOM {
		panic(fmt.Sprintf("\"use\" only accepts atoms as parameter! (line %d)", node.Token.Line))
	}

	if len(node.Arguments) == 2 && node.Arguments[1].Token.Value == "native" {
		return doNativeUse(node.Arguments[0].Token.Value, currentFile, variables)
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

	if importsHas(file) {
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
		storeGlobalFns(k, v)
	}

	importsMu.Lock()
	imports = append(imports, file)
	importsMu.Unlock()

	return Node{
		Value:    0,
		NodeType: 0,
	}
}

// DoVariableCall calls a fn variable
func DoVariableCall(node parser.Node, val Node, variables *map[string]Node) Node {
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
			// set context in case of currying
			ctx := make(map[string]Node, 0)

			if fn.Context != nil {
				for k, v := range fn.Context {
					ctx[k] = v
				}
			}

			retFn := ret.Value.(FnLiteral)
			if retFn.Context != nil {
				for k, v := range retFn.Context {
					ctx[k] = v
				}
			}

			getArgsByParameterList(node.Arguments, variables, fn.Parameters, &ctx, node.Token.Line)

			ret.Value = FnLiteral{
				Parameters: ret.Value.(FnLiteral).Parameters,
				Body:       ret.Value.(FnLiteral).Body,
				Context:    ctx,
			}
		}

		return ret
	} else if val.NodeType == NODETYPE_STRUCT {
		structDef := val.Value.(StructDef)
		arg := resolve(node.Arguments, variables)

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

func doType(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables)

	if len(args) != 1 {
		panic(fmt.Sprintf("Expected one argument for type! (line %d)", node.Token.Line))
	}

	t := ""

	switch args[0].NodeType {
	case NODETYPE_INT:
		t = "int"
	case NODETYPE_FLOAT:
		t = "float"
	case NODETYPE_STRING:
		t = "string"
	case NODETYPE_BOOL:
		t = "bool"
	case NODETYPE_ATOM:
		t = "atom"
	case NODETYPE_FN:
		t = "fn"
	case NODETYPE_LIST:
		t = "list"
	case NODETYPE_MAP:
		t = "map"
	case NODETYPE_STRUCT:
		t = "struct"
	}

	return Node{
		Value:    t,
		NodeType: NODETYPE_ATOM,
	}
}

func doCall(node parser.Node, variables *map[string]Node) Node {
	if val, ok := loadGlobalFns(node.Token.Value); ok {
		return DoVariableCall(node, val, variables)
	} else if val, ok := (*variables)[node.Token.Value]; ok {
		return DoVariableCall(node, val, variables)
	} else if val, ok := loadNativeFns(node.Token.Value); ok {
		args := resolve(node.Arguments, variables)
		return val(args, variables)
	} else if node.Token.Type == lexer.ANONYMOUS_FN {
		fn := getValueFromNode(node.Arguments[0], variables)

		node.Arguments = node.Arguments[1:]

		return DoVariableCall(node, fn, variables)
	}

	switch node.Token.Value {
	case USE:
		return doUse(node, (*variables)[EXEC_FILE].Value.(string), (*variables)[SELF].Value.(int), variables)
	case TYPE:
		return doType(node, variables)
	case MAP_ACCESS:
		return accessMap(parser.Node{
			Type:      parser.ACTION_MAP_ACCESS,
			Arguments: node.Arguments,
			Token:     lexer.Token{},
		}, variables)
	case BUILTIN_MATH:
		return builtInMath(node, variables)
	case BUILTIN_COMPARE:
		return builtInCompare(node, variables)
	case BUILTIN_COMPARE_ARITHMETIC:
		return builtInCompareArithmetic(node, variables)
	case BUILTIN_BOOL:
		return builtInBool(node, variables)
	case BUILTIN_BITWISE:
		return builtInBitwise(node, variables)

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

func resolve(nodes []parser.Node, variables *map[string]Node) []Node {
	args := make([]Node, 0)

	for _, node := range nodes {
		args = append(args, getValueFromNode(node, variables))
	}

	return args
}

func getArgsByParameterList(nodes []parser.Node, variables *map[string]Node, parameters []string, targetMap *map[string]Node, line uint) {
	if len(parameters) == 1 && len(nodes) > 1 {
		arg := ListNode{Values: resolve(nodes, variables)}

		(*targetMap)[parameters[0]] = Node{
			Value:    arg,
			NodeType: NODETYPE_LIST,
		}
	} else if len(parameters) == len(nodes) {
		arg := resolve(nodes, variables)
		for i := range nodes {
			(*targetMap)[parameters[i]] = arg[i]
		}
	} else if len(parameters) < len(nodes) {
		arg := resolve(nodes, variables)
		i := 0
		for i = range parameters[:len(parameters)-1] {
			(*targetMap)[parameters[i]] = arg[i]
		}

		(*targetMap)[parameters[i+1]] = Node{
			Value:    ListNode{Values: arg[i+1:]},
			NodeType: NODETYPE_LIST,
		}
	} else {
		panic(fmt.Sprintf("Argument mismatch! (line %d)", line))
	}
}

func getArgs(nodes []parser.Node, parameters []string, variables *map[string]Node, line uint) map[string]Node {
	targetMap := make(map[string]Node, 0)

	for k, v := range *variables {
		targetMap[k] = v
	}

	getArgsByParameterList(nodes, variables, parameters, &targetMap, line)

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
