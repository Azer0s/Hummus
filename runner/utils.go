package runner

import "github.com/Azer0s/Hummus/interpreter"

func setupVars(vars *map[string]interpreter.Node, dir string) {
	(*vars)[interpreter.EXEC_FILE] = interpreter.StringNode(dir)
	(*vars)[interpreter.SELF] = interpreter.IntNode(0)
}
