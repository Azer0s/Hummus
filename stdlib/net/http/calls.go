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
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :static")
		interpreter.EnsureType(&args, 2, interpreter.NODETYPE_STRING, CALL+" :static")

		fs := http.FileServer(http.Dir(args[2].Value.(string)))
		http.Handle(args[1].Value.(string), fs)

		return interpreter.Nothing
	case "handle":
		return doHandle(args, variables)
	case "serve":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :serve")

		err := http.ListenAndServe(args[1].Value.(string), nil)

		if err != nil {
			panic(err)
		}

		return interpreter.Nothing
	default:
		panic("Unrecognized mode")
	}
}

func doHandle(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :handle")
	interpreter.EnsureType(&args, 2, interpreter.NODETYPE_FN, CALL+" :handle")

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

	return interpreter.Nothing
}

func getReqMap(request *http.Request) interpreter.Node {
	return interpreter.NodeMap(map[string]interpreter.Node{
		"method": interpreter.StringNode(request.Method),
		"proto":  interpreter.StringNode(request.Proto),
		"header": getHeaders(request),
		"body":   getBody(request),
		"params": getParams(request),
	})
}

func getParams(request *http.Request) interpreter.Node {
	return getNodeByMapStringArray(request.URL.Query())
}

func getBody(request *http.Request) interpreter.Node {
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		return interpreter.StringNode("")
	}

	return interpreter.StringNode(string(body))
}

func getHeaders(request *http.Request) interpreter.Node {
	return getNodeByMapStringArray(request.Header)
}

func getNodeByMapStringArray(values map[string][]string) interpreter.Node {
	node := make(map[string]interpreter.Node, 0)

	for k, strings := range values {
		vals := make([]interpreter.Node, 0)

		for _, s := range strings {
			vals = append(vals, interpreter.StringNode(s))
		}

		node[k] = interpreter.NodeList(vals)
	}

	return interpreter.NodeMap(node)
}
