package runner

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"io/ioutil"
	"os"
	"path"
)

// RunFile runs a file by filename
func RunFile(filename string) interpreter.Node {
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	execFile := path.Join(p, filename)

	vars := make(map[string]interpreter.Node, 0)

	vars[interpreter.EXEC_FILE] = interpreter.Node{
		Value:    execFile,
		NodeType: interpreter.NODETYPE_STRING,
	}

	return interpreter.Run(parser.Parse(lexer.LexString(string(b))), &vars)
}
