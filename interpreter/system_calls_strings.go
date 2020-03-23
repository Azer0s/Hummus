package interpreter

import (
	"github.com/Azer0s/Hummus/parser"
	"strings"
)

func doSystemCallStrings(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

	mode := args[0].Value.(string)

	switch mode {
	case "concat":
		if args[1].NodeType != NODETYPE_LIST {
			panic(SYSTEM_STRING + " :concat only accepts lists!")
		}

		str := make([]string, 0)

		for _, value := range args[1].Value.(ListNode).Values {
			if value.NodeType != NODETYPE_STRING {
				panic(SYSTEM_STRING + " :concat only accepts lists of string!")
			}

			str = append(str, value.Value.(string))
		}

		return Node{
			Value:    strings.Join(str, ""),
			NodeType: NODETYPE_STRING,
		}
	case "len":
		if args[1].NodeType != NODETYPE_STRING {
			panic(SYSTEM_STRING + " :len only accepts strings!")
		}

		return Node{
			Value:    len(args[1].Value.(string)),
			NodeType: NODETYPE_INT,
		}
	case "nth":
		if args[1].NodeType != NODETYPE_STRING {
			panic(SYSTEM_STRING + " :nth expects a string as the first argument!")
		}

		if args[2].NodeType != NODETYPE_INT {
			panic(SYSTEM_STRING + " :nth expects an int as the second argument!")
		}

		return Node{
			Value:    string(args[1].Value.(string)[args[2].Value.(int)]),
			NodeType: NODETYPE_STRING,
		}
	default:
		panic("Unrecognized mode")
	}
}
