package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"strconv"
	"strings"
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

	switch mode {
	case "string":
		return interpreter.Node{
			Value:    interpreter.DumpNode(args[1]),
			NodeType: interpreter.NODETYPE_STRING,
		}

	case "atom":
		val := interpreter.DumpNode(args[1])
		val = strings.ReplaceAll(val, " ", "_")
		val = strings.ReplaceAll(val, "(", "")
		val = strings.ReplaceAll(val, ")", "")

		return interpreter.Node{
			Value:    val,
			NodeType: interpreter.NODETYPE_ATOM,
		}

	case "int":
		val, err := strconv.Atoi(interpreter.DumpNode(args[1]))
		return interpreter.OptionalNode(val, interpreter.NODETYPE_INT, err != nil)

	case "float":
		val, err := strconv.ParseFloat(interpreter.DumpNode(args[1]), 64)
		return interpreter.OptionalNode(val, interpreter.NODETYPE_FLOAT, err != nil)

	case "bool":
		val, err := strconv.ParseBool(interpreter.DumpNode(args[1]))
		return interpreter.OptionalNode(val, interpreter.NODETYPE_BOOL, err != nil)

	case "identity":
		return args[1]

	default:
		panic("Unrecognized mode")
	}
}
