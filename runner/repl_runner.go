package runner

import (
	"bufio"
	"fmt"
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"os"
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

// RunRepl runs the run-eval-print-loop
func RunRepl() {
	vars := make(map[string]interpreter.Node, 0)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("=> ")

		text := ""
		var tokens []lexer.Token

		for {
			t, _ := reader.ReadString('\n')

			text += t
			tokens = lexer.LexString(text)

			if checkClose(tokens) {
				break
			}
			fmt.Print(" ..")
		}

		var nodes []parser.Node

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
			fmt.Println(fmt.Sprintf("%v", interpreter.Run(nodes, &vars).Value))
		}
	}
}
