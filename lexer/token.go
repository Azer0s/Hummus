package lexer

// TokenType lexeme type
type TokenType uint8

// Token a lexeme
type Token struct {
	Value string    `json:"value"`
	Type  TokenType `json:"type"`
	Line  uint      `json:"line"`
}

const (
	// OPEN_BRACE (
	OPEN_BRACE TokenType = 1
	// CLOSE_BRACE )
	CLOSE_BRACE TokenType = 2
	// IDENTIFIER an identifier
	IDENTIFIER TokenType = 3
	// INT an integer
	INT TokenType = 4
	// FLOAT a floating point literal
	FLOAT TokenType = 5
	// STRING a string
	STRING TokenType = 6
	// BOOL a boolean
	BOOL TokenType = 7
	// ATOM an atom
	ATOM TokenType = 8
	// IDENTIFIER_DEF "def"
	IDENTIFIER_DEF TokenType = 9
	// IDENTIFIER_FN "fn"
	IDENTIFIER_FN TokenType = 10
	// IDENTIFIER_STRUCT "struct"
	IDENTIFIER_STRUCT TokenType = 11
	// IDENTIFIER_FOR "for"
	IDENTIFIER_FOR TokenType = 12
	// IDENTIFIER_IF "if"
	IDENTIFIER_IF TokenType = 13
)
