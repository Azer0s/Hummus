package main

import (
	"fmt"

	"github.com/Azer0s/Hummus/interpreter"
)

func main() {}

// CALL debug functions
var CALL string = "nativeCall!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	fmt.Println("Hello from native!")
	return interpreter.Nothing
}
