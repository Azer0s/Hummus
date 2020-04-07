package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"os"
	"os/exec"
	"strings"
)

// CALL string functions
var CALL string = "--system-do-os!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "exit":
		return exit(args)

	case "env":
		return env(args)

	case "env-all":
		return envAll()

	case "args":
		return getArgs()

	case "cmd-args":
		return cmdArgs(args)

	case "cmd":
		return cmd(args)

	default:
		panic("Unrecognized mode")
	}
}

func exit(args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_INT {
		panic(CALL + " :exit only accepts ints!")
	}

	os.Exit(args[1].Value.(int))

	return interpreter.Nothing
}

func env(args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_STRING {
		panic(CALL + " :env only accepts strings!")
	}

	return interpreter.StringNode(os.Getenv(args[1].Value.(string)))
}

func envAll() interpreter.Node {
	vals := make(map[string]interpreter.Node, 0)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)

		vals[pair[0]] = interpreter.StringNode(pair[1])
	}

	return interpreter.NodeMap(vals)
}

func getArgs() interpreter.Node {
	nodes := make([]interpreter.Node, 0)

	for _, arg := range os.Args {
		nodes = append(nodes, interpreter.StringNode(arg))
	}

	return interpreter.NodeList(nodes)
}

func cmdArgs(args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_STRING {
		panic(CALL + " :env expects a string as the first argument!")
	}

	if args[2].NodeType != interpreter.NODETYPE_LIST {
		panic(CALL + " :env expects a list as the second argument!")
	}

	cmdArgs := make([]string, 0)

	for _, value := range args[2].Value.(interpreter.ListNode).Values {
		cmdArgs = append(cmdArgs, interpreter.DumpNode(value))
	}

	b, err := exec.Command(args[1].Value.(string), cmdArgs...).CombinedOutput()

	out := string(b)

	return interpreter.OptionalNode(interpreter.StringNode(out), err != nil)
}

func cmd(args []interpreter.Node) interpreter.Node {
	if args[1].NodeType != interpreter.NODETYPE_STRING {
		panic(CALL + " :env expects a string as the first argument!")
	}

	b, err := exec.Command(args[1].Value.(string)).CombinedOutput()
	out := string(b)

	return interpreter.OptionalNode(interpreter.StringNode(out), err != nil)
}
