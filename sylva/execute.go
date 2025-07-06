package sylva

import (
	"fmt"
	"sylva/parser"

	"github.com/antlr4-go/antlr/v4"
)

func Execute(code, file string) error {
	input := antlr.NewInputStream(code)
	lexer := parser.NewSylvaLexer(input)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewSylvaParser(tokens)

	tree := p.Program()
	visitor := &SylvaVisitor{
		File: file,
	}
	visitor.Visit(tree)

	commands := visitor.Commands
	fmt.Println("Commands:")
	for _, command := range commands {
		fmt.Println(command.StringRepresentation())
	}
	fmt.Println()

	runtime := CreateSylvaRuntime()
	runtime.ProvideSourceCode(code)

	if err := runtime.ConvertToBytecode(commands); err != nil {
		fmt.Println("Error while converting to bytecode:", err)
		return err
	}
	fmt.Println("Bytecode:", runtime.Bytecode)

	if err := runtime.ExecuteUntilDone(); err != nil {
		fmt.Println("Error:", err)
		return err
	} else {
		fmt.Println("Result:", runtime.Registers[runtime.RegisterNames["$varres"]])
	}

	return nil
}
