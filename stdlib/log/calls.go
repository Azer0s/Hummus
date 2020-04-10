package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	log "github.com/sirupsen/logrus"
)

// CALL string functions
var CALL string = "--system-do-log!"

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	// noinit
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "create-hook":
		return doCreateHook(args)
	case "register-hook":
		return doRegisterHook(args)
	case "json":
		return doJson(args)
	case "set-level":
		return doSetLevel(args)
	case "log":
		return doLog(args, variables)
	case "log-props":
		return doLogProps(args, variables)
	default:
		panic("Unrecognized mode")
	}
}

func doCreateHook(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_FN, CALL+" :create-hook")
	if len(args[1].Value.(interpreter.FnLiteral).Parameters) != 1 {
		panic(CALL + " :create-hook expects a function with 1 parameter!")
	}

	return interpreter.IntNode(interpreter.StoreObject(&stdHook{callback: args[1]}))
}

func doRegisterHook(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_INT, CALL+" :register-hook")
	val, ok := interpreter.LoadObject(args[1].Value.(int))

	if !ok {
		return interpreter.BoolNode(false)
	}

	if hook, ok := val.(log.Hook); ok {
		log.AddHook(hook)
		return interpreter.BoolNode(true)
	}

	return interpreter.BoolNode(false)
}

func doJson(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_BOOL, CALL+" :json")

	if args[1].Value.(bool) {
		log.SetFormatter(&log.JSONFormatter{})
		return interpreter.Nothing
	}

	log.SetFormatter(&log.TextFormatter{})
	return interpreter.Nothing
}

func doSetLevel(args []interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_ATOM, CALL+" :set-level")
	lvl, err := log.ParseLevel(args[1].Value.(string))

	if err != nil {
		return interpreter.BoolNode(false)
	}

	log.SetLevel(lvl)

	return interpreter.BoolNode(true)
}

func doLogProps(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_ATOM, CALL+" :log-props")
	interpreter.EnsureType(&args, 2, interpreter.NODETYPE_MAP, CALL+" :log-props")

	lvl, err := log.ParseLevel(args[1].Value.(string))

	if err != nil {
		return interpreter.BoolNode(false)
	}

	entry := log.WithFields(log.Fields{
		"pid": (*variables)[interpreter.SELF].Value,
	})

	for k, v := range args[2].Value.(interpreter.MapNode).Values {
		if v.NodeType >= interpreter.NODETYPE_FN {
			entry = entry.WithField(k, interpreter.DumpNode(v))
			continue
		}

		entry = entry.WithField(k, v.Value)
	}

	val := make([]interface{}, 0)

	if args[3].NodeType == interpreter.NODETYPE_LIST {
		for _, value := range args[3].Value.(interpreter.ListNode).Values {
			val = append(val, interpreter.DumpNode(value))
		}
	} else {
		val = append(val, interpreter.DumpNode(args[3]))
	}

	return execLog(entry, lvl, val)
}

func doLog(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	interpreter.EnsureType(&args, 1, interpreter.NODETYPE_ATOM, CALL+" :log")

	lvl, err := log.ParseLevel(args[1].Value.(string))

	if err != nil {
		return interpreter.BoolNode(false)
	}

	entry := log.WithFields(log.Fields{
		"pid": (*variables)[interpreter.SELF].Value,
	})

	val := make([]interface{}, 0)

	if args[2].NodeType == interpreter.NODETYPE_LIST {
		for _, value := range args[2].Value.(interpreter.ListNode).Values {
			val = append(val, interpreter.DumpNode(value))
		}
	} else {
		val = append(val, interpreter.DumpNode(args[2]))
	}

	return execLog(entry, lvl, val)
}

func execLog(entry *log.Entry, lvl log.Level, val []interface{}) interpreter.Node {
	switch lvl {
	case log.PanicLevel:
		entry.Panic(val...)
	case log.FatalLevel:
		entry.Fatal(val...)
	case log.ErrorLevel:
		entry.Error(val...)
	case log.WarnLevel:
		entry.Warn(val...)
	case log.InfoLevel:
		entry.Info(val...)
	case log.DebugLevel:
		entry.Debug(val...)
	case log.TraceLevel:
		entry.Trace(val...)
	}

	return interpreter.Nothing
}
