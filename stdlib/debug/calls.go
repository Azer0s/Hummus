package main

import (
	"fmt"
	"github.com/Azer0s/Hummus/interpreter"
)

// CALL debug functions
var CALL string = "--system-do-debug!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "dump":
		fmt.Println(interpreter.DumpNode(args[1]))
		return interpreter.Nothing
	case "dump_state":
		for k, v := range *variables {
			fmt.Println(fmt.Sprintf("%s => %s", k, interpreter.DumpNode(v)))
		}
		return interpreter.Nothing
	default:
		panic("Unrecognized mode")
	}
}
