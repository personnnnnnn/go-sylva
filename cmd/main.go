package main

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sylva/parser"
	"sylva/sylva"
)

func main() {
	data := "res = 'Hello, '\n+ 'World!'"

	input := antlr.NewInputStream(data)
	lexer := parser.NewSylvaLexer(input)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewSylvaParser(tokens)

	tree := p.Program()
	visitor := &sylva.SylvaVisitor{
		File: "<stdin>",
	}
	visitor.Visit(tree)

	commands := visitor.Commands
	fmt.Println("Commands:")
	for _, command := range commands {
		fmt.Println(command.StringRepresentation())
	}
	fmt.Println()

	runtime := sylva.CreateSylvaRuntime()
	runtime.ProvideSourceCode(data)

	if err := runtime.ConvertToBytecode(commands); err != nil {
		fmt.Println("Error while converting to bytecode: ", err)
		return
	}
	fmt.Println("Bytecode:", runtime.Bytecode)

	if err := runtime.ExecuteUntilDone(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", runtime.Registers[runtime.RegisterNames["$varres"]])
	}
}
