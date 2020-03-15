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

func fail(expected string, got lexer.Token) {
	panic(fmt.Sprintf("Expected %s, got %s (line %d)!", expected, got.Value, got.Line))
}

func parseParameters(i *int, current *lexer.Token, tokens []lexer.Token) []Node {
	nodes := make([]Node, 0)

	for current.Type == lexer.IDENTIFIER {
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
		nodeType = LITERAL_BOOL
		break
	case lexer.ATOM:
		nodeType = LITERAL_ATOM
		break
	default:
		fail("a literal", current)
	}

	return Node{
		Type:      nodeType,
		Arguments: nil,
		Token:     current,
	}
}

func parseCall(i *int, current *lexer.Token, tokens []lexer.Token) Node {
	if current.Type != lexer.IDENTIFIER {
		fail("identifier", *current)
	}

	call := Node{
		Type:      ACTION_CALL,
		Arguments: make([]Node, 0),
		Token:     *current,
	}

	next(i, current, tokens)

	for current.Type != lexer.CLOSE_BRACE {
		if current.Type == lexer.OPEN_BRACE {
			next(i, current, tokens)
			call.Arguments = append(call.Arguments, parseCall(i, current, tokens))
		} else if current.Type == lexer.IDENTIFIER {
			call.Arguments = append(call.Arguments, parseParameters(i, current, tokens)...)
		} else if current.Type >= lexer.INT && current.Type <= lexer.ATOM {
			call.Arguments = append(call.Arguments, parseLiteral(*current))
			next(i, current, tokens)
		} else {
			fail("a valid function argument", *current)
		}
	}

	next(i, current, tokens)

	return call
}

func parseFunction(i *int, current *lexer.Token, tokens []lexer.Token) Node {
	// Current is already "fn"
	fn := Node{
		Type:      LITERAL_FN,
		Arguments: make([]Node, 0),
		Token:     lexer.Token{},
	}

	next(i, current, tokens)

	parameters := Node{
		Type:      PARAMETERS,
		Arguments: parseParameters(i, current, tokens),
		Token:     lexer.Token{},
	}

	fn.Arguments = append(fn.Arguments, parameters)

	for current.Type != lexer.CLOSE_BRACE {
		if current.Type != lexer.OPEN_BRACE {
			fail("(", *current)
		}
		next(i, current, tokens)
		fn.Arguments = append(fn.Arguments, doParse(i, current, tokens))
	}

	if current.Type != lexer.CLOSE_BRACE {
		fail(")", *current)
	}
	next(i, current, tokens)
	return fn
}

func doParse(i *int, current *lexer.Token, tokens []lexer.Token) (node Node) {
	if current.Type == lexer.IDENTIFIER_DEF {
		node = Node{
			Type:      ACTION_DEF,
			Arguments: make([]Node, 0),
			Token:     lexer.Token{},
		}
		next(i, current, tokens)

		node.Arguments = append(node.Arguments, parseParameters(i, current, tokens)...)

		if current.Type != lexer.OPEN_BRACE {
			node.Arguments = append(node.Arguments, parseLiteral(*current))
			next(i, current, tokens)
		} else {
			next(i, current, tokens)

			if current.Type == lexer.IDENTIFIER_FN {
				node.Arguments = append(node.Arguments, parseFunction(i, current, tokens))
			}

			//TODO: Parse struct, fn or call
		}

		if current.Type != lexer.CLOSE_BRACE {
			fail(")", *current)
		}
		return
	} else if current.Type == lexer.IDENTIFIER_FOR {
		//TODO Parse for
	} else if current.Type == lexer.IDENTIFIER_IF {
		//TODO Parse if
	} else if current.Type == lexer.ATOM {
		//TODO Parse map access
	} else {
		node = parseCall(i, current, tokens)
	}

	return
}

// Parse parse tokens and return an AST
func Parse(tokens []lexer.Token) []Node {
	nodes := make([]Node, 0)

	for i := 0; i < len(tokens); i++ {
		current := tokens[i]

		if current.Type != lexer.OPEN_BRACE {
			fail("(", current)
		}

		next(&i, &current, tokens)

		nodes = append(nodes, doParse(&i, &current, tokens))
	}

	return nodes
}
