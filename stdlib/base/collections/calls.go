package main

import (
	"github.com/Azer0s/Hummus/interpreter"
)

// CALL the built in value make system function
var CALL string = "--system-do-make!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "list":
		return list(args)

	case "exists":
		return exists(args)

	case "keys":
		return keys(args)

	case "range":
		return doRange(args)

	default:
		panic("Unrecognized mode")
	}
}

func list(args []interpreter.Node) interpreter.Node {
	if args[1].NodeType == interpreter.NODETYPE_LIST {
		return interpreter.Node{
			Value:    args[1].Value,
			NodeType: interpreter.NODETYPE_LIST,
		}
	}

	return interpreter.Node{
		Value:    interpreter.ListNode{Values: args[1:]},
		NodeType: interpreter.NODETYPE_LIST,
	}
}

func exists(args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_ATOM {
		panic("Expected an atom as parameter for exists? !")
	}

	if args[2].NodeType != interpreter.NODETYPE_MAP {
		panic("Expected a map as parameter for exists? !")
	}

	_, ok := args[2].Value.(interpreter.MapNode).Values[args[1].Value.(string)]

	return interpreter.Node{
		Value:    ok,
		NodeType: interpreter.NODETYPE_BOOL,
	}
}

func keys(args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_MAP {
		panic("Expected a map as parameter for keys!")
	}

	keys := make([]interpreter.Node, 0)
	for s := range args[1].Value.(interpreter.MapNode).Values {
		keys = append(keys, interpreter.Node{
			Value:    s,
			NodeType: interpreter.NODETYPE_ATOM,
		})
	}

	return interpreter.Node{
		Value:    interpreter.ListNode{Values: keys},
		NodeType: interpreter.NODETYPE_LIST,
	}
}

func doRange(args []interpreter.Node) interpreter.Node {
	from := args[1]
	to := args[2]

	if from.NodeType != interpreter.NODETYPE_INT || to.NodeType != interpreter.NODETYPE_INT {
		panic("Expected an int as parameter for range!")
	}

	f := from.Value.(int)
	t := to.Value.(int)

	list := interpreter.ListNode{Values: make([]interpreter.Node, 0)}

	if f > t {
		for i := t; i >= t; i-- {
			list.Values = append(list.Values, interpreter.Node{
				Value:    i,
				NodeType: interpreter.NODETYPE_INT,
			})
		}
	} else {
		for i := f; i <= t; i++ {
			list.Values = append(list.Values, interpreter.Node{
				Value:    i,
				NodeType: interpreter.NODETYPE_INT,
			})
		}
	}

	return interpreter.Node{
		Value:    list,
		NodeType: interpreter.NODETYPE_LIST,
	}
}
