package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"math/rand"
)

// CALL random functions
var CALL string = "--system-do-random!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "intn":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_INT, CALL+":int")
		return interpreter.IntNode(rand.Intn(args[1].Value.(int)))
	case "int":
		return interpreter.IntNode(rand.Int())
	case "float":
		return interpreter.FloatNode(rand.Float64())
	case "string":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_INT, CALL+":string")
		return interpreter.StringNode(randomString(args[1].Value.(int)))
	case "stringc":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_INT, CALL+":stringc")
		interpreter.EnsureType(&args, 2, interpreter.NODETYPE_STRING, CALL+":stringc")
		return interpreter.StringNode(randomStringWithCharset(args[1].Value.(int), args[2].Value.(string)))
	}

	return interpreter.Nothing
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return randomStringWithCharset(length, charset)
}
