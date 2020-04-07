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
	interpreter.CopyVariableState(variables, &ctx)

	switch mode {
	case "slice":
		return doSlice(args)
	case "len":
		return doLen(args)
	case "nth":
		return doNth(args)
	case "each":
		return doEach(ctx, args)
	case "map":
		return doMap(ctx, args)
	case "flatmap":
		return doFlatmap(args)
	case "filter":
		return doFilter(ctx, args)
	case "reduce":
		return doReduce(ctx, args)
	default:
		panic("Unrecognized mode")
	}
}

func doSlice(args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_INT && (args[1].NodeType != interpreter.NODETYPE_ATOM && args[1].Value != "-") {
		panic(CALL + " :slice expects an int or :- as first argument!")
	}

	if args[2].NodeType != interpreter.NODETYPE_INT && (args[2].NodeType != interpreter.NODETYPE_ATOM && args[2].Value != "-") {
		panic(CALL + " :slice expects an int or :- as second argument!")
	}

	if args[3].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " :slice expects a list as third argument!")
	}

	if args[1].NodeType == interpreter.NODETYPE_INT && args[2].NodeType == interpreter.NODETYPE_INT {
		return interpreter.NodeList(args[3].Value.(interpreter.ListNode).Values[args[1].Value.(int):args[2].Value.(int)])
	} else if args[1].NodeType == interpreter.NODETYPE_INT && args[2].NodeType == interpreter.NODETYPE_ATOM {
		return interpreter.NodeList(args[3].Value.(interpreter.ListNode).Values[args[1].Value.(int):])
	} else if args[1].NodeType == interpreter.NODETYPE_ATOM && args[2].NodeType == interpreter.NODETYPE_INT {
		return interpreter.NodeList(args[3].Value.(interpreter.ListNode).Values[:args[2].Value.(int)])
	} else {
		return args[3]
	}
}

func doLen(args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " :len expects a list as first argument!")
	}

	return interpreter.IntNode(len(args[1].Value.(interpreter.ListNode).Values))
}

func doNth(args []interpreter.Node) interpreter.Node {
	if args[2].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " :nth expects a list as second argument!")
	}
	list := args[2].Value.(interpreter.ListNode)

	if args[1].NodeType != interpreter.NODETYPE_INT {
		panic(CALL + " :nth expects an int as first argument!")
	}

	return list.Values[args[1].Value.(int)]
}

func doEach(ctx map[string]interpreter.Node, args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " :each expects a list as first argument!")
	}
	list := args[1].Value.(interpreter.ListNode)

	if args[2].NodeType != interpreter.NODETYPE_FN {
		panic(CALL + " :each expects a function as second argument!")
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

	return interpreter.Nothing
}

func doMap(ctx map[string]interpreter.Node, args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " :map expects a list as first argument!")
	}
	list := args[1].Value.(interpreter.ListNode)

	if args[2].NodeType != interpreter.NODETYPE_FN {
		panic(CALL + " :map expects a function as second argument!")
	}
	fn := args[2].Value.(interpreter.FnLiteral)

	if len(fn.Parameters) != 1 {
		panic("Enumerate map should have one parameter in execution function!")
	}

	mapResult := make([]interpreter.Node, 0)

	for _, value := range list.Values {
		ctx[SYSTEM_ENUMERATE_VAL] = value

		mapResult = append(mapResult, interpreter.DoVariableCall(parser.Node{
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

	return interpreter.NodeList(mapResult)
}

func doFlatmap(args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " :filter expects a list as first argument!")
	}

	return interpreter.NodeList(flatmapNode(args[1]))
}

func flatmapNode(arg interpreter.Node) []interpreter.Node {

	if arg.NodeType == interpreter.NODETYPE_LIST {
		list := make([]interpreter.Node, 0)

		for _, value := range arg.Value.(interpreter.ListNode).Values {
			list = append(list, flatmapNode(value)...)
		}

		return list
	}

	return []interpreter.Node{arg}
}

func doFilter(ctx map[string]interpreter.Node, args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " :filter expects a list as first argument!")
	}
	list := args[1].Value.(interpreter.ListNode)

	if args[2].NodeType != interpreter.NODETYPE_FN {
		panic(CALL + " :filter expects a function as second argument!")
	}
	fn := args[2].Value.(interpreter.FnLiteral)
	if len(fn.Parameters) != 1 {
		panic("Enumerate filter should have one parameter in execution function!")
	}

	filterResult := make([]interpreter.Node, 0)

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
			filterResult = append(filterResult, value)
		}
	}

	return interpreter.NodeList(filterResult)
}

func doReduce(ctx map[string]interpreter.Node, args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " :reduce expects a list as first argument!")
	}
	list := args[1].Value.(interpreter.ListNode)

	if args[2].NodeType != interpreter.NODETYPE_FN {
		panic(CALL + " :reduce expects a function as second argument!")
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
