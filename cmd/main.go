package main

import (
	"fmt"
	"sylva/sylva"
	// "github.com/antlr4-go/antlr/v4"
	// "sylva/parser"
	// "sylva/sylva"
)

func main() {
	// data := "1 + 2 * 3"

	// input := antlr.NewInputStream(data)
	// lexer := parser.NewSylvaLexer(input)
	// tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	// p := parser.NewSylvaParser(tokens)

	// tree := p.Expr()
	// visitor := &sylva.SylvaVisitor{BaseSylvaVisitor: parser.BaseSylvaVisitor{}}
	// result := visitor.Visit(tree)

	runtime := sylva.CreateSylvaRuntime()
	commands := []sylva.Command{
		&sylva.LoadCommand{Register: "$a", Value: "\"Hello, \""},
		&sylva.LoadCommand{Register: "$b", Value: "\"World!\""},
		&sylva.BinOpCommand{Operation: "concat", Register: "$res", A: "$a", B: "$b"},
		&sylva.FreeCommand{Register: "$a"},
		&sylva.FreeCommand{Register: "$b"},
	}

	if err := runtime.ConvertToBytecode(commands); err != nil {
		fmt.Println("Error while converting to byteocde: ", err)
	}
	fmt.Println("Bytecode:", runtime.Bytecode)

	if err := runtime.ExecuteUntilDone(); err != nil {
		fmt.Println("Error:", err, "ip", runtime.IP)
	} else {
		fmt.Println("Result:", runtime.Registers[runtime.RegisterNames["$res"]])
	}
}
