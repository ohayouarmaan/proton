/////////////////////////////////////////////////////////////////////////
// WIP-----------------
// POSITRON: Virtual machine (emulates coulomb assembly)
// POSITRON: Will also be able to interpret proton
/////////////////////////////////////////////////////////////////////////

package positron

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ohayouarmaan/proton/lexer"
	"github.com/ohayouarmaan/proton/parser"
)

type Positron struct {
	Pc    int
	Stack []interface{}
}

func (p *Positron) Compile_to_coulomb(ast *[]parser.Statement) {

}

func (p *Positron) visit_expression(node parser.Expression) lexer.Literal_Value {
	if node.Type == "BinaryExp" {
		return p.visit_bin_expression(node)
	} else if node.Type == "UnaryExp" {
		p.visit_unary_expression(node)
	} else if node.Type == "LiteralValue" {
		return p.visit_literal_expression(node)
	}
	panic("UNIMPLEMENTED")
}

func (p *Positron) Visit_statement(node parser.Statement) {
	if node.Type == "PrintStmt" {
		fmt.Println(p.visit_expression(node.Value.(parser.PrintStatement).Value).Value)
	}
}

func (p *Positron) visit_bin_expression(node parser.Expression) lexer.Literal_Value {
	if val, ok := node.Value.(parser.BinaryOp); ok {
		lhs := val.Lhs
		operator := val.Operator
		rhs := val.Rhs

		lhs_parsed := p.visit_expression(lhs)
		rhs_parsed := p.visit_expression(rhs)
		switch operator {
		case string(lexer.Plus):
			return lexer.Literal_Value{
				Value: lhs_parsed.Value.(float64) + rhs_parsed.Value.(float64),
				Type:  "number",
			}
		case string(lexer.Multiply):
			return lexer.Literal_Value{
				Value: lhs_parsed.Value.(float64) * rhs_parsed.Value.(float64),
				Type:  "number",
			}
		case string(lexer.Minus):
			return lexer.Literal_Value{
				Value: lhs_parsed.Value.(float64) - rhs_parsed.Value.(float64),
				Type:  "number",
			}
		case string(lexer.Divide):
			return lexer.Literal_Value{
				Value: lhs_parsed.Value.(float64) / rhs_parsed.Value.(float64),
				Type:  "number",
			}
		case string(lexer.Modulo):
			if val, ok := lhs_parsed.Value.(float64); ok && math.Mod(val, 1) == 0 {
				if val2, ok2 := rhs_parsed.Value.(float64); ok2 {
					return lexer.Literal_Value{
						Value: math.Mod(val, val2),
						Type:  "number",
					}
				}
				panic("Can not parse to float64")
			} else {
				panic("operator 'a' can not be a float in the modulo operator. in line: " + strconv.Itoa(node.Line))
			}
		}
	}
	panic("UNIMPLEMENTED")
}

func (p *Positron) visit_unary_expression(node parser.Expression) lexer.Literal_Value {
	if value, ok := node.Value.(parser.UnaryOp); ok {
		if value.Operator == "-" {
			exp := p.visit_expression(value.Rhs)
			return lexer.Literal_Value{
				Value: -exp.Value.(float64),
				Type:  "number",
			}
		}
	}
	panic("UNREACHABLE")
}

func (p *Positron) visit_literal_expression(node parser.Expression) lexer.Literal_Value {
	if val, ok := node.Value.(lexer.Literal_Value); ok {
		return val
	}
	panic("INVALID LITERAL VALUE")
}
