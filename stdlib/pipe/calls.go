package main

import (
	"fmt"
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
)

var CALL string = "--system-do-pipe!"

func Init(variables *map[string]interpreter.Node) {
	// noinit
}

func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "combine":
		if args[1].NodeType != interpreter.NODETYPE_LIST {
			panic(CALL + " :concat only accepts lists!")
		}

		list := args[1].Value.(interpreter.ListNode)

		fns := make([]string, 0)
		ctx := make(map[string]interpreter.Node, 0)

		count := 0

		for _, value := range list.Values {
			fn := fmt.Sprintf("f%d", count)
			count++

			fns = append(fns, fn)
			ctx[fn] = value
		}

		reversed := make([]parser.Node, 0)

		// reverse order
		for i := range fns {
			n := fns[len(fns)-1-i]
			reversed = append(reversed, parser.Node{
				Type:      parser.ACTION_CALL,
				Arguments: nil,
				Token: lexer.Token{
					Value: n,
					Type:  lexer.IDENTIFIER,
					Line:  0,
				},
			})
		}

		reversed[len(reversed)-1].Arguments = make([]parser.Node, 0)
		for _, parameter := range list.Values[0].Value.(interpreter.FnLiteral).Parameters {
			reversed[len(reversed)-1].Arguments = append(reversed[len(reversed)-1].Arguments, parser.Node{
				Type:      parser.IDENTIFIER,
				Arguments: nil,
				Token: lexer.Token{
					Value: parameter,
					Type:  lexer.IDENTIFIER,
					Line:  0,
				},
			})
		}

		for i := len(reversed) - 1; i > 0; i-- {
			reversed[i-1].Arguments = make([]parser.Node, 0)
			reversed[i-1].Arguments = append(reversed[i-1].Arguments, reversed[i])
		}

		return interpreter.Node{
			Value: interpreter.FnLiteral{
				Parameters: list.Values[0].Value.(interpreter.FnLiteral).Parameters,
				Body: []parser.Node{
					reversed[0],
				},
				Context: ctx,
			},
			NodeType: interpreter.NODETYPE_FN,
		}
	default:
		panic("Unrecognized mode")
	}
}
