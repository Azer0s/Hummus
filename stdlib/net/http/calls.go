package main

import (
	"fmt"
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"io/ioutil"
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
		return doHandle(args, variables)
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

func doHandle(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_STRING {
		panic(CALL + " :handle expects a string as the first argument!")
	}

	if args[2].NodeType != interpreter.NODETYPE_FN {
		panic(CALL + " :handle expects a fn as the second argument!")
	}

	ctx := make(map[string]interpreter.Node, 0)
	interpreter.CopyVariableState(variables, &ctx)

	fn := args[2]

	http.HandleFunc(args[1].Value.(string), func(writer http.ResponseWriter, request *http.Request) {
		req := "request"

		callCtx := make(map[string]interpreter.Node, 0)
		interpreter.CopyVariableState(variables, &callCtx)

		callCtx[req] = getReqMap(request)

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
}

func getReqMap(request *http.Request) interpreter.Node {
	return interpreter.Node{
		Value: interpreter.MapNode{Values: map[string]interpreter.Node{
			"method": {
				Value:    request.Method,
				NodeType: interpreter.NODETYPE_STRING,
			},
			"proto": {
				Value:    request.Proto,
				NodeType: interpreter.NODETYPE_STRING,
			},
			"header": getHeaders(request),
			"body":   getBody(request),
			"params": getParams(request),
		}},
		NodeType: interpreter.NODETYPE_MAP,
	}
}

func getParams(request *http.Request) interpreter.Node {
	return getNodeByMapStringArray(request.URL.Query())
}

func getBody(request *http.Request) interpreter.Node {
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		return interpreter.Node{
			Value:    "",
			NodeType: interpreter.NODETYPE_STRING,
		}
	}

	return interpreter.Node{
		Value:    string(body),
		NodeType: interpreter.NODETYPE_STRING,
	}
}

func getHeaders(request *http.Request) interpreter.Node {
	return getNodeByMapStringArray(request.Header)
}

func getNodeByMapStringArray(values map[string][]string) interpreter.Node {
	node := make(map[string]interpreter.Node, 0)

	for k, strings := range values {
		vals := interpreter.ListNode{Values: make([]interpreter.Node, 0)}

		for _, s := range strings {
			vals.Values = append(vals.Values, interpreter.Node{
				Value:    s,
				NodeType: interpreter.NODETYPE_STRING,
			})
		}

		node[k] = interpreter.Node{
			Value:    vals,
			NodeType: interpreter.NODETYPE_LIST,
		}
	}

	return interpreter.Node{
		Value:    interpreter.MapNode{Values: node},
		NodeType: interpreter.NODETYPE_MAP,
	}
}
