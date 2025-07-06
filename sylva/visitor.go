package sylva

import (
	"fmt"
	"sylva/parser"

	"github.com/antlr4-go/antlr/v4"
)

// visit calls return the register in which they stored their result
type SylvaVisitor struct {
	parser.BaseSylvaVisitor
	File          string
	Commands      []Command
	registerTopID int
}

func (v *SylvaVisitor) GetRegister() string {
	v.registerTopID++
	return fmt.Sprintf("$i%v", v.registerTopID)
}

func (v *SylvaVisitor) Visit(tree antlr.ParseTree) any {
	return tree.Accept(v)
}

func (v *SylvaVisitor) VisitMulExpr(ctx *parser.MulExprContext) any {
	left := v.Visit(ctx.Expr(0)).(string)
	right := v.Visit(ctx.Expr(1)).(string)
	op := ctx.GetOp().GetText()
	var opStr string
	res := v.GetRegister()
	switch op {
	case "*":
		opStr = "mul"
	case "/":
		opStr = "div"
	case "%":
		opStr = "mod"
	default:
		return fmt.Errorf("there is no mul operand %v", op)
	}
	v.Commands = append(v.Commands, &BinOpCommand{
		Operation: opStr,
		Register:  res,
		A:         left,
		B:         right,
	})
	v.Commands = append(v.Commands, &FreeCommand{Register: left})
	v.Commands = append(v.Commands, &FreeCommand{Register: right})
	return res
}

func (v *SylvaVisitor) GetDebugData(ctx antlr.ParserRuleContext) CommandDebugData {
	start := ctx.GetStart()
	stop := ctx.GetStop()

	startLine := start.GetLine()
	startColumn := start.GetColumn()
	endLine := stop.GetLine()
	endColumn := stop.GetColumn() + len(stop.GetText())

	startIdx := start.GetStart()
	endIdx := stop.GetStop()

	debug := CommandDebugData{
		FileLocation: v.File,
		Start: ProgramLocation{
			Line:   startLine,
			Column: startColumn,
			Idx:    startIdx,
		},
		End: ProgramLocation{
			Line:   endLine,
			Column: endColumn,
			Idx:    endIdx,
		},
		FunctionDeclaration: struct {
			FunctionName       string
			FunctionLineNumber int
		}{
			FunctionName:       "<global>",
			FunctionLineNumber: 0,
		},
	}

	return debug
}

func (v *SylvaVisitor) VisitAddExpr(ctx *parser.AddExprContext) any {
	left := v.Visit(ctx.Expr(0)).(string)
	right := v.Visit(ctx.Expr(1)).(string)
	op := ctx.GetOp().GetText()
	var opStr string
	res := v.GetRegister()
	switch op {
	case "+":
		opStr = "add"
	case "-":
		opStr = "sub"
	default:
		return fmt.Errorf("there is no add operand %v", op)
	}
	debug := v.GetDebugData(ctx)
	v.Commands = append(v.Commands, &BinOpCommand{
		Operation: opStr,
		Register:  res,
		A:         left,
		B:         right,
		DebugData: debug,
	})
	v.Commands = append(v.Commands, &FreeCommand{Register: left, DebugData: debug})
	v.Commands = append(v.Commands, &FreeCommand{Register: right, DebugData: debug})
	return res
}

func (v *SylvaVisitor) VisitUnaryOp(ctx *parser.UnaryOpContext) any {
	debug := v.GetDebugData(ctx)
	// why is opToken nil???
	opToken := ctx.GetOp()
	op := opToken.GetText()
	reg := v.GetRegister()
	x := v.Visit(ctx.Expr()).(string)
	switch op {
	case "+":
		v.Commands = append(v.Commands, &LoadCommand{Register: reg, Value: x, DebugData: debug})
	case "-":
		v.Commands = append(v.Commands, &UmnCommand{Register: reg, Value: x, DebugData: debug})
	}
	v.Commands = append(v.Commands, &FreeCommand{Register: x, DebugData: debug})
	return reg
}

func (v *SylvaVisitor) VisitIntValue(ctx *parser.IntValueContext) any {
	debug := v.GetDebugData(ctx)
	str := ctx.INT().GetText()
	reg := v.GetRegister()
	v.Commands = append(v.Commands, &LoadCommand{Register: reg, Value: str, DebugData: debug})
	return reg
}

func (v *SylvaVisitor) VisitFloatValue(ctx *parser.FloatValueContext) any {
	debug := v.GetDebugData(ctx)
	str := ctx.FLOAT().GetText()
	reg := v.GetRegister()
	v.Commands = append(v.Commands, &LoadCommand{Register: reg, Value: str, DebugData: debug})
	return reg
}

func (v *SylvaVisitor) VisitStringValue(ctx *parser.StringValueContext) any {
	debug := v.GetDebugData(ctx)
	str := ctx.STRING().GetText()
	str = str[1 : len(str)-1]
	str = fmt.Sprintf("\"%v\"", str)
	// TODO: multiline strings
	reg := v.GetRegister()
	v.Commands = append(v.Commands, &LoadCommand{Register: reg, Value: str, DebugData: debug})
	return reg
}

func (v *SylvaVisitor) VisitBoolValue(ctx *parser.BoolValueContext) any {
	debug := v.GetDebugData(ctx)
	str := ctx.BOOL().GetText()
	reg := v.GetRegister()
	v.Commands = append(v.Commands, &LoadCommand{Register: reg, Value: str, DebugData: debug})
	return reg
}

func (v *SylvaVisitor) VisitParensValue(ctx *parser.ParensValueContext) any {
	return v.Visit(ctx.Expr())
}

func (v *SylvaVisitor) VisitValueExpr(ctx *parser.ValueExprContext) any {
	return v.Visit(ctx.Value())
}

func (v *SylvaVisitor) VisitConcatExpr(ctx *parser.ConcatExprContext) any {
	debug := v.GetDebugData(ctx)
	left := v.Visit(ctx.Expr(0)).(string)
	right := v.Visit(ctx.Expr(1)).(string)
	reg := v.GetRegister()
	v.Commands = append(v.Commands, &BinOpCommand{Operation: "concat", Register: reg, A: left, B: right, DebugData: debug})
	v.Commands = append(v.Commands, &FreeCommand{Register: left, DebugData: debug})
	v.Commands = append(v.Commands, &FreeCommand{Register: right, DebugData: debug})
	return reg
}

func (v *SylvaVisitor) VisitExprStatement(ctx *parser.ExprStatementContext) any {
	debug := v.GetDebugData(ctx)
	reg := v.Visit(ctx.Expr()).(string)
	v.Commands = append(v.Commands, &FreeCommand{Register: reg, DebugData: debug})
	return "<NULL>"
}

func (v *SylvaVisitor) VisitSetStatement(ctx *parser.SetStatementContext) any {
	debug := v.GetDebugData(ctx)
	reg := v.Visit(ctx.Expr()).(string)
	varReg := fmt.Sprintf("$var%v", ctx.ID().GetText())
	v.Commands = append(v.Commands, &LoadCommand{Register: varReg, Value: reg, DebugData: debug})
	v.Commands = append(v.Commands, &FreeCommand{Register: reg, DebugData: debug})
	return "<NULL>"
}

func (v *SylvaVisitor) VisitVarAccessValue(ctx *parser.VarAccessValueContext) any {
	debug := v.GetDebugData(ctx)
	str := ctx.ID().GetText()
	reg := v.GetRegister()
	varReg := fmt.Sprintf("$var%v", str)
	v.Commands = append(v.Commands, &LoadCommand{Register: reg, Value: varReg, DebugData: debug})
	return reg
}

func (v *SylvaVisitor) VisitFullProgram(ctx *parser.FullProgramContext) any {
	for _, statement := range ctx.AllStatement() {
		v.Visit(statement)
	}
	return "<NULL>"
}

func (v *SylvaVisitor) VisitListValue(ctx *parser.ListValueContext) any {
	reg := v.GetRegister()
	debug := v.GetDebugData(ctx)
	v.Commands = append(v.Commands, &ListCommand{
		Register:  reg,
		DebugData: debug,
	})
	return reg
}
