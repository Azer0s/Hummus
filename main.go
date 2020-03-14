package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
)

func main() {
	lexer.LexString("22")
	lexer.LexString("3.1415")
	lexer.LexString("3.")
	lexer.LexString(":func")
	lexer.LexString("\"Hello world\"")
	lexer.LexString("((fn x (* x x)) 4)")

	interpreter.Run(parser.Parse(lexer.LexString(`
(def a 20)
(def b true)
(def c "Hello world")
(def d 3.1414)
(def e :ok)
`)), make(map[string]interpreter.Node, 0))

	parser.Parse(lexer.LexString(`
(def times (fn x y
  (* x y)))
`))

	parser.Parse(lexer.LexString(`
(def Animal (struct ; Declare struct
  :name
  :age
  :race
))

(def tom (Animal ; use struct
  "Tom"
  1
  :cat
))

(out (:name tom) " is a " (' (:race tom)))
`))
}
