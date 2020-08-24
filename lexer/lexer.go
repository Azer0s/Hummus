package lexer

import (
	"math"
	"unicode"
)

func set(value string, tokenType TokenType, buffer *[]rune, typePtr *TokenType) {
	*buffer = []rune(value)
	*typePtr = tokenType
}

func next(i *int, current *rune, code string) {
	*i += 1

	if *i >= len(code) {
		*current = math.MaxInt32
		return
	}

	*current = rune(code[*i])
}

// LexString tokenizes a string
func LexString(code string) []Token {
	tokens := make([]Token, 0)
	line := uint(1)

	for i := 0; i < len(code); i++ {
		buffer := make([]rune, 0)
		current := rune(code[i])
		tokenType := TokenType(uint8(255))

		if current == '(' {
			set("(", OPEN_BRACE, &buffer, &tokenType)
		} else if current == ')' {
			set(")", CLOSE_BRACE, &buffer, &tokenType)
		} else if unicode.IsDigit(current) {
			tokenType = INT

			for unicode.IsDigit(current) {
				buffer = append(buffer, current)
				next(&i, &current, code)
			}

			if current == '.' {
				tokenType = FLOAT
				buffer = append(buffer, current)
				bufferLength := len(buffer)
				next(&i, &current, code)

				for unicode.IsDigit(current) {
					buffer = append(buffer, current)
					next(&i, &current, code)
				}

				if len(buffer) == bufferLength {
					buffer = append(buffer, '0')
				}
			}

			i--
		} else if current == ':' {
			tokenType = ATOM
			next(&i, &current, code)

			for !unicode.IsSpace(current) && current != math.MaxInt32 && current != ')' && current != '(' {
				buffer = append(buffer, current)
				next(&i, &current, code)
			}

			i--
		} else if current == '"' {
			tokenType = STRING
			next(&i, &current, code)
			for current != '"' && (!unicode.IsSpace(current) || current == ' ') && current != math.MaxInt32 {
				if current == '\\' {
					next(&i, &current, code)
					switch current {
					case 'n':
						buffer = append(buffer, '\n')
						break
					case 't':
						buffer = append(buffer, '\t')
						break
					case 'r':
						buffer = append(buffer, '\r')
						break
					case 'v':
						buffer = append(buffer, '\v')
						break
					case '"':
						buffer = append(buffer, '"')
						break
					}
					next(&i, &current, code)
				} else {
					buffer = append(buffer, current)
					next(&i, &current, code)
				}
			}
		} else if current == ';' {
			for current != '\n' && current != math.MaxInt32 {
				next(&i, &current, code)
			}
			i--
		} else if !unicode.IsSpace(current) {
			tokenType = IDENTIFIER

			for !unicode.IsSpace(current) && current != math.MaxInt32 && current != ')' && current != '(' {
				buffer = append(buffer, current)
				next(&i, &current, code)
			}

			i--

			switch string(buffer) {
			case "def":
				tokenType = IDENTIFIER_DEF
				break
			case "fn":
				tokenType = IDENTIFIER_FN
				break
			case "if":
				tokenType = IDENTIFIER_IF
				break
			case "struct":
				tokenType = IDENTIFIER_STRUCT
				break
			case "macro":
				tokenType = IDENTIFIER_MACRO
				break
			case "for":
				tokenType = IDENTIFIER_FOR
				break
			case "true":
				tokenType = BOOL
				break
			case "false":
				tokenType = BOOL
				break
			}
		} else if current == '\n' {
			line++
		}

		if tokenType != TokenType(uint8(255)) {
			tokens = append(tokens, Token{
				Value: string(buffer),
				Type:  tokenType,
				Line:  line,
			})
		}
	}

	return tokens
}
