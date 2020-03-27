package main

import (
	"github.com/Azer0s/Hummus/interpreter"
)

// CALL the built in value make system function
var CALL string = "--system-do-make!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "list":
		if args[1].NodeType == interpreter.NODETYPE_LIST {
			return interpreter.Node{
				Value:    args[1].Value,
				NodeType: interpreter.NODETYPE_LIST,
			}
		}

		return interpreter.Node{
			Value:    interpreter.ListNode{Values: args[1:]},
			NodeType: interpreter.NODETYPE_LIST,
		}
	case "range":
		from := args[1]
		to := args[2]

		if from.NodeType != interpreter.NODETYPE_INT || to.NodeType != interpreter.NODETYPE_INT {
			panic("Expected an int as parameter for range!")
		}

		f := from.Value.(int)
		t := to.Value.(int)

		list := interpreter.ListNode{Values: make([]interpreter.Node, 0)}

		if f > t {
			for i := t; i >= t; i-- {
				list.Values = append(list.Values, interpreter.Node{
					Value:    i,
					NodeType: interpreter.NODETYPE_INT,
				})
			}
		} else {
			for i := f; i <= t; i++ {
				list.Values = append(list.Values, interpreter.Node{
					Value:    i,
					NodeType: interpreter.NODETYPE_INT,
				})
			}
		}

		return interpreter.Node{
			Value:    list,
			NodeType: interpreter.NODETYPE_LIST,
		}

	default:
		panic("Unrecognized mode")
	}
}
