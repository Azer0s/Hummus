package interpreter

import (
	"crypto/rand"
	"fmt"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"log"
)

func doSystemCallPipe(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

	mode := args[0].Value.(string)

	switch mode {
	case "combine":
		if args[1].NodeType != NODETYPE_LIST {
			panic(SYSTEM_STRING + " :concat only accepts lists!")
		}

		list := args[1].Value.(ListNode)

		fns := make([]string, 0)
		ctx := make(map[string]Node, 0)

		count := 0

		for _, value := range list.Values {
			b := make([]byte, 16)
			_, err := rand.Read(b)
			if err != nil {
				log.Fatal(err)
			}
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
		for _, parameter := range list.Values[0].Value.(FnLiteral).Parameters {
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

		return Node{
			Value: FnLiteral{
				Parameters: list.Values[0].Value.(FnLiteral).Parameters,
				Body: []parser.Node{
					reversed[0],
				},
				Context: ctx,
			},
			NodeType: NODETYPE_FN,
		}
	default:
		panic("Unrecognized mode")
	}
}
