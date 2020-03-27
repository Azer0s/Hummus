package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
)

// CALL enumeration functions
var CALL string = "--system-do-enumerate!"

// SYSTEM_ENUMERATE_VAL variable where mfr values are stored
const SYSTEM_ENUMERATE_VAL string = "--system-do-enumerate-val"

// SYSTEM_ACCUMULATE_VAL variable where reduce state is stored
const SYSTEM_ACCUMULATE_VAL string = "--system-do-accumulate-val"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	ctx := make(map[string]interpreter.Node, 0)
	for k, v := range *variables {
		ctx[k] = v
	}

	switch mode {
	case "nth":
		if args[2].NodeType != interpreter.NODETYPE_LIST {
			panic(CALL + " expects a list as second argument!")
		}
		list := args[2].Value.(interpreter.ListNode)

		if args[1].NodeType != interpreter.NODETYPE_INT {
			panic(CALL + " :nth expects an int as first argument!")
		}

		return list.Values[args[1].Value.(int)]
	case "each":
		return doEach(ctx, args)
	case "map":
		return doMap(ctx, args)
	case "filter":
		return doFilter(ctx, args)
	case "reduce":
		return doReduce(ctx, args)
	default:
		panic("Unrecognized mode")
	}
}

func doEach(ctx map[string]interpreter.Node, args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " expects a list as first argument!")
	}
	list := args[1].Value.(interpreter.ListNode)

	if args[2].NodeType != interpreter.NODETYPE_FN {
		panic(CALL + " expects a function as second argument!")
	}
	fn := args[2].Value.(interpreter.FnLiteral)
	if len(fn.Parameters) != 1 {
		panic("Enumerate each should have one parameter in execution function!")
	}

	for _, value := range list.Values {
		ctx[SYSTEM_ENUMERATE_VAL] = value

		interpreter.DoVariableCall(parser.Node{
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

	return interpreter.Node{
		Value:    0,
		NodeType: 0,
	}
}

func doMap(ctx map[string]interpreter.Node, args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " expects a list as first argument!")
	}
	list := args[1].Value.(interpreter.ListNode)

	if args[2].NodeType != interpreter.NODETYPE_FN {
		panic(CALL + " expects a function as second argument!")
	}
	fn := args[2].Value.(interpreter.FnLiteral)

	if len(fn.Parameters) != 1 {
		panic("Enumerate map should have one parameter in execution function!")
	}

	mapResult := interpreter.ListNode{Values: make([]interpreter.Node, 0)}

	for _, value := range list.Values {
		ctx[SYSTEM_ENUMERATE_VAL] = value

		mapResult.Values = append(mapResult.Values, interpreter.DoVariableCall(parser.Node{
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

	return interpreter.Node{
		Value:    mapResult,
		NodeType: interpreter.NODETYPE_LIST,
	}
}

func doFilter(ctx map[string]interpreter.Node, args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " expects a list as first argument!")
	}
	list := args[1].Value.(interpreter.ListNode)

	if args[2].NodeType != interpreter.NODETYPE_FN {
		panic(CALL + " expects a function as second argument!")
	}
	fn := args[2].Value.(interpreter.FnLiteral)
	if len(fn.Parameters) != 1 {
		panic("Enumerate filter should have one parameter in execution function!")
	}

	filterResult := interpreter.ListNode{Values: make([]interpreter.Node, 0)}

	for _, value := range list.Values {
		ctx[SYSTEM_ENUMERATE_VAL] = value

		res := interpreter.DoVariableCall(parser.Node{
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

		if res.NodeType != interpreter.NODETYPE_BOOL {
			panic("Enumerate filter result must be a bool!")
		}

		if res.Value.(bool) {
			filterResult.Values = append(filterResult.Values, value)
		}
	}

	return interpreter.Node{
		Value:    filterResult,
		NodeType: interpreter.NODETYPE_LIST,
	}
}

func doReduce(ctx map[string]interpreter.Node, args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " expects a list as first argument!")
	}
	list := args[1].Value.(interpreter.ListNode)

	if args[2].NodeType != interpreter.NODETYPE_FN {
		panic(CALL + " expects a function as second argument!")
	}
	fn := args[2].Value.(interpreter.FnLiteral)
	if len(fn.Parameters) != 2 {
		panic("Enumerate reduce should have two parameters in execution function!")
	}

	ctx[SYSTEM_ACCUMULATE_VAL] = args[3]

	for _, value := range list.Values {
		ctx[SYSTEM_ENUMERATE_VAL] = value

		res := interpreter.DoVariableCall(parser.Node{
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
}
