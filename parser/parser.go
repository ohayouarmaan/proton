/////////////////////////////////////////////////////////////////////////
// WIP-----------------
/////////////////////////////////////////////////////////////////////////

package parser

import (
	"github.com/ohayouarmaan/proton/lexer"
)

type Parser struct {
	Tokens      []lexer.Token
	Current_Idx int
}

type Expression struct {
	Value interface{}
	Type  string
}

type BinaryOp struct {
	lhs      Expression
	operator string
	rhs      Expression
}

type UnaryOp struct {
	operator lexer.TokenType
	rhs      Expression
}

func (p *Parser) Parse() Expression {
	if p.Tokens[p.Current_Idx].TokenType == lexer.Num {
		if p.Tokens[p.Current_Idx+1].TokenType == lexer.Plus {
			return p.Parse_binop()
		}
	}
	return Expression{}
}

func (p *Parser) create_binary_op(tts []lexer.TokenType, precedent_function func() Expression) Expression {
	expr := precedent_function()
	for p.match_tokens(tts) {
		operator_type := p.previous_tokens()
		rhs := precedent_function()
		expr = Expression{
			Value: BinaryOp{
				lhs:      expr,
				operator: string(operator_type),
				rhs:      rhs,
			},
			Type: "BinaryOp",
		}
	}
	return expr
}
func (p *Parser) Parse_binop() Expression {
	return p.create_binary_op([]lexer.TokenType{lexer.Plus, lexer.Minus}, p.parse_factor)
}

func (p *Parser) parse_factor() Expression {
	return p.create_binary_op([]lexer.TokenType{lexer.Divide, lexer.Multiply}, p.parse_unary)
}

func (p *Parser) parse_unary() Expression {
	if p.match_tokens([]lexer.TokenType{lexer.NotUnary, lexer.Minus}) {
		operator := p.previous_tokens()
		return Expression{
			Value: UnaryOp{
				operator: operator,
				rhs:      p.parse_unary(),
			},
			Type: "UnaryExp",
		}
	}
	return p.parse_primary()
}

func (p *Parser) parse_primary() Expression {
	if p.Tokens[p.Current_Idx].Meta_data != nil {
		p.Current_Idx += 1
		return Expression{
			Value: p.Tokens[p.Current_Idx-1].Meta_data.Value,
			Type:  "LiteralValue",
		}
	}
	panic("invalid token.")
}

func (p *Parser) match_tokens(tts []lexer.TokenType) bool {
	for _, tt := range tts {
		if p.Tokens[p.Current_Idx].TokenType == tt {
			p.Current_Idx += 1
			return true
		}
	}
	return false
}

func (p *Parser) previous_tokens() lexer.TokenType {
	return p.Tokens[p.Current_Idx-1].TokenType
}
