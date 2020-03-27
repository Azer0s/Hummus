package main

import (
	"fmt"
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"net/http"
)

// CALL HTTP server functions
var CALL string = "--system-do-http!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "static":
		if args[1].NodeType != interpreter.NODETYPE_STRING {
			panic(CALL + " :static expects a string as the first argument!")
		}

		if args[2].NodeType != interpreter.NODETYPE_STRING {
			panic(CALL + " :static expects a string as the second argument!")
		}

		fs := http.FileServer(http.Dir(args[2].Value.(string)))
		http.Handle(args[1].Value.(string), fs)

		return interpreter.Node{
			Value:    0,
			NodeType: 0,
		}
	case "handle":
		if args[1].NodeType != interpreter.NODETYPE_STRING {
			panic(CALL + " :handle expects a string as the first argument!")
		}

		if args[2].NodeType != interpreter.NODETYPE_FN {
			panic(CALL + " :handle expects a fn as the second argument!")
		}

		ctx := make(map[string]interpreter.Node, 0)
		for k, v := range *variables {
			ctx[k] = v
		}

		fn := args[2]

		http.HandleFunc(args[1].Value.(string), func(writer http.ResponseWriter, request *http.Request) {
			req := "request"

			callCtx := make(map[string]interpreter.Node, 0)
			for k, v := range *variables {
				callCtx[k] = v
			}

			callCtx[req] = interpreter.Node{
				Value: interpreter.MapNode{Values: map[string]interpreter.Node{
					"method": {
						Value:    request.Method,
						NodeType: interpreter.NODETYPE_STRING,
					},
					"body": {
						Value:    request.Body,
						NodeType: interpreter.NODETYPE_STRING,
					},
				}},
				NodeType: interpreter.NODETYPE_MAP,
			}

			res := interpreter.DoVariableCall(parser.Node{
				Type: 0,
				Arguments: []parser.Node{
					{
						Type:      parser.IDENTIFIER,
						Arguments: nil,
						Token: lexer.Token{
							Value: req,
							Type:  0,
							Line:  0,
						},
					},
				},
				Token: lexer.Token{},
			}, fn, &callCtx)

			fmt.Fprintf(writer, interpreter.DumpNode(res))
		})

		return interpreter.Node{
			Value:    0,
			NodeType: 0,
		}
	case "serve":
		if args[1].NodeType != interpreter.NODETYPE_STRING {
			panic(CALL + " :serve expects a string as the first argument!")
		}

		err := http.ListenAndServe(args[1].Value.(string), nil)

		if err != nil {
			panic(err)
		}

		return interpreter.Node{
			Value:    0,
			NodeType: 0,
		}
	default:
		panic("Unrecognized mode")
	}
}
