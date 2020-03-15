package main

import (
	"fmt"
	"github.com/Azer0s/Hummus/hummus"
	"github.com/Azer0s/Hummus/lexer"
)

func main() {
	lexer.LexString("22")
	lexer.LexString("3.1415")
	lexer.LexString("3.")
	lexer.LexString(":func")
	lexer.LexString("\"Hello world\"")
	lexer.LexString("((fn x (* x x)) 4)")

	fmt.Println(hummus.RunFile("examples/def_mult.hummus").Value)
}
