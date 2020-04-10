package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	log "github.com/sirupsen/logrus"
)

type stdHook struct {
	callback interpreter.Node
}

const hookVal string = "hook"

func (hook *stdHook) Fire(entry *log.Entry) error {
	fields := make(map[string]interpreter.Node, 0)

	for k, v := range entry.Data {
		switch v.(type) {
		case float64:
			fields[k] = interpreter.FloatNode(v.(float64))
		case int:
			fields[k] = interpreter.IntNode(v.(int))
		case string:
			fields[k] = interpreter.StringNode(v.(string))
		case bool:
			fields[k] = interpreter.BoolNode(v.(bool))
		}
	}

	ctx := make(map[string]interpreter.Node, 0)
	ctx[hookVal] = interpreter.NodeMap(map[string]interpreter.Node{
		"level":   interpreter.StringNode(entry.Level.String()),
		"time":    interpreter.StringNode(entry.Time.String()),
		"message": interpreter.StringNode(entry.Message),
		"data":    interpreter.NodeMap(fields),
	})

	interpreter.DoVariableCall(parser.Node{
		Type: 0,
		Arguments: []parser.Node{
			{
				Type:      parser.IDENTIFIER,
				Arguments: nil,
				Token: lexer.Token{
					Value: hookVal,
					Type:  0,
					Line:  0,
				},
			},
		},
		Token: lexer.Token{},
	}, hook.callback, &ctx)

	return nil
}

func (hook *stdHook) Levels() []log.Level {
	return log.AllLevels
}
