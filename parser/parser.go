/////////////////////////////////////////////////////////////////////////
// WIP-----------------
/////////////////////////////////////////////////////////////////////////

package parser

import (
	"fmt"

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

type Statement struct {
	Value interface{}
	Type  string
}

type VarDeclarationStatement struct {
	Name  string
	Type  lexer.TokenType
	Value Expression
}

type PrintStatement struct {
	Value Expression
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

func (p *Parser) ParseProgram() []Statement {
	stmts := []Statement{}
	for p.Tokens[p.Current_Idx].TokenType != lexer.EOF {
		stmt := p.parse_stmt()
		if stmt != nil {
			stmts = append(stmts, *stmt)
		} else {
			fmt.Println("currentToken: ", p.Tokens[p.Current_Idx])
		}
	}
	return stmts
}

func (p *Parser) parse_stmt() *Statement {
	if p.Tokens[p.Current_Idx].TokenType == lexer.Dec {
		p.Current_Idx += 1
		return p.parse_var_dec_stmt()
	} else if p.Tokens[p.Current_Idx].TokenType == lexer.Print {
		fmt.Println("Parsing print statement")
		p.Current_Idx += 1
		return p.parse_print_stmt()
	}
	return nil
}

func (p *Parser) parse_print_stmt() *Statement {
	value := p.Parse_binop()
	p.consume(lexer.SemiColon, "Expected a ';' after a print statement.")
	generated_stmt := Statement{
		Value: PrintStatement{
			Value: value,
		},
		Type: "PrintStmt",
	}
	return &generated_stmt
}

func (p *Parser) parse_var_dec_stmt() *Statement {
	if p.consume(lexer.Identifier, "Expected an Identifier") {
		new_var_name := p.previous_token()
		if p.match_tokens([]lexer.TokenType{lexer.Integer, lexer.String}) {
			new_var_type := p.previous_token_type()
			if p.match_tokens([]lexer.TokenType{lexer.EqualTo, lexer.String}) {
				new_var_value := p.Parse_binop()
				p.consume(lexer.SemiColon, "Expected a ';' after a variable assignment.")
				return &Statement{
					Value: VarDeclarationStatement{
						Name:  new_var_name.Lexeme,
						Type:  new_var_type,
						Value: new_var_value,
					},
					Type: "VariableDeclaration",
				}
			} else {
				p.consume(lexer.SemiColon, "Expected a ';' after a variable declaration.")
				return &Statement{
					Value: VarDeclarationStatement{
						Name: new_var_name.Lexeme,
						Type: new_var_type,
						Value: Expression{
							Value: lexer.Literal_Value{
								Value: nil,
								Type:  "Nil",
							},
							Type: "LiteralValue",
						},
					},
					Type: "VariableDeclaration",
				}
			}
		}
	}
	return nil
}

func (p *Parser) create_binary_op(tts []lexer.TokenType, precedent_function func() Expression) Expression {
	expr := precedent_function()
	for p.match_tokens(tts) {
		operator_type := p.previous_token_type()
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
		operator := p.previous_token_type()
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

func (p *Parser) previous_token_type() lexer.TokenType {
	return p.Tokens[p.Current_Idx-1].TokenType
}

func (p *Parser) previous_token() lexer.Token {
	return p.Tokens[p.Current_Idx-1]
}

func (p *Parser) consume(expected_type lexer.TokenType, error_message string) bool {
	if p.Tokens[p.Current_Idx].TokenType == expected_type {
		p.Current_Idx += 1
		return true
	}
	panic(error_message)
}

func (p *Parser) get_current_token() lexer.Token {
	return p.Tokens[p.Current_Idx]
}
