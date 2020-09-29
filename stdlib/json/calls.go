package main

import (
	"encoding/json"
	"fmt"
	"github.com/Azer0s/Hummus/interpreter"
	"reflect"
)

func main() {}

// CALL HTTP server functions
var CALL string = "--system-do-json!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

func doMarshal(obj interpreter.Node) interface{} {
	switch obj.NodeType {
	case interpreter.NODETYPE_INT:
		return obj.Value.(int)
	case interpreter.NODETYPE_FLOAT:
		return obj.Value.(float64)
	case interpreter.NODETYPE_STRING:
		return obj.Value.(string)
	case interpreter.NODETYPE_BOOL:
		return obj.Value.(bool)
	case interpreter.NODETYPE_ATOM:
		return obj.Value.(string)
	case interpreter.NODETYPE_FN:
		return "[fn]"
	case interpreter.NODETYPE_LIST:
		l := make([]interface{}, 0)
		for _, value := range obj.Value.(interpreter.ListNode).Values {
			l = append(l, doMarshal(value))
		}
		return l
	case interpreter.NODETYPE_MAP:
		m := make(map[string]interface{})
		for s, node := range obj.Value.(interpreter.MapNode).Values {
			m[s] = doMarshal(node)
		}
		return m
	case interpreter.NODETYPE_STRUCT:
		return "[struct]"
	case interpreter.NODETYPE_MACRO:
		return "[macro]"
	}

	return interpreter.Node{}
}

func doTopLevelMarshal(obj interpreter.Node) interpreter.Node {
	if obj.NodeType != interpreter.NODETYPE_LIST && obj.NodeType != interpreter.NODETYPE_MAP {
		panic(CALL + ":marshal expects a list or a map as the 1st argument!")
	}

	str, _ := json.Marshal(doMarshal(obj))

	return interpreter.StringNode(string(str))
}

func doUnmarshal(i interface{}) interpreter.Node {
	switch i.(type) {
	case []interface{}:
		a := interpreter.ListNode{Values: []interpreter.Node{}}

		for _, val := range i.([]interface{}) {
			a.Values = append(a.Values, doUnmarshal(val))
		}

		return interpreter.Node{
			Value:    a,
			NodeType: interpreter.NODETYPE_LIST,
		}

	case map[string]interface{}:
		m := interpreter.MapNode{Values: map[string]interpreter.Node{}}

		for key, val := range i.(map[string]interface{}) {
			m.Values[key] = doUnmarshal(val)
		}

		return interpreter.Node{
			Value:    m,
			NodeType: interpreter.NODETYPE_MAP,
		}

	case int:
		return interpreter.IntNode(i.(int))

	case string:
		return interpreter.StringNode(i.(string))

	case bool:
		return interpreter.BoolNode(i.(bool))

	case float64:
		return interpreter.FloatNode(i.(float64))

	default:
		panic(fmt.Sprintf("Unrecognized type %s!", reflect.TypeOf(i).String()))
	}
}

func doTopLevelUnmarshal(arg interpreter.Node) interpreter.Node {
	interpreter.EnsureSingleType(&arg, 1, interpreter.NODETYPE_STRING, CALL+" :unmarshal")

	var v interface{}
	err := json.Unmarshal([]byte(arg.Value.(string)), &v)

	if err != nil {
		fmt.Println(err.Error())
		return interpreter.NodeMap(map[string]interpreter.Node{
			"value": interpreter.Nothing,
			"error": interpreter.BoolNode(true),
		})
	}

	return interpreter.NodeMap(map[string]interpreter.Node{
		"value": doUnmarshal(v),
		"error": interpreter.BoolNode(false),
	})
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "valid":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :valid")
		return interpreter.BoolNode(json.Valid([]byte(args[1].Value.(string))))

	case "marshal":
		return doTopLevelMarshal(args[1])

	case "unmarshal":
		return doTopLevelUnmarshal(args[1])

	default:
		panic("Unrecognized mode")
	}
}
