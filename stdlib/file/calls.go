package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"io/ioutil"
	"math"
	"os"
)

func main() {}

// CALL string functions
var CALL string = "--system-do-file!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

func convertDecimalToOctal(number int) int {
	octal := 0
	counter := 1
	remainder := 0
	for number != 0 {
		remainder = number % 8
		number = number / 8
		octal += remainder * counter
		counter *= 10
	}
	return octal
}

func convertOctaToDecimal(number int) int {
	decimal := 0
	counter := 0.0
	remainder := 0

	for number != 0 {
		remainder = number % 10
		decimal += remainder * int(math.Pow(8.0, counter))
		number = number / 10
		counter++
	}
	return decimal
}

func stats(f os.FileInfo, err error) interpreter.Node {
	if err != nil {
		return interpreter.NodeMap(map[string]interpreter.Node{
			"value": interpreter.StringNode(err.Error()),
			"error": interpreter.BoolNode(true),
		})
	}

	return interpreter.NodeMap(map[string]interpreter.Node{
		"value": interpreter.NodeMap(map[string]interpreter.Node{
			"name":        interpreter.StringNode(f.Name()),
			"size":        interpreter.IntNode(int(f.Size())),
			"dir":         interpreter.BoolNode(f.IsDir()),
			"modified":    interpreter.IntNode(int(f.ModTime().Unix())),
			"mode":        interpreter.IntNode(convertDecimalToOctal(int(f.Mode()))),
			"mode-string": interpreter.StringNode(f.Mode().String()),
		}),
		"error": interpreter.BoolNode(false),
	})

}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "read-string":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :read-string")
		b, err := ioutil.ReadFile(args[1].Value.(string))

		if err != nil {
			return interpreter.NodeMap(map[string]interpreter.Node{
				"value": interpreter.StringNode(err.Error()),
				"error": interpreter.BoolNode(true),
			})
		}

		return interpreter.NodeMap(map[string]interpreter.Node{
			"value": interpreter.StringNode(string(b)),
			"error": interpreter.BoolNode(false),
		})

	case "read-bytes":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :read-bytes")
		b, err := ioutil.ReadFile(args[1].Value.(string))

		if err != nil {
			return interpreter.NodeMap(map[string]interpreter.Node{
				"value": interpreter.StringNode(err.Error()),
				"error": interpreter.BoolNode(true),
			})
		}

		return interpreter.NodeMap(map[string]interpreter.Node{
			"value": interpreter.IntNode(interpreter.StoreObject(b)),
			"error": interpreter.BoolNode(false),
		})

	case "exists":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :exists")
		_, err := os.Open(args[1].Value.(string))
		return interpreter.BoolNode(!os.IsNotExist(err))

	case "stats":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :stats")
		f, err := os.Stat(args[1].Value.(string))
		return stats(f, err)

	case "create":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :create")
		interpreter.EnsureType(&args, 2, interpreter.NODETYPE_INT, CALL+" :create")

		f, err := os.OpenFile(args[1].Value.(string), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(convertOctaToDecimal(args[2].Value.(int))))

		if err != nil {
			return interpreter.NodeMap(map[string]interpreter.Node{
				"value": interpreter.StringNode(err.Error()),
				"error": interpreter.BoolNode(true),
			})
		}

		return stats(f.Stat())

	case "delete":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :delete")
		err := os.Remove(args[1].Value.(string))
		return interpreter.BoolNode(err == nil)

	case "write-string":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :write-string")
		interpreter.EnsureType(&args, 2, interpreter.NODETYPE_STRING, CALL+" :write-string")
		interpreter.EnsureType(&args, 3, interpreter.NODETYPE_INT, CALL+" :write-string")

		err := ioutil.WriteFile(args[1].Value.(string), []byte(args[2].Value.(string)), os.FileMode(convertOctaToDecimal(args[3].Value.(int))))

		return interpreter.BoolNode(err == nil)

	case "write-bytes":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :write-bytes")
		interpreter.EnsureType(&args, 2, interpreter.NODETYPE_INT, CALL+" :write-bytes")
		interpreter.EnsureType(&args, 3, interpreter.NODETYPE_INT, CALL+" :write-bytes")

		b, ok := interpreter.LoadObject(args[2].Value.(int))

		if !ok {
			panic("Byte array doesn't exist or was already freed!")
		}

		err := ioutil.WriteFile(args[1].Value.(string), b.([]byte), os.FileMode(convertOctaToDecimal(args[3].Value.(int))))

		return interpreter.BoolNode(err == nil)
	default:
		panic("Unrecognized mode")
	}
}
