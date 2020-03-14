package main

import (
	"encoding/json"
	"fmt"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
)

func main() {
	lexer.LexString("22")
	lexer.LexString("3.1415")
	lexer.LexString("3.")
	lexer.LexString(":func")
	lexer.LexString("\"Hello world\"")

	b, _ := json.Marshal(lexer.LexString("((fn x (* x x)) 4)"))
	fmt.Println(string(b))

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
