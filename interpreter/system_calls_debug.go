package interpreter

import (
	"fmt"
	"github.com/Azer0s/Hummus/parser"
)

func doSystemCallDebug(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

	mode := args[0].Value.(string)

	switch mode {
	case "dump":
		fmt.Println(DumpNode(args[1]))
		return Node{
			Value:    0,
			NodeType: 0,
		}
	case "dump_state":
		for k, v := range *variables {
			fmt.Println(fmt.Sprintf("%s => %s", k, DumpNode(v)))
		}
		return Node{
			Value:    0,
			NodeType: 0,
		}
	case "line":
		return Node{
			Value:    int(node.Token.Line),
			NodeType: 0,
		}
	default:
		panic("Unrecognized mode")
	}
}
