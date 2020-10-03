package main

import (
	"fmt"
	"github.com/Azer0s/Hummus/interpreter"
	"strconv"
)

func main() {}

// CALL string functions
var CALL string = "--system-do-bytes!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
}

func getByteArrayByPseudoPtr(pseudoPtr int) []byte {
	b, ok := interpreter.LoadObject(pseudoPtr)

	if !ok {
		panic("Byte array doesn't exist or was already freed!")
	}

	return b.([]byte)
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "from-string":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :from-string")
		return interpreter.IntNode(interpreter.StoreObject([]byte(args[1].Value.(string))))

	case "to-atoms":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_INT, CALL+" :to-atoms")

		l := make([]interpreter.Node, 0)
		for _, b2 := range getByteArrayByPseudoPtr(args[1].Value.(int)) {
			l = append(l, interpreter.AtomNode("0x"+fmt.Sprintf("%x", b2)))
		}

		return interpreter.NodeList(l)

	case "to-string":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_INT, CALL+" :to-string")
		return interpreter.StringNode(string(getByteArrayByPseudoPtr(args[1].Value.(int))))

	case "from-atoms":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_LIST, CALL+" :from-atoms")

		bytes := make([]byte, 0)
		for i, value := range args[1].Value.(interpreter.ListNode).Values {
			interpreter.EnsureSingleType(&value, i+2, interpreter.NODETYPE_ATOM, CALL+" :from-atoms")
			b, err := strconv.ParseUint(value.Value.(string), 0, 8)

			if err != nil {
				panic(err)
			}

			bytes = append(bytes, byte(b))
		}

		return interpreter.IntNode(interpreter.StoreObject(bytes))
	default:
		panic("Unrecognized mode")
	}
}
