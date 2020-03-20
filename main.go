package main

import (
	"github.com/Azer0s/Hummus/runner"
)

func main() {
	//lexer.LexString("22")
	//lexer.LexString("3.1415")
	//lexer.LexString("3.")
	//lexer.LexString(":func")
	//lexer.LexString("\"Hello world\"")
	//lexer.LexString("((fn x (* x x)) 4)")

	runner.RunFile("examples/fib.hummus")
}
