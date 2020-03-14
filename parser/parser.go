package parser

import (
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

func parseParameters(i *int, current *lexer.Token, tokens []lexer.Token) []Node {
	nodes := make([]Node, 0)

	for (*current).Type == lexer.IDENTIFIER {
		nodes = append(nodes, Node{
			Type:      IDENTIFIER,
			Arguments: nil,
			Token:     *current,
		})

		next(i, current, tokens)
	}

	return nodes
}

func parseLiteral(current lexer.Token) Node {
	var nodeType NodeType

	switch current.Type {
	case lexer.INT:
		nodeType = LITERAL_INT
		break
	case lexer.FLOAT:
		nodeType = LITERAL_FLOAT
		break
	case lexer.STRING:
		nodeType = LITERAL_STRING
		break
	case lexer.BOOL:
		nodeType = LITERAL_BOOl
		break
	case lexer.ATOM:
		nodeType = LITERAL_ATOM
		break
	}

	return Node{
		Type:      nodeType,
		Arguments: nil,
		Token:     current,
	}
}

func Parse(tokens []lexer.Token) []Node {
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
				Token:     lexer.Token{},
			}
			next(&i, &current, tokens)

			node.Arguments = append(node.Arguments, parseParameters(&i, &current, tokens)...)

			if current.Type != lexer.OPEN_BRACE {
				node.Arguments = append(node.Arguments, parseLiteral(current))
				next(&i, &current, tokens)
			} else {
				//TODO: Parse struct, fn or call
			}

			nodes = append(nodes, node)
		}
	}

	return nodes
}
