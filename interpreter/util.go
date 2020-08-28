package interpreter

import (
	"fmt"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
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

func parserNodeToAstList(node parser.Node) Node {
	switch node.Type {
	case parser.ACTION_CALL:
		astList := []Node{AtomNode("call"), AtomNode(node.Token.Value)}

		for _, argument := range node.Arguments {
			astList = append(astList, parserNodeToAstList(argument))
		}

		return NodeList(astList)
	case parser.IDENTIFIER:
		return NodeList([]Node{AtomNode("identifier"), AtomNode(node.Token.Value)})
	case parser.LITERAL_FN:
		paramList := []Node{AtomNode("parameters")}
		for _, param := range node.Arguments[0].Arguments {
			paramList = append(paramList, AtomNode(param.Token.Value))
		}

		astList := []Node{AtomNode("fn"), NodeList(paramList)}

		for _, argument := range node.Arguments[1:] {
			astList = append(astList, parserNodeToAstList(argument))
		}

		return NodeList(astList)
	case parser.LITERAL_ATOM:
		return NodeList([]Node{AtomNode("atom"), AtomNode(node.Token.Value)})
	case parser.LITERAL_BOOL:
		b, err := strconv.ParseBool(node.Token.Value)

		if err != nil {
			panic(err)
		}

		return NodeList([]Node{AtomNode("bool"), BoolNode(b)})
	case parser.LITERAL_FLOAT:
		f, err := strconv.ParseFloat(node.Token.Value, 64)

		if err != nil {
			panic(err)
		}

		return NodeList([]Node{AtomNode("float"), FloatNode(f)})
	case parser.LITERAL_INT:
		i, err := strconv.Atoi(node.Token.Value)

		if err != nil {
			panic(err)
		}

		return NodeList([]Node{AtomNode("int"), IntNode(i)})
	case parser.LITERAL_STRING:
		return NodeList([]Node{AtomNode("string"), StringNode(node.Token.Value)})
	}
	return Node{}
}

func astListToParserNode(list ListNode) []parser.Node {
	if len(list.Values) == 0 {
		panic("Empty list cannot be converted to AST node!")
	}

	if list.Values[0].NodeType == NODETYPE_LIST {
		astNodes := make([]parser.Node, 0)
		for _, value := range list.Values {
			astNodes = append(astNodes, astListToParserNode(value.Value.(ListNode))...)
		}

		return astNodes
	}

	//TODO: for, def
	if list.Values[0].NodeType != NODETYPE_ATOM {
		panic("Atom required in macro for AST node type!")
	}

	if len(list.Values) < 2 {
		panic("At least two values equired in macro for AST node!")
	}

	switch list.Values[0].Value.(string) {
	case "def":
		//TODO: Ensure types
		return []parser.Node{
			{
				Type: parser.ACTION_DEF,
				Arguments: []parser.Node{
					{
						Type:      parser.IDENTIFIER,
						Arguments: nil,
						Token: lexer.Token{
							Value: list.Values[1].Value.(string),
							Type:  lexer.IDENTIFIER,
							Line:  0,
						},
					},
					astListToParserNode(list.Values[2].Value.(ListNode))[0],
				},
				Token: lexer.Token{},
			},
		}
	case "call":
		args := make([]parser.Node, 0)

		if len(list.Values) > 2 {
			for i, node := range list.Values[2:] {
				EnsureSingleType(&node, i+2, NODETYPE_LIST, "Macro call definition")

				args = append(args, astListToParserNode(node.Value.(ListNode))[0])
			}
		}

		return []parser.Node{
			{
				Type:      parser.ACTION_CALL,
				Arguments: args,
				Token: lexer.Token{
					Value: list.Values[1].Value.(string),
					Type:  lexer.IDENTIFIER,
					Line:  0,
				},
			},
		}
	case "if":
		EnsureSingleType(&list.Values[1], 2, NODETYPE_LIST, "Macro if definition")
		EnsureSingleType(&list.Values[2], 3, NODETYPE_LIST, "Macro if definition")

		if len(list.Values) == 3 {
			return []parser.Node{
				{
					Type: parser.ACTION_IF,
					Arguments: []parser.Node{
						astListToParserNode(list.Values[1].Value.(ListNode))[0],
						{
							Type:      parser.ACTION_BRANCH,
							Arguments: astListToParserNode(list.Values[2].Value.(ListNode)),
						},
					},
					Token: lexer.Token{},
				},
			}
		}

		EnsureSingleType(&list.Values[3], 4, NODETYPE_LIST, "Macro if definition")

		return []parser.Node{
			{
				Type: parser.ACTION_IF,
				Arguments: []parser.Node{
					astListToParserNode(list.Values[1].Value.(ListNode))[0],
					{
						Type:      parser.ACTION_BRANCH,
						Arguments: astListToParserNode(list.Values[2].Value.(ListNode)),
					},
					{
						Type:      parser.ACTION_BRANCH,
						Arguments: astListToParserNode(list.Values[3].Value.(ListNode)),
					},
				},
				Token: lexer.Token{},
			},
		}
	case "fn":
		//TODO: Ensure types
		paramList := make([]parser.Node, 0)

		for _, node := range list.Values[1].Value.(ListNode).Values[1:] {
			paramList = append(paramList, parser.Node{
				Type:      parser.IDENTIFIER,
				Arguments: nil,
				Token: lexer.Token{
					Value: node.Value.(string),
					Type:  lexer.IDENTIFIER,
					Line:  0,
				},
			})
		}

		actionList := make([]parser.Node, 0)

		for _, node := range list.Values[2:] {
			actionList = append(actionList, astListToParserNode(node.Value.(ListNode))...)
		}

		fnDef := parser.Node{
			Type: parser.LITERAL_FN,
			Arguments: []parser.Node{
				{
					Type:      parser.PARAMETERS,
					Arguments: paramList,
					Token:     lexer.Token{},
				},
			},
			Token: lexer.Token{},
		}

		fnDef.Arguments = append(fnDef.Arguments, actionList...)

		return []parser.Node{
			fnDef,
		}
	case "for_iter":
		//TODO: Ensure types
		actionList := make([]parser.Node, 0)

		for _, node := range list.Values[2:] {
			actionList = append(actionList, astListToParserNode(node.Value.(ListNode))...)
		}

		forIter := []parser.Node{
			{
				Type:      parser.IDENTIFIER,
				Arguments: nil,
				Token: lexer.Token{
					Value: list.Values[1].Value.(ListNode).Values[0].Value.(string),
					Type:  lexer.IDENTIFIER,
					Line:  0,
				},
			},
			astListToParserNode(list.Values[1].Value.(ListNode).Values[1].Value.(ListNode))[0],
		}

		return []parser.Node{
			{
				Type:      parser.ACTION_FOR,
				Arguments: append(forIter, actionList...),
				Token:     lexer.Token{},
			},
		}
	case "for_loop":
		//TODO: Ensure types
		actionList := make([]parser.Node, 0)

		for _, node := range list.Values[2:] {
			actionList = append(actionList, astListToParserNode(node.Value.(ListNode))...)
		}

		return []parser.Node{
			{
				Type: parser.ACTION_WHILE,
				Arguments: append([]parser.Node{
					astListToParserNode(list.Values[1].Value.(ListNode))[0],
				}, actionList...),
				Token: lexer.Token{},
			},
		}
	case "identifier":
		if list.Values[1].NodeType != NODETYPE_ATOM {
			panic("AST node in macro declared as identifier (atom), but isn't!")
		}

		return []parser.Node{
			{
				Type:      parser.IDENTIFIER,
				Arguments: nil,
				Token: lexer.Token{
					Value: list.Values[1].Value.(string),
					Type:  lexer.IDENTIFIER,
					Line:  0,
				},
			},
		}
	case "atom":
		if list.Values[1].NodeType != NODETYPE_ATOM {
			panic("AST node in macro declared as atom, but isn't!")
		}

		return []parser.Node{
			{
				Type:      parser.LITERAL_ATOM,
				Arguments: nil,
				Token: lexer.Token{
					Value: list.Values[1].Value.(string),
					Type:  lexer.ATOM,
					Line:  0,
				},
			},
		}
	case "bool":
		if list.Values[1].NodeType != NODETYPE_BOOL {
			panic("AST node in macro declared as bool, but isn't!")
		}

		return []parser.Node{
			{
				Type:      parser.LITERAL_BOOL,
				Arguments: nil,
				Token: lexer.Token{
					Value: strconv.FormatBool(list.Values[1].Value.(bool)),
					Type:  lexer.BOOL,
					Line:  0,
				},
			},
		}
	case "int":
		if list.Values[1].NodeType != NODETYPE_INT {
			panic("AST node in macro declared as int, but isn't!")
		}

		return []parser.Node{
			{
				Type:      parser.LITERAL_INT,
				Arguments: nil,
				Token: lexer.Token{
					Value: strconv.Itoa(list.Values[1].Value.(int)),
					Type:  lexer.INT,
					Line:  0,
				},
			},
		}
	case "float":
		if list.Values[1].NodeType != NODETYPE_FLOAT {
			panic("AST node in macro declared as float, but isn't!")
		}

		return []parser.Node{
			{
				Type:      parser.LITERAL_FLOAT,
				Arguments: nil,
				Token: lexer.Token{
					Value: strconv.FormatFloat(list.Values[1].Value.(float64), 'f', 6, 64),
					Type:  lexer.FLOAT,
					Line:  0,
				},
			},
		}
	case "string":
		if list.Values[1].NodeType != NODETYPE_STRING {
			panic("AST node in macro declared as string, but isn't!")
		}

		return []parser.Node{
			{
				Type:      parser.LITERAL_STRING,
				Arguments: nil,
				Token: lexer.Token{
					Value: list.Values[1].Value.(string),
					Type:  lexer.STRING,
					Line:  0,
				},
			},
		}
	default:
		panic("Unknown AST type " + list.Values[0].Value.(string))
	}
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
