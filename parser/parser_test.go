package parser_test

import (
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"testing"
)

func TestParseCall(t *testing.T) {
	tokens := lexer.LexString("(out \"Hello\")")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_CALL {
		t.FailNow()
	}
}

func TestParseNestedCall(t *testing.T) {
	tokens := lexer.LexString("(out (get-result :foo))")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_CALL || node[0].Arguments[0].Type != parser.ACTION_CALL {
		t.FailNow()
	}
}

func TestParseNestedCallWithOtherArgs(t *testing.T) {
	tokens := lexer.LexString("(out (get-result :foo) 12.4)")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_CALL ||
		node[0].Arguments[0].Type != parser.ACTION_CALL ||
		node[0].Arguments[1].Type != parser.LITERAL_FLOAT {
		t.FailNow()
	}
}

func TestParseDef(t *testing.T) {
	tokens := lexer.LexString("(def a \"Hello\")")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_DEF ||
		node[0].Arguments[0].Type != parser.IDENTIFIER ||
		node[0].Arguments[1].Type != parser.LITERAL_STRING {
		t.FailNow()
	}
}

func TestParseDefFnLiteral(t *testing.T) {
	tokens := lexer.LexString("(def a (fn x (* x x)))")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_DEF ||
		node[0].Arguments[0].Type != parser.IDENTIFIER ||
		node[0].Arguments[1].Type != parser.LITERAL_FN {
		t.FailNow()
	}
}

func TestParseDefIdentity(t *testing.T) {
	tokens := lexer.LexString("(def a (identity b))")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_DEF ||
		node[0].Arguments[0].Type != parser.IDENTIFIER ||
		node[0].Arguments[1].Type != parser.ACTION_CALL {
		t.FailNow()
	}
}

func TestParseUseStdlib(t *testing.T) {
	tokens := lexer.LexString("(use :<base>)")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_CALL ||
		node[0].Token.Value != "use" ||
		node[0].Arguments[0].Type != parser.LITERAL_ATOM {
		t.FailNow()
	}
}

func TestParseUseFile(t *testing.T) {
	tokens := lexer.LexString("(use :./examples/supervisor)")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_CALL ||
		node[0].Token.Value != "use" ||
		node[0].Arguments[0].Type != parser.LITERAL_ATOM {
		t.FailNow()
	}
}

func TestParseUseNative(t *testing.T) {
	tokens := lexer.LexString("(use :./regex/calls.so :native)")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_CALL ||
		node[0].Token.Value != "use" ||
		node[0].Arguments[0].Type != parser.LITERAL_ATOM ||
		node[0].Arguments[1].Type != parser.LITERAL_ATOM {
		t.FailNow()
	}
}

func TestParseCallLiteralFn(t *testing.T) {
	tokens := lexer.LexString("((fn x (* x x)) 4)")
	node := parser.Parse(tokens)

	if node[0].Token.Type != lexer.ANONYMOUS_FN ||
		node[0].Arguments[0].Type != parser.LITERAL_FN ||
		node[0].Arguments[1].Type != parser.LITERAL_INT {
		t.FailNow()
	}
}

func TestParseWhile(t *testing.T) {
	tokens := lexer.LexString("(for true (out \"Hello\"))")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_WHILE ||
		node[0].Arguments[0].Type != parser.LITERAL_BOOL ||
		node[0].Arguments[1].Type != parser.ACTION_CALL {
		t.FailNow()
	}
}

func TestParseWhileCall(t *testing.T) {
	tokens := lexer.LexString("(for (is_valid) (out \"Hello\"))")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_WHILE ||
		node[0].Arguments[0].Type != parser.ACTION_CALL ||
		node[0].Arguments[1].Type != parser.ACTION_CALL {
		t.FailNow()
	}
}

func TestParseFor(t *testing.T) {
	tokens := lexer.LexString("(for (range a b) (out a))")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_FOR ||
		node[0].Arguments[0].Type != parser.IDENTIFIER ||
		node[0].Arguments[1].Type != parser.IDENTIFIER ||
		node[0].Arguments[2].Type != parser.ACTION_CALL {
		t.FailNow()
	}
}

func TestParseForCall(t *testing.T) {
	tokens := lexer.LexString("(for (range a (.. 0 10)) (out a))")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_FOR ||
		node[0].Arguments[0].Type != parser.IDENTIFIER ||
		node[0].Arguments[1].Type != parser.ACTION_CALL ||
		node[0].Arguments[2].Type != parser.ACTION_CALL {
		t.FailNow()
	}
}

func TestParseMap(t *testing.T) {
	tokens := lexer.LexString("({} (:foo \"Bar\") (:hello \"World\"))")
	node := parser.Parse(tokens)

	if node[0].Type != parser.ACTION_MAP ||
		node[0].Arguments[0].Type != parser.ACTION_MAP_ACCESS ||
		node[0].Arguments[1].Type != parser.ACTION_MAP_ACCESS {
		t.FailNow()
	}
}
