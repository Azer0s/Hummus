package parser

import (
	"encoding/json"
	"fmt"
	"github.com/Azer0s/Hummus/lexer"
)

func next(i *int, current *lexer.Token, tokens []lexer.Token) {
	*i += 1

	if *i >= len(tokens) {
		*current = lexer.Token{}
		return
	}

	*current = tokens[*i]
}

func parseParameter(i *int, current *lexer.Token, tokens []lexer.Token) Node {
	return Node{
		Type:      0,
		Arguments: nil,
		Value:     lexer.Token{},
	}
}

func Parse(tokens []lexer.Token) []Node {
	b, _ := json.Marshal(tokens)
	fmt.Println(string(b))
	nodes := make([]Node, 0)

	for i := 0; i < len(tokens); i++ {
		current := tokens[i]

		if current.Type != lexer.OPEN_BRACE {
			panic(fmt.Sprintf("Unexpected token %s (line %d)!", current.Value, current.Line))
		}

		next(&i, &current, tokens)

		if current.Type == lexer.IDENTIFIER_DEF {
			node := Node{
				Type:      ACTION_DEF,
				Arguments: make([]Node, 0),
				Value:     lexer.Token{},
			}
			next(&i, &current, tokens)

			//TODO: Implement assignment
			node.Arguments = append(node.Arguments, parseParameter(&i, &current, tokens))
		}
	}

	return nodes
}
