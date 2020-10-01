package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"regexp"
)

// CALL concurrency functions
var CALL string = "--system-do-regex!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "compile":
		return doCompile(args)

	case "ismatch":
		return doIsMatch(args)

	case "matches":
		return doMatches(args)

	case "replace":
		return doReplace(args)
	}

	return interpreter.Nothing
}

func getReg(args []interpreter.Node) *regexp.Regexp {
	if args[1].NodeType != interpreter.NODETYPE_STRING {
		r, ok := interpreter.LoadObject(args[1].Value.(int))

		if !ok {
			return nil
		}

		return r.(*regexp.Regexp)
	}

	r, err := regexp.Compile(args[1].Value.(string))

	if err != nil {
		return nil
	}

	return r
}

func doCompile(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :compile")

	regex, err := regexp.Compile(args[1].Value.(string))

	if err != nil {
		return interpreter.OptionNode(interpreter.Nothing, true)
	}

	return interpreter.OptionNode(interpreter.IntNode(interpreter.StoreObject(regex)), false)
}

func doIsMatch(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureTypes(&args, 1, []interpreter.NodeType{interpreter.NODETYPE_STRING, interpreter.NODETYPE_INT}, CALL+" :ismatch")
	interpreter.EnsureType(&args, 2, interpreter.NODETYPE_STRING, CALL+" :ismatch")

	strToMatch := args[2].Value.(string)

	reg := getReg(args)
	if reg == nil {
		return interpreter.BoolNode(false)
	}

	return interpreter.BoolNode(reg.MatchString(strToMatch))
}

func doMatches(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureTypes(&args, 1, []interpreter.NodeType{interpreter.NODETYPE_STRING, interpreter.NODETYPE_INT}, CALL+" :matches")
	interpreter.EnsureType(&args, 2, interpreter.NODETYPE_STRING, CALL+" :matches")

	strToMatch := args[2].Value.(string)
	matches := make([]interpreter.Node, 0)

	reg := getReg(args)
	if reg == nil {
		return interpreter.OptionNode(interpreter.Nothing, true)
	}

	for _, s := range reg.FindStringSubmatch(strToMatch) {
		matches = append(matches, interpreter.StringNode(s))
	}
	return interpreter.OptionNode(interpreter.NodeList(matches), !reg.MatchString(strToMatch))
}

func doReplace(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureTypes(&args, 1, []interpreter.NodeType{interpreter.NODETYPE_STRING, interpreter.NODETYPE_INT}, CALL+" :replace")
	interpreter.EnsureType(&args, 2, interpreter.NODETYPE_STRING, CALL+" :replace")
	interpreter.EnsureType(&args, 3, interpreter.NODETYPE_STRING, CALL+" :replace")

	reg := getReg(args)
	if reg == nil {
		return interpreter.OptionNode(interpreter.StringNode(""), true)
	}

	return interpreter.StringNode(reg.ReplaceAllString(args[2].Value.(string), args[3].Value.(string)))
}
