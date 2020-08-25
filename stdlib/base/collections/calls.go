package main

import (
	"github.com/Azer0s/Hummus/interpreter"
)

func main() {}

// CALL the built in value make system function
var CALL string = "--system-do-collections!"

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

	case "nothing":
		return interpreter.NodeList(make([]interpreter.Node, 0))

	case "map-add":
		return doMapAdd(args)

	case "list-add":
		return doListAdd(args)

	default:
		panic("Unrecognized mode")
	}
}

func doMapAdd(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_MAP, CALL+" :map-add")
	interpreter.EnsureType(&args, 2, interpreter.NODETYPE_ATOM, CALL+" :map-add")

	m := args[1].Value.(interpreter.MapNode).Values
	m[args[2].Value.(string)] = args[3]

	return interpreter.NodeMap(m)
}

func doListAdd(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_LIST, CALL+" :map-add")

	l := args[1].Value.(interpreter.ListNode).Values
	l = append(l, args[2])

	return interpreter.NodeList(l)
}

func list(args []interpreter.Node) interpreter.Node {
	if args[1].NodeType == interpreter.NODETYPE_LIST {
		return interpreter.Node{
			//I don't use interpreter.NodeList here because I'd have to cast args[1].Value to a slice
			Value:    args[1].Value,
			NodeType: interpreter.NODETYPE_LIST,
		}
	}

	return interpreter.NodeList(args[1:])
}

func exists(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_ATOM, CALL+" :exists")
	interpreter.EnsureType(&args, 2, interpreter.NODETYPE_MAP, CALL+" :exists")

	_, ok := args[2].Value.(interpreter.MapNode).Values[args[1].Value.(string)]

	return interpreter.BoolNode(ok)
}

func keys(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_MAP, CALL+" :keys")

	keys := make([]interpreter.Node, 0)
	for s := range args[1].Value.(interpreter.MapNode).Values {
		keys = append(keys, interpreter.AtomNode(s))
	}

	return interpreter.NodeList(keys)
}

func doRange(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_INT, CALL+" :range")
	interpreter.EnsureType(&args, 2, interpreter.NODETYPE_INT, CALL+" :range")

	f := args[1].Value.(int)
	t := args[2].Value.(int)

	list := make([]interpreter.Node, 0)

	if f > t {
		for i := t; i >= t; i-- {
			list = append(list, interpreter.IntNode(i))
		}
	} else {
		for i := f; i <= t; i++ {
			list = append(list, interpreter.IntNode(i))
		}
	}

	return interpreter.NodeList(list)
}
