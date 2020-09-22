package interpreter

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"io/ioutil"
	"path"
	"path/filepath"
	"plugin"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// BasePath base path from which to import stdlib
var BasePath string

// LibBasePath path from which to import libraries
var LibBasePath string

type nativeFn func([]Node, *map[string]Node) Node

//noinspection GoSnakeCaseUsage
const (
	// USE include function
	USE string = "use"
	// TYPE type function
	TYPE string = "type"
	// MAP_ACCESS map access function
	MAP_ACCESS string = "[]"
	// FREE object store free function
	FREE string = "free"
	// EXEC_FILE current file
	EXEC_FILE string = "EXEC-FILE"
	// SELF current process id
	SELF string = "self"
)

var literalRe *regexp.Regexp
var libVersionRe *regexp.Regexp

func init() {
	re, err := regexp.Compile("\\|(\\w+)\\|")

	if err != nil {
		panic(err)
	}

	literalRe = re

	re, err = regexp.Compile("^@([\\w-]+)/([\\w\\.-]+)$")

	if err != nil {
		panic(err)
	}

	libVersionRe = re
}

var localFnsMu = &sync.RWMutex{}
var localFns = make(map[string]Node, 0)

var globalFnsMu = &sync.RWMutex{}
var globalFns = make(map[string]Node, 0)

var nativeFnsMu = &sync.RWMutex{}
var nativeFns = make(map[string]nativeFn, 0)

var importsMu = &sync.RWMutex{}
var imports = make([]string, 0)

var localImportsMu = &sync.RWMutex{}
var localImports = make([]string, 0) //we just do filename + importname

var localNativeFnsMu = &sync.RWMutex{}
var localNativeFns = make(map[string]nativeFn, 0)

var objectsMu = &sync.RWMutex{}
var objects = make(map[int]interface{}, 0)

var objectIdCountMu = &sync.Mutex{}
var objectIdCount = 0

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

func localImportHash(importName, currentFile string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(importName+currentFile)))
}

func localImportsHas(currentFile string) func(str string) bool {
	return func(str string) bool {
		has := false
		importHash := localImportHash(str, currentFile)

		localImportsMu.RLock()
		for _, a := range localImports {
			if a == importHash {
				has = true
				break
			}
		}
		localImportsMu.RUnlock()

		return has
	}
}

func loadLocalFns(key, currentFile string) (Node, bool) {
	localFnsMu.RLock()
	val, ok := localFns[localImportHash(key, currentFile)]
	localFnsMu.RUnlock()

	return val, ok
}

func loadLocalNativeFns(key, currentFile string) (nativeFn, bool) {
	localNativeFnsMu.RLock()
	val, ok := localNativeFns[localImportHash(key, currentFile)]
	localNativeFnsMu.RUnlock()

	return val, ok
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

func storeLocalFns(currentFile string) func(key string, val Node) {
	return func(key string, val Node) {
		localFnsMu.Lock()
		localFns[localImportHash(key, currentFile)] = val
		localFnsMu.Unlock()
	}
}

func storeLocalNativeFns(currentFile string) func(key string, val nativeFn) {
	return func(key string, val nativeFn) {
		localNativeFnsMu.Lock()
		localNativeFns[localImportHash(key, currentFile)] = val
		localNativeFnsMu.Unlock()
	}
}

func loadNativeFns(key string) (nativeFn, bool) {
	nativeFnsMu.RLock()
	val, ok := nativeFns[key]
	nativeFnsMu.RUnlock()

	return val, ok
}

func storeNativeFns(key string, val nativeFn) {
	nativeFnsMu.Lock()
	nativeFns[key] = val
	nativeFnsMu.Unlock()
}

// LoadObject load an object from the object store by its ptr
func LoadObject(key int) (interface{}, bool) {
	objectsMu.RLock()
	val, ok := objects[key]
	objectsMu.RUnlock()

	return val, ok
}

// StoreObject store an object and return a pseudo ptr to it
func StoreObject(val interface{}) int {
	objectIdCountMu.Lock()
	id := objectIdCount
	objectIdCount++
	objectIdCountMu.Unlock()

	objectsMu.Lock()
	objects[id] = val
	objectsMu.Unlock()

	return id
}

func getLiteralFn(parserNode parser.Node, variables *map[string]Node) (nodeType NodeType, nodeValue interface{}) {
	nodeType = NODETYPE_FN

	parameters := make([]string, 0)
	for _, argument := range parserNode.Arguments[0].Arguments {
		parameters = append(parameters, argument.Token.Value)
	}

	nodeValue = FnLiteral{
		Parameters: parameters,
		Body:       parserNode.Arguments[1:],
		Context:    map[string]Node{EXEC_FILE: (*variables)[EXEC_FILE]},
	}

	return
}

func getIdentifier(parserNode parser.Node, variables *map[string]Node) Node {
	if val, ok := (*variables)[parserNode.Token.Value]; ok {
		return val
	} else if val, ok := loadLocalFns(parserNode.Token.Value, (*variables)[EXEC_FILE].Value.(string)); ok {
		return val
	} else if val, ok := loadGlobalFns(parserNode.Token.Value); ok {
		return val
	}

	panic(fmt.Sprintf("Unknown variable %s! (line %d)", parserNode.Token.Value, parserNode.Token.Line))
}

func getValueFromNode(parserNode parser.Node, variables *map[string]Node) (node Node) {
	node = Nothing

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
		node.NodeType, node.Value = getLiteralFn(parserNode, variables)
		break
	case parser.IDENTIFIER:
		node = getIdentifier(parserNode, variables)
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
	mapVals := make(map[string]Node, 0)

	for _, argument := range node.Arguments {
		mapVals[argument.Arguments[0].Token.Value] = getValueFromNode(argument.Arguments[1], variables)
	}

	return NodeMap(mapVals)
}

func getStructDef(node parser.Node) StructDef {
	structDef := StructDef{Parameters: make([]string, 0)}

	for _, argument := range node.Arguments[1].Arguments {
		structDef.Parameters = append(structDef.Parameters, argument.Token.Value)
	}

	return structDef
}

func getMacroDef(parserNode parser.Node) MacroDef {
	parameters := make([]MacroParameter, 0)
	for _, argument := range parserNode.Arguments[1].Arguments[0].Arguments {
		if literalRe.MatchString(argument.Token.Value) {
			parameters = append(parameters, MacroParameter{
				Parameter: literalRe.FindStringSubmatch(argument.Token.Value)[1],
				Literal:   true,
			})
		} else {
			parameters = append(parameters, MacroParameter{
				Parameter: argument.Token.Value,
				Literal:   false,
			})
		}
	}

	return MacroDef{
		Parameters: parameters,
		Body:       parserNode.Arguments[1].Arguments[1:],
	}
}

func defineVariable(node parser.Node, variables *map[string]Node) Node {
	name := node.Arguments[0].Token.Value
	variable := Nothing

	//I am doing this here, just like I do macros here because you can only define macros and structs in `def`
	if node.Arguments[1].Type == parser.STRUCT_DEF {
		(*variables)[name] = Node{
			Value:    getStructDef(node),
			NodeType: NODETYPE_STRUCT,
		}
	} else if node.Arguments[1].Type == parser.MACRO_DEF {
		(*variables)[name] = Node{
			Value:    getMacroDef(node),
			NodeType: NODETYPE_MACRO,
		}
	} else {
		(*variables)[name] = getValueFromNode(node.Arguments[1], variables)
	}

	return variable
}

func doNativeUse(name, currentFile string, variables *map[string]Node, hasImport func(importName string) bool, storeNativeFn func(key string, val nativeFn), saveImportedFile func(file string)) Node {
	dir, err := filepath.Abs(filepath.Dir(currentFile))

	if err != nil {
		panic(err)
	}

	file := path.Join(dir, name)

	if hasImport(file) {
		return Nothing
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

	storeNativeFn(dName, dFn)

	saveImportedFile(file)

	return Nothing
}

func doFileUse(node parser.Node, currentFile string, pid int, hasImport func(importName string) bool, storeState func(key string, val Node), saveImportedFile func(file string)) Node {
	stdlib, _ := regexp.Compile("<([^>]+)>")

	file := node.Arguments[0].Token.Value

	if stdlib.MatchString(file) {
		dir, err := filepath.Abs(BasePath)

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

	if hasImport(file) {
		return Nothing
	}

	b, err := ioutil.ReadFile(file)

	if err != nil {
		panic(err)
	}

	vars := make(map[string]Node, 0)
	vars[EXEC_FILE] = StringNode(file)
	vars[SELF] = IntNode(pid)

	_ = Run(parser.Parse(lexer.LexString(string(b))), &vars)

	for k, v := range vars {
		storeState(k, v)
	}

	saveImportedFile(file)

	return Nothing
}

func getLibPath(libImport string) (projectJson string, libPath string) {
	if libVersionRe.MatchString(libImport) {
		groups := libVersionRe.FindStringSubmatch(libImport)

		if groups[2] != "master" {
			libPath = groups[1] + "@" + strings.ReplaceAll(groups[2], ".", "_")
			projectJson = path.Join(LibBasePath, libPath, "project.json")
		} else {
			libPath = groups[1]
			projectJson = path.Join(LibBasePath, libPath, "project.json")
		}
		return
	}

	libPath = libImport[1:]
	projectJson = path.Join(LibBasePath, libPath, "project.json")
	return
}

func getLibEntryPath(libraryName, currentFile string) string {
	projectJson, libPath := getLibPath(libraryName)

	b, err := ioutil.ReadFile(projectJson)

	if err != nil {
		panic(err)
	}

	settings := make(map[string]interface{})

	err = json.Unmarshal(b, &settings)

	if err != nil {
		panic(err)
	}

	p, err := filepath.Rel(
		filepath.Dir(currentFile),
		path.Join(LibBasePath, libPath, settings["output"].(string), ReplaceEnd(settings["entry"].(string), ".hummus", "", 1)))

	if err != nil {
		panic(err)
	}

	return p
}

func doUse(node parser.Node, currentFile string, pid int, variables *map[string]Node) Node {
	if node.Arguments[0].Type != parser.LITERAL_ATOM {
		panic(fmt.Sprintf("\"use\" only accepts atoms as parameter! (line %d)", node.Token.Line))
	}

	if len(node.Arguments) == 2 && node.Arguments[1].Token.Value == "native" {
		return doNativeUse(node.Arguments[0].Token.Value, currentFile, variables, importsHas, storeNativeFns, func(file string) {
			importsMu.Lock()
			imports = append(imports, file)
			importsMu.Unlock()
		})
	}

	if len(node.Arguments) == 3 &&
		((node.Arguments[1].Token.Value == "native" && node.Arguments[2].Token.Value == "local") ||
			(node.Arguments[1].Token.Value == "local" && node.Arguments[2].Token.Value == "native")) {
		return doNativeUse(node.Arguments[0].Token.Value, currentFile, variables, localImportsHas(currentFile), storeLocalNativeFns(currentFile), func(file string) {
			localImportsMu.Lock()
			localImports = append(localImports, localImportHash(file, currentFile))
			localImportsMu.Unlock()
		})
	}

	if node.Arguments[0].Token.Value[0] == '@' && LibBasePath != "" {
		node.Arguments[0].Token.Value = getLibEntryPath(node.Arguments[0].Token.Value, currentFile)
	}

	if len(node.Arguments) == 2 && node.Arguments[1].Token.Value == "local" {
		return doFileUse(node, currentFile, pid, localImportsHas(currentFile), storeLocalFns(currentFile), func(file string) {
			localImportsMu.Lock()
			localImports = append(localImports, localImportHash(file, currentFile))
			localImportsMu.Unlock()
		})
	}

	return doFileUse(node, currentFile, pid, importsHas, storeGlobalFns, func(file string) {
		importsMu.Lock()
		imports = append(imports, file)
		importsMu.Unlock()
	})
}

func doFnCall(node parser.Node, val Node, variables *map[string]Node) Node {
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
}

func doStructCreation(node parser.Node, val Node, variables *map[string]Node) Node {
	structDef := val.Value.(StructDef)
	arg := resolve(node.Arguments, variables)

	if len(structDef.Parameters) != len(arg) {
		panic(fmt.Sprintf("Struct argument mismatch! (line %d)", node.Token.Line))
	}

	mapVals := make(map[string]Node, 0)

	for i := range structDef.Parameters {
		mapVals[structDef.Parameters[i]] = arg[i]
	}

	return NodeMap(mapVals)
}

func doMacroCall(node parser.Node, val Node, variables *map[string]Node) Node {
	macro := val.Value.(MacroDef)

	targetMap := make(map[string]Node, 0)
	CopyVariableState(variables, &targetMap)

	if len(macro.Parameters) == 1 && len(node.Arguments) > 1 {
		//1 macro argument, many inputs => list

		if macro.Parameters[0].Literal {
			l := make([]Node, 0)
			for _, argument := range node.Arguments {
				l = append(l, parserNodeToAstList(argument))
			}

			targetMap[macro.Parameters[0].Parameter] = NodeList(l)
		} else {
			targetMap[macro.Parameters[0].Parameter] = NodeList(resolve(node.Arguments, variables))
		}
	} else if len(macro.Parameters) == len(node.Arguments) {
		//arguments 1:1

		for i := range macro.Parameters {
			if macro.Parameters[i].Literal {
				targetMap[macro.Parameters[i].Parameter] = parserNodeToAstList(node.Arguments[i])
			} else {
				targetMap[macro.Parameters[i].Parameter] = getValueFromNode(node.Arguments[i], variables)
			}
		}
	} else if len(macro.Parameters) < len(node.Arguments) && len(macro.Parameters) > 0 {
		//arguments 1:1 until end, last variable is a list

		i := 0
		for i = range macro.Parameters[:len(macro.Parameters)-1] {
			if macro.Parameters[i].Literal {
				targetMap[macro.Parameters[i].Parameter] = parserNodeToAstList(node.Arguments[i])
			} else {
				targetMap[macro.Parameters[i].Parameter] = getValueFromNode(node.Arguments[i], variables)
			}
		}

		if macro.Parameters[i+1].Literal {
			l := make([]Node, 0)
			for _, argument := range node.Arguments[i+1:] {
				l = append(l, parserNodeToAstList(argument))
			}

			targetMap[macro.Parameters[i+1].Parameter] = NodeList(l)
		} else {
			targetMap[macro.Parameters[i+1].Parameter] = NodeList(resolve(node.Arguments[i+1:], variables))
		}
	} else {
		panic(fmt.Sprintf("Argument mismatch! (line %d)", node.Token.Line))
	}

	ret := Run(macro.Body, &targetMap)

	if ret.NodeType != NODETYPE_LIST {
		panic(fmt.Sprintf("Macro %s doesn't return a list! (line %d)", node.Token.Value, node.Token.Line))
	}

	//Convert ret to Ast and run it
	ast := astListToParserNode(ret.Value.(ListNode))

	return Run(ast, variables)
}

// DoVariableCall calls a fn variable
func DoVariableCall(node parser.Node, val Node, variables *map[string]Node) Node {
	switch val.NodeType {
	case NODETYPE_FN:
		return doFnCall(node, val, variables)
	case NODETYPE_STRUCT:
		return doStructCreation(node, val, variables)
	case NODETYPE_MACRO:
		return doMacroCall(node, val, variables)
	default:
		panic(fmt.Sprintf("Variable %s is not callable! (line %d)", node.Token.Value, node.Token.Line))
	}
}

func doType(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables)

	if len(args) != 1 {
		panic(fmt.Sprintf("Expected one argument for type! (line %d)", node.Token.Line))
	}

	return AtomNode(args[0].NodeType.String())
}

func doCall(node parser.Node, variables *map[string]Node) Node {
	currentFile := (*variables)[EXEC_FILE].Value.(string)
	if val, ok := (*variables)[node.Token.Value]; ok {
		return DoVariableCall(node, val, variables)
	} else if val, ok := loadLocalFns(node.Token.Value, currentFile); ok {
		return DoVariableCall(node, val, variables)
	} else if val, ok := loadLocalNativeFns(node.Token.Value, currentFile); ok {
		args := resolve(node.Arguments, variables)
		return val(args, variables)
	} else if val, ok := loadGlobalFns(node.Token.Value); ok {
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
	case FREE:
		return doFree(node, variables)
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

func doFree(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables)

	if args[0].NodeType != NODETYPE_INT {
		panic(fmt.Sprintf("First argument in free must be an int! (line %d)", node.Arguments[0].Token.Line))
	}

	objectsMu.Lock()
	defer objectsMu.Unlock()

	delete(objects, args[0].Value.(int))

	return Nothing
}

func doIf(node parser.Node, variables *map[string]Node) Node {
	cond := getValueFromNode(node.Arguments[0], variables)

	if cond.NodeType == NODETYPE_BOOL && cond.Value.(bool) {
		return getValueFromNode(node.Arguments[1].Arguments[0], variables)
	}

	if len(node.Arguments) == 3 {
		return getValueFromNode(node.Arguments[2].Arguments[0], variables)
	}

	return Nothing
}

func doForLoop(node parser.Node, variables *map[string]Node) {
	context := make(map[string]Node, 0)
	CopyVariableState(variables, &context)

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
	CopyVariableState(variables, &context)

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
		(*targetMap)[parameters[0]] = NodeList(resolve(nodes, variables))
	} else if len(parameters) == len(nodes) {
		arg := resolve(nodes, variables)
		for i := range nodes {
			(*targetMap)[parameters[i]] = arg[i]
		}
	} else if len(parameters) < len(nodes) && len(parameters) > 0 {
		arg := resolve(nodes, variables)
		i := 0
		for i = range parameters[:len(parameters)-1] {
			(*targetMap)[parameters[i]] = arg[i]
		}

		(*targetMap)[parameters[i+1]] = NodeList(arg[i+1:])
	} else {
		panic(fmt.Sprintf("Argument mismatch! (line %d)", line))
	}
}

func getArgs(nodes []parser.Node, parameters []string, variables *map[string]Node, line uint) map[string]Node {
	targetMap := make(map[string]Node, 0)
	CopyVariableState(variables, &targetMap)

	getArgsByParameterList(nodes, variables, parameters, &targetMap, line)

	return targetMap
}

// Run run an AST
func Run(nodes []parser.Node, variables *map[string]Node) (returnVal Node) {
	returnVal = Nothing

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
