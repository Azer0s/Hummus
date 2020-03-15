package hummus

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"io/ioutil"
)

func RunFile(filename string) interpreter.Node {
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return interpreter.Run(parser.Parse(lexer.LexString(string(b))), make(map[string]interpreter.Node, 0))
}
