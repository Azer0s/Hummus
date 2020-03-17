package main

import (
	"fmt"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/runner"
)

func main() {
	lexer.LexString("22")
	lexer.LexString("3.1415")
	lexer.LexString("3.")
	lexer.LexString(":func")
	lexer.LexString("\"Hello world\"")
	lexer.LexString("((fn x (* x x)) 4)")

	fmt.Println(runner.RunFile("examples/maps.hummus").Value)
}
