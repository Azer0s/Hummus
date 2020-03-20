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

func expectClose(i *int, current *lexer.Token, tokens []lexer.Token) {
	if current.Type != lexer.CLOSE_BRACE {
		fail(")", *current)
	}

	next(i, current, tokens)
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

func parseMapAccess(i *int, current *lexer.Token, tokens []lexer.Token, allowLiterals bool) (node Node) {
	node = Node{
		Type:      ACTION_MAP_ACCESS,
		Arguments: make([]Node, 0),
		Token:     lexer.Token{},
	}
	node.Arguments = append(node.Arguments, parseLiteral(*current))
	next(i, current, tokens)

	if current.Type == lexer.IDENTIFIER {
		node.Arguments = append(node.Arguments, Node{
			Type:      IDENTIFIER,
			Arguments: nil,
			Token:     *current,
		})
		next(i, current, tokens)
	} else if current.Type == lexer.OPEN_BRACE {
		next(i, current, tokens)
		node.Arguments = append(node.Arguments, parseCall(i, current, tokens))
	} else if current.Type >= lexer.INT && current.Type <= lexer.ATOM && allowLiterals {
		node.Arguments = append(node.Arguments, parseLiteral(*current))
		next(i, current, tokens)
	} else {
		fail("a valid map access", *current)
	}

	expectClose(i, current, tokens)

	return
}

func parseMap(i *int, current *lexer.Token, tokens []lexer.Token) (node Node) {
	node = Node{
		Type:      ACTION_MAP,
		Arguments: make([]Node, 0),
		Token:     lexer.Token{},
	}

	next(i, current, tokens)

	for current.Type != lexer.CLOSE_BRACE {
		next(i, current, tokens)
		node.Arguments = append(node.Arguments, parseMapAccess(i, current, tokens, true))
	}

	expectClose(i, current, tokens)

	return
}

func parseCall(i *int, current *lexer.Token, tokens []lexer.Token) Node {
	if current.Value == "{}" {
		return parseMap(i, current, tokens)
	}

	if current.Type == lexer.IDENTIFIER_IF {
		return parseIf(i, current, tokens)
	}

	call := Node{
		Type:      ACTION_CALL,
		Arguments: make([]Node, 0),
		Token:     *current,
	}

	if current.Type == lexer.OPEN_BRACE {
		next(i, current, tokens)
		if current.Type != lexer.IDENTIFIER_FN {
			fail("fn", *current)
		}

		call.Arguments = append(call.Arguments, parseFunction(i, current, tokens))
		*i--

		call.Token = lexer.Token{
			Type: lexer.ANONYMOUS_FN,
		}
	} else if current.Type != lexer.IDENTIFIER {
		fail("identifier", *current)
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
		if current.Type >= lexer.INT && current.Type <= lexer.ATOM {
			fn.Arguments = append(fn.Arguments, parseLiteral(*current))
			next(i, current, tokens)

			if current.Type != lexer.CLOSE_BRACE {
				fail("a closing brace after a literal return", *current)
			}
			continue
		} else if current.Type == lexer.IDENTIFIER {
			fn.Arguments = append(fn.Arguments, Node{
				Type:      IDENTIFIER,
				Arguments: nil,
				Token:     *current,
			})

			next(i, current, tokens)

			if current.Type != lexer.CLOSE_BRACE {
				fail("a closing brace after an variable return", *current)
			}
			continue
		}

		if current.Type != lexer.OPEN_BRACE {
			fail("(", *current)
		}
		next(i, current, tokens)
		fn.Arguments = append(fn.Arguments, doParse(i, current, tokens, false))
	}

	expectClose(i, current, tokens)
	return fn
}

func parseIf(i *int, current *lexer.Token, tokens []lexer.Token) (node Node) {
	node = Node{
		Type:      ACTION_IF,
		Arguments: make([]Node, 0),
		Token:     lexer.Token{},
	}

	next(i, current, tokens)

	if current.Type == lexer.BOOL {
		node.Arguments = append(node.Arguments, parseLiteral(*current))
		next(i, current, tokens)
	} else if current.Type == lexer.IDENTIFIER {
		node.Arguments = append(node.Arguments, Node{
			Type:      IDENTIFIER,
			Arguments: nil,
			Token:     *current,
		})
		next(i, current, tokens)
	} else if current.Type == lexer.OPEN_BRACE {
		next(i, current, tokens)
		node.Arguments = append(node.Arguments, parseCall(i, current, tokens))
	}

	node.Arguments = append(node.Arguments, parseBranch(i, current, tokens))

	if current.Type != lexer.CLOSE_BRACE {
		node.Arguments = append(node.Arguments, parseBranch(i, current, tokens))
	}

	expectClose(i, current, tokens)

	return
}

func parseBranch(i *int, current *lexer.Token, tokens []lexer.Token) (node Node) {
	node = Node{
		Type:      ACTION_BRANCH,
		Arguments: make([]Node, 0),
		Token:     lexer.Token{},
	}

	if current.Type >= lexer.INT && current.Type <= lexer.ATOM {
		node.Arguments = append(node.Arguments, parseLiteral(*current))
		next(i, current, tokens)
	} else if current.Type == lexer.IDENTIFIER {
		node.Arguments = append(node.Arguments, Node{
			Type:      IDENTIFIER,
			Arguments: nil,
			Token:     *current,
		})
		next(i, current, tokens)
	} else if current.Type == lexer.OPEN_BRACE {
		next(i, current, tokens)
		node.Arguments = append(node.Arguments, doParse(i, current, tokens, false))
	}

	return
}

func parseFor(i *int, current *lexer.Token, tokens []lexer.Token) (node Node) {
	node = Node{
		Type:      ACTION_WHILE,
		Arguments: make([]Node, 0),
		Token:     lexer.Token{},
	}
	next(i, current, tokens)

	if current.Type == lexer.IDENTIFIER {
		node.Arguments = append(node.Arguments, Node{
			Type:      IDENTIFIER,
			Arguments: nil,
			Token:     *current,
		})
		next(i, current, tokens)
	} else if current.Type == lexer.BOOL {
		node.Arguments = append(node.Arguments, parseLiteral(*current))
		next(i, current, tokens)
	} else if current.Type == lexer.OPEN_BRACE {
		next(i, current, tokens)

		if current.Value == "range" {
			node.Type = ACTION_FOR
			next(i, current, tokens)

			if current.Type != lexer.IDENTIFIER {
				fail("identifier", *current)
			}

			node.Arguments = append(node.Arguments, Node{
				Type:      IDENTIFIER,
				Arguments: nil,
				Token:     *current,
			})
			next(i, current, tokens)

			if current.Type == lexer.IDENTIFIER {
				node.Arguments = append(node.Arguments, Node{
					Type:      IDENTIFIER,
					Arguments: nil,
					Token:     *current,
				})
				next(i, current, tokens)
			} else if current.Type == lexer.OPEN_BRACE {
				next(i, current, tokens)
				node.Arguments = append(node.Arguments, parseCall(i, current, tokens))
			} else {
				fail("an identifier or a call", *current)
			}

			expectClose(i, current, tokens)
		} else {
			node.Arguments = append(node.Arguments, parseCall(i, current, tokens))
		}
	}

	for current.Type != lexer.CLOSE_BRACE {
		if current.Type != lexer.OPEN_BRACE {
			fail("(", *current)
		}
		next(i, current, tokens)
		node.Arguments = append(node.Arguments, doParse(i, current, tokens, false))
	}

	expectClose(i, current, tokens)

	return
}

func parseDef(i *int, current *lexer.Token, tokens []lexer.Token, canMacro bool) (node Node) {
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
		} else if current.Type == lexer.IDENTIFIER_MACRO {
			if !canMacro {
				panic(fmt.Sprintf("Macros can only be defined in root scope! (line %d)", current.Line))
			}

			//TODO: Parse macro - macros come last (after the interpreter is done)
		} else if current.Type == lexer.IDENTIFIER_STRUCT {
			next(i, current, tokens)
			atoms := Node{
				Type:      STRUCT_DEF,
				Arguments: make([]Node, 0),
				Token:     lexer.Token{},
			}

			for current.Type != lexer.CLOSE_BRACE {
				if current.Type != lexer.ATOM {
					fail("an atom", *current)
				}
				next(i, current, tokens)
				atoms.Arguments = append(atoms.Arguments, parseLiteral(*current))
			}

			node.Arguments = append(node.Arguments, atoms)
		} else {
			node.Arguments = append(node.Arguments, parseCall(i, current, tokens))
		}
	}

	expectClose(i, current, tokens)

	return
}

func doParse(i *int, current *lexer.Token, tokens []lexer.Token, canMacro bool) Node {
	switch current.Type {
	case lexer.IDENTIFIER_DEF:
		return parseDef(i, current, tokens, canMacro)
	case lexer.IDENTIFIER_FOR:
		return parseFor(i, current, tokens)
	case lexer.IDENTIFIER_IF:
		return parseIf(i, current, tokens)
	case lexer.ATOM:
		return parseMapAccess(i, current, tokens, false)
	default:
		return parseCall(i, current, tokens)
	}
}

// Parse parse tokens and return an AST
func Parse(tokens []lexer.Token) []Node {
	nodes := make([]Node, 0)

	for i := 0; i < len(tokens); {
		current := tokens[i]

		if current.Type != lexer.OPEN_BRACE {
			fail("(", current)
		}

		next(&i, &current, tokens)

		nodes = append(nodes, doParse(&i, &current, tokens, true))
	}

	return nodes
}
