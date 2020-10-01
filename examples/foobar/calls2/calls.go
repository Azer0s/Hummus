package main

import (
	"github.com/Azer0s/Hummus/interpreter"
)

// CALL debug functions
var CALL string = "calls2!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	doPrint()
	return interpreter.Nothing
}
