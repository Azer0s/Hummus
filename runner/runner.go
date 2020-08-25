package runner

import (
	"fmt"
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"github.com/carmark/pseudo-terminal-go/terminal"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func printHelp() {
	fmt.Println(`
Hummus REPL v1.0

(exit) .. Exits the REPL
(help) .. Prints the help page
 `)
}

func checkClose(tokens []lexer.Token) bool {
	buffer := 0

	for _, token := range tokens {
		if token.Type == lexer.CLOSE_BRACE {
			buffer--
		}

		if token.Type == lexer.OPEN_BRACE {
			buffer++
		}
	}

	return buffer == 0
}

func dumpRepl(node interpreter.Node) string {
	ret := ""

	if node.NodeType == interpreter.NODETYPE_LIST {
		ret += "("

		for _, value := range node.Value.(interpreter.ListNode).Values {
			ret += dumpRepl(value) + " "
		}

		ret = strings.TrimSuffix(ret, " ")
		ret += ")"
	} else if node.NodeType == interpreter.NODETYPE_MAP {
		ret += "("

		for k, v := range node.Value.(interpreter.MapNode).Values {
			ret += fmt.Sprintf("%s => %s ", k, dumpRepl(v))
		}

		ret = strings.TrimSuffix(ret, " ")
		ret += ")"
	} else if node.NodeType == interpreter.NODETYPE_FN {
		ret += "[fn "

		for _, parameter := range node.Value.(interpreter.FnLiteral).Parameters {
			ret += parameter + " "
		}

		ret = strings.TrimSuffix(ret, " ")
		ret += "]"
	} else if node.NodeType == interpreter.NODETYPE_STRUCT {
		ret += "[struct "

		for _, parameter := range node.Value.(interpreter.StructDef).Parameters {
			ret += parameter + " "
		}

		ret = strings.TrimSuffix(ret, " ")
		ret += "]"
	} else if node.NodeType == interpreter.NODETYPE_STRING {
		ret = fmt.Sprintf("\"%v\"", node.Value)
	} else if node.NodeType == interpreter.NODETYPE_ATOM {
		ret = fmt.Sprintf(":%v", node.Value)
	} else {
		ret = fmt.Sprintf("%v", node.Value)
	}

	return ret
}

func getParserType(tokenType lexer.TokenType) parser.NodeType {
	switch tokenType {
	case lexer.IDENTIFIER:
		return parser.IDENTIFIER
	case lexer.INT:
		return parser.LITERAL_INT
	case lexer.FLOAT:
		return parser.LITERAL_FLOAT
	case lexer.STRING:
		return parser.LITERAL_STRING
	case lexer.BOOL:
		return parser.LITERAL_BOOL
	case lexer.ATOM:
		return parser.LITERAL_ATOM
	}

	return 0
}

func doFailsafeRepl(term *terminal.Terminal, vars *map[string]interpreter.Node) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:", r)
		}
	}()

	text := ""
	var tokens []lexer.Token

	term.SetPrompt("=> ")

	for {
		t, _ := term.ReadLine()

		text += t + "\n"
		tokens = lexer.LexString(text)

		if checkClose(tokens) {
			break
		}

		term.SetPrompt(" ..")
	}

	var nodes []parser.Node

	if len(tokens) == 0 {
		return
	}

	if len(tokens) == 1 && tokens[0].Type == lexer.IDENTIFIER || tokens[0].Type >= lexer.INT && tokens[0].Type <= lexer.ATOM {
		nodes = []parser.Node{
			{
				Type:      getParserType(tokens[0].Type),
				Arguments: nil,
				Token:     tokens[0],
			},
		}
	} else {
		nodes = parser.Parse(tokens)
	}

	if len(nodes) == 1 && nodes[0].Type == parser.ACTION_CALL && nodes[0].Token.Value == "exit" {
		os.Exit(0)
	} else if len(nodes) == 1 && nodes[0].Type == parser.ACTION_CALL && nodes[0].Token.Value == "help" {
		printHelp()
	} else {
		fmt.Println(dumpRepl(interpreter.Run(nodes, vars)))
	}
}

// RunRepl runs the run-eval-print-loop
func RunRepl() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	vars := make(map[string]interpreter.Node, 0)
	setupVars(&vars, dir)

	if err != nil {
		panic(err)
	}

	term, err := terminal.NewWithStdInOut()
	if err != nil {
		panic(err)
	}
	defer term.ReleaseFromStdInOut()

	for {
		doFailsafeRepl(term, &vars)
	}
}

func getExecFile(p, filename string) string {
	if path.IsAbs(filename) {
		return filename
	}

	return path.Join(p, filename)
}

// RunFile runs a file by filename
func RunFile(filename string) interpreter.Node {
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	execFile := getExecFile(p, filename)

	vars := make(map[string]interpreter.Node, 0)
	setupVars(&vars, execFile)

	return interpreter.Run(parser.Parse(lexer.LexString(string(b))), &vars)
}

// RunString run Hummus code from a string
func RunString(code string) interpreter.Node {
	vars := make(map[string]interpreter.Node, 0)
	setupVars(&vars, os.Args[0])

	return interpreter.Run(parser.Parse(lexer.LexString(code)), &vars)
}
