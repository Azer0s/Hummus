package interpreter

import (
	"fmt"
	"github.com/Azer0s/Hummus/parser"
	"strconv"
)

const (
	SYSTEM_MATH string = "--system-do-math!"
)

func getLiteralFromNode(parserNode parser.Node) (node Node) {
	node = Node{
		Value:    nil,
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
	}

	return
}

func defineVariable(node parser.Node, variables *map[string]Node) Node {
	name := node.Arguments[0].Token.Value
	variable := Node{
		Value:    nil,
		NodeType: 0,
	}

	(*variables)[name] = getLiteralFromNode(node.Arguments[1])
	return variable
}

func doVariableCall(node parser.Node, val Node, variables *map[string]Node) Node {
	if val.NodeType == NODETYPE_FN {
		fn := val.Value.(FnLiteral)
		args := getArgs(node.Arguments, fn.Parameters, variables, node.Token.Line)
		//TODO: Set context for function returns (so...if a function was returned, the state of the args shall be saved
		return Run(fn.Body, args)
	} else {
		panic(fmt.Sprintf("Variable %s is not callable! (line %d)", node.Token, node.Token.Line))
	}

	return Node{}
}

func resolve(nodes []parser.Node, variables *map[string]Node, line uint) []Node {
	args := make([]Node, 0)

	for _, node := range nodes {
		if node.Type >= 5 && node.Type <= 10 {
			args = append(args, getLiteralFromNode(node))
		} else if node.Type == parser.IDENTIFIER {
			args = append(args, (*variables)[node.Token.Value])
		} else {
			panic(fmt.Sprintf("Unresolved value %s! (line %d)", node.Token.Value, line))
		}
	}

	return args
}

func getArgs(nodes []parser.Node, parameters []string, variables *map[string]Node, line uint) map[string]Node {
	targetMap := make(map[string]Node, 0)

	for k, v := range *variables {
		targetMap[k] = v
	}

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
	fmt.Println(vals)

	//TODO: Do calculation

	switch mode {
	case "*":
		break
	case "/":
		break
	case "+":
		break
	case "-":
		break
	}

	return Node{}
}

func doCall(node parser.Node, variables *map[string]Node) Node {
	if val, ok := (*variables)[node.Token.Value]; ok {
		return doVariableCall(node, val, variables)
	} else if node.Token.Value == SYSTEM_MATH {
		return doSystemCallMath(node, variables)
	}

	return Node{
		Value:    nil,
		NodeType: 0,
	}
}

// Run run an AST
func Run(nodes []parser.Node, variables map[string]Node) (returnVal Node) {
	for _, node := range nodes {
		switch node.Type {
		case parser.ACTION_DEF:
			returnVal = defineVariable(node, &variables)
			break
		case parser.ACTION_CALL:
			returnVal = doCall(node, &variables)
			break
		}
	}

	return
}
