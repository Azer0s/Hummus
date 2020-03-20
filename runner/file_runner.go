package runner

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"io/ioutil"
)

// RunFile runs a file by filename
func RunFile(filename string) interpreter.Node {
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	vars := make(map[string]interpreter.Node, 0)
	return interpreter.Run(parser.Parse(lexer.LexString(string(b))), &vars)
}
