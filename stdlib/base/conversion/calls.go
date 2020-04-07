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
		return interpreter.StringNode(interpreter.DumpNode(args[1]))

	case "atom":
		val := interpreter.DumpNode(args[1])
		val = strings.ReplaceAll(val, " ", "_")
		val = strings.ReplaceAll(val, "(", "")
		val = strings.ReplaceAll(val, ")", "")

		return interpreter.AtomNode(val)

	case "int":
		val, err := strconv.Atoi(interpreter.DumpNode(args[1]))
		return interpreter.OptionalNode(interpreter.IntNode(val), err != nil)

	case "float":
		val, err := strconv.ParseFloat(interpreter.DumpNode(args[1]), 64)
		return interpreter.OptionalNode(interpreter.FloatNode(val), err != nil)

	case "bool":
		val, err := strconv.ParseBool(interpreter.DumpNode(args[1]))
		return interpreter.OptionalNode(interpreter.BoolNode(val), err != nil)

	case "identity":
		return args[1]

	default:
		panic("Unrecognized mode")
	}
}
