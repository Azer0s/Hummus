package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"strings"
)

// CALL string functions
var CALL string = "--system-do-strings!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "concat":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_LIST, CALL+" :concat")

		str := make([]string, 0)

		for _, value := range args[1].Value.(interpreter.ListNode).Values {
			if value.NodeType != interpreter.NODETYPE_STRING {
				panic(CALL + " :concat only accepts lists of string!")
			}

			str = append(str, value.Value.(string))
		}

		return interpreter.StringNode(strings.Join(str, ""))
	case "len":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :len")
		return interpreter.IntNode(len(args[1].Value.(string)))
	case "nth":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_INT, CALL+" :nth")
		interpreter.EnsureType(&args, 2, interpreter.NODETYPE_STRING, CALL+" :nth")

		return interpreter.StringNode(string(args[2].Value.(string)[args[1].Value.(int)]))
	case "lower":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :lower")
		return interpreter.StringNode(strings.ToLower(args[1].Value.(string)))
	case "upper":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :upper")
		return interpreter.StringNode(strings.ToUpper(args[1].Value.(string)))
	default:
		panic("Unrecognized mode")
	}
}
