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
		if args[1].NodeType != interpreter.NODETYPE_LIST {
			panic(CALL + " :concat only accepts lists!")
		}

		str := make([]string, 0)

		for _, value := range args[1].Value.(interpreter.ListNode).Values {
			if value.NodeType != interpreter.NODETYPE_STRING {
				panic(CALL + " :concat only accepts lists of string!")
			}

			str = append(str, value.Value.(string))
		}

		return interpreter.Node{
			Value:    strings.Join(str, ""),
			NodeType: interpreter.NODETYPE_STRING,
		}
	case "len":
		if args[1].NodeType != interpreter.NODETYPE_STRING {
			panic(CALL + " :len only accepts strings!")
		}

		return interpreter.Node{
			Value:    len(args[1].Value.(string)),
			NodeType: interpreter.NODETYPE_INT,
		}
	case "nth":
		if args[1].NodeType != interpreter.NODETYPE_INT {
			panic(CALL + " :nth expects an int as the second argument!")
		}

		if args[2].NodeType != interpreter.NODETYPE_STRING {
			panic(CALL + " :nth expects a string as the first argument!")
		}

		return interpreter.Node{
			Value:    string(args[2].Value.(string)[args[1].Value.(int)]),
			NodeType: interpreter.NODETYPE_STRING,
		}
	default:
		panic("Unrecognized mode")
	}
}
