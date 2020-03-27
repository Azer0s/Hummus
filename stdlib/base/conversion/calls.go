package main

import (
	"github.com/Azer0s/Hummus/interpreter"
)

// CALL conversion functions
var CALL string = "--system-do-convert!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	if args[1].NodeType == interpreter.NODETYPE_LIST {
		panic(CALL + " doesn't accept lists!")
	}

	switch mode {
	case "string":
		return interpreter.Node{
			Value:    interpreter.DumpNode(args[1]),
			NodeType: interpreter.NODETYPE_STRING,
		}
	case "identity":
		return args[1]
	default:
		panic("Unrecognized mode")
	}
}
