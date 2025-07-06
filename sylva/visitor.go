package sylva

import (
	"fmt"
	"strconv"
	"sylva/parser"

	"github.com/antlr4-go/antlr/v4"
)

// this will be replaced by a command generator
type SylvaVisitor struct {
	parser.BaseSylvaVisitor
}

func (v *SylvaVisitor) Visit(tree antlr.ParseTree) any {
	return tree.Accept(v)
}

func (v *SylvaVisitor) VisitMulExpr(ctx *parser.MulExprContext) any {
	left := v.Visit(ctx.Expr(0)).(Value)
	right := v.Visit(ctx.Expr(1)).(Value)
	op := ctx.GetOp().GetText()
	switch op {
	case "*":
		if res, err := Mul(left, right); err != nil {
			return err
		} else {
			return res
		}
	case "/":
		if res, err := Div(left, right); err != nil {
			return err
		} else {
			return res
		}
	default:
		return fmt.Errorf("there is no mul operand %v", op)
	}
}
func (v *SylvaVisitor) VisitAddExpr(ctx *parser.AddExprContext) any {
	left := v.Visit(ctx.Expr(0)).(Value)
	right := v.Visit(ctx.Expr(1)).(Value)
	op := ctx.GetOp().GetText()
	switch op {
	case "+":
		if res, err := Add(left, right); err != nil {
			return err
		} else {
			return res
		}
	case "-":
		if res, err := Sub(left, right); err != nil {
			return err
		} else {
			return res
		}
	default:
		return fmt.Errorf("there is no add operand %v", op)
	}
}
func (v *SylvaVisitor) VisitUnaryOp(ctx *parser.UnaryOpContext) any {
	val := v.Visit(ctx.Expr()).(Value)
	op := ctx.GetOp().GetText()
	switch op {
	case "+":
		return val
	case "-":
		if res, err := Umn(val); err != nil {
			return err
		} else {
			return res
		}
	default:
		return fmt.Errorf("there is no unary operand %v", op)
	}
}
func (*SylvaVisitor) VisitIntValue(ctx *parser.IntValueContext) any {
	str := ctx.INT().GetText()
	if i, err := strconv.Atoi(str); err != nil {
		return fmt.Errorf("error while parsing int: %v", err)
	} else {
		return i
	}
}
func (*SylvaVisitor) VisitFloatValue(ctx *parser.FloatValueContext) any {
	str := ctx.FLOAT().GetText()
	if f, err := strconv.ParseFloat(str, 64); err != nil {
		return fmt.Errorf("error while parsing float: %v", err)
	} else {
		return f
	}
}
func (*SylvaVisitor) VisitStringValue(ctx *parser.StringValueContext) any {
	// get the string contents, without the enclosing quotes
	rawStr := ctx.STRING().GetText()[1:]
	rawStr = rawStr[:len(rawStr)-1]
	return rawStr
}
func (*SylvaVisitor) VisitBoolValue(ctx *parser.BoolValueContext) any {
	return ctx.BOOL().GetText() == "true"
}
func (v *SylvaVisitor) VisitParensValue(ctx *parser.ParensValueContext) any {
	return v.Visit(ctx.Expr())
}
func (v *SylvaVisitor) VisitValueExpr(ctx *parser.ValueExprContext) any {
	return v.Visit(ctx.Value())
}
func (v *SylvaVisitor) VisitConcatExpr(ctx *parser.ConcatExprContext) any {
	return ""
}
