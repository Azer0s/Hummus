package main

import (
	"bufio"
	"fmt"
	"github.com/Azer0s/Hummus/interpreter"
	"os"
)

var reader = bufio.NewReader(os.Stdin)

// CALL the built in io function
var CALL string = "--system-do-io!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "console-out":
		if args[1].NodeType <= interpreter.NODETYPE_ATOM {
			fmt.Print(args[1].Value)
		} else {
			panic(CALL + " :console-out only accepts int, float, bool, string or atom!")
		}

		return interpreter.Node{
			Value:    0,
			NodeType: 0,
		}
	case "console-in":
		t, _ := reader.ReadString('\n')

		return interpreter.Node{
			Value:    t,
			NodeType: interpreter.NODETYPE_STRING,
		}
	case "console-clear":
		print("\033[H\033[2J")

		return interpreter.Node{
			Value:    0,
			NodeType: 0,
		}
	default:
		panic("Unrecognized mode")
	}
}
