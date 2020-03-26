package runner

import "github.com/Azer0s/Hummus/interpreter"

func setupVars(vars *map[string]interpreter.Node, dir string) {
	(*vars)[interpreter.EXEC_FILE] = interpreter.Node{
		Value:    dir,
		NodeType: interpreter.NODETYPE_STRING,
	}
	(*vars)[interpreter.SELF] = interpreter.Node{
		Value:    0,
		NodeType: interpreter.NODETYPE_INT,
	}
}
