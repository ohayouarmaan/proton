/////////////////////////////////////////////////////////////////////////
// WIP-----------------
/////////////////////////////////////////////////////////////////////////

package parser

import (
	"fmt"
	"strconv"

	"github.com/ohayouarmaan/proton/lexer"
)

type Parser struct {
	Tokens      []lexer.Token
	Current_Idx int
	File_Name   string
}

type Expression struct {
	Value interface{}
	Line  int
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

type Block struct {
	statements []Statement
}

type FunctionArgument struct {
	Type         lexer.Token
	DefaultValue lexer.Literal_Value
}
type FunctionStatement struct {
	Name          string
	Arguments     map[string]FunctionArgument
	FunctionBlock Block
}

type BinaryOp struct {
	Lhs      Expression
	Operator string
	Rhs      Expression
}

type UnaryOp struct {
	Operator lexer.TokenType
	Rhs      Expression
}

func (p *Parser) ParseProgram() []Statement {
	stmts := []Statement{}
	for !p.match_tokens([]lexer.TokenType{lexer.EOF}) {
		stmt := p.parse_stmt()
		if stmt != nil {
			stmts = append(stmts, *stmt)
		}
	}
	return stmts
}

func (p *Parser) Load_File_Name(fn string) {
	p.File_Name = fn
}

func (p *Parser) parse_stmt() *Statement {
	if p.get_current_token().TokenType == lexer.Dec {
		p.Current_Idx += 1
		return p.parse_var_dec_stmt()
	} else if p.match_tokens([]lexer.TokenType{lexer.LBrace}) {
		return p.parse_block_stmt()
	} else if p.Tokens[p.Current_Idx].TokenType == lexer.Print {
		p.Current_Idx += 1
		return p.parse_print_stmt()
	} else if p.match_tokens([]lexer.TokenType{lexer.Process}) {
		return p.parse_function_statement()
	}
	return nil
}

func (p *Parser) parse_block_stmt() *Statement {
	stmts := []Statement{}
	for !p.match_tokens([]lexer.TokenType{lexer.RBrace}) {
		stmts = append(stmts, *p.parse_stmt())
	}
	return &Statement{
		Value: Block{
			statements: stmts,
		},
		Type: "BlockStatement",
	}
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
		if p.match_tokens([]lexer.TokenType{lexer.Num, lexer.String}) {
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
							Line: p.get_current_token().Line,
						},
					},
					Type: "VariableDeclaration",
				}
			}
		} else {
			panic("Expected a type definition after a variable identifier " + p.File_Name + ":line:" + strconv.Itoa(p.get_current_token().Line+1))
		}
	}
	return nil
}

func (p *Parser) parse_function_statement() *Statement {
	p.consume(lexer.Identifier, "Expected a identifier after proc keyword")
	fn_name := p.previous_token()
	p.consume(lexer.LParen, "Expected a '(' after function identifier")
	args := make(map[string]FunctionArgument)
	for p.match_tokens([]lexer.TokenType{lexer.Identifier}) {
		fmt.Println("HI: ", p.get_current_token())
		arg_name := p.previous_token()
		// Make it generic to have any type instead of hardcoding types.
		// Using the above technique we will be able to have user defined types as well.
		if p.match_tokens([]lexer.TokenType{lexer.Num, lexer.String}) {
			arg_type := p.previous_token()
			fmt.Println("Wwow", arg_type)
			if p.match_tokens([]lexer.TokenType{lexer.EqualTo}) {
				if p.match_tokens([]lexer.TokenType{lexer.Num, lexer.Str}) {
					arg_default_value := p.previous_token()
					args[arg_name.Lexeme] = FunctionArgument{
						Type: arg_type,
						DefaultValue: lexer.Literal_Value{
							Value: arg_default_value,
							Type:  string(arg_default_value.TokenType),
						},
					}
				}
			} else {
				args[arg_name.Lexeme] = FunctionArgument{
					Type: arg_type,
					DefaultValue: lexer.Literal_Value{
						Value: nil,
						Type:  "Nil",
					},
				}
			}
		} else {
			panic("You must add argument type in your function's arguments " + p.File_Name + ":line:" + strconv.Itoa(p.get_current_token().Line))
		}
		if p.match_tokens([]lexer.TokenType{lexer.Comma}) {
			continue
		} else {
			p.consume(lexer.RParen, "Missing a closing ')'")
		}
	}
	p.consume(lexer.LBrace, "Expected a '{' after your arguments")
	block_stmt := p.parse_block_stmt().Value.(Block)
	return &Statement{
		Value: FunctionStatement{
			Name:          fn_name.Lexeme,
			Arguments:     args,
			FunctionBlock: block_stmt,
		},
	}
}

func (p *Parser) create_binary_op(tts []lexer.TokenType, precedent_function func() Expression) Expression {
	expr := precedent_function()
	for p.match_tokens(tts) {
		operator_type := p.previous_token_type()
		rhs := precedent_function()
		expr = Expression{
			Value: BinaryOp{
				Lhs:      expr,
				Operator: string(operator_type),
				Rhs:      rhs,
			},
			Line: p.get_current_token().Line,
			Type: "BinaryExp",
		}
	}
	return expr
}
func (p *Parser) Parse_binop() Expression {
	return p.create_binary_op([]lexer.TokenType{lexer.Plus, lexer.Minus}, p.parse_factor)
}

func (p *Parser) parse_factor() Expression {
	return p.create_binary_op([]lexer.TokenType{lexer.Divide, lexer.Multiply, lexer.Modulo}, p.parse_unary)
}

func (p *Parser) parse_unary() Expression {
	if p.match_tokens([]lexer.TokenType{lexer.NotUnary, lexer.Minus}) {
		operator := p.previous_token_type()
		return Expression{
			Value: UnaryOp{
				Operator: operator,
				Rhs:      p.parse_unary(),
			},
			Line: p.get_current_token().Line,
			Type: "UnaryExp",
		}
	}
	return p.parse_primary()
}

func (p *Parser) parse_primary() Expression {
	if p.match_tokens([]lexer.TokenType{lexer.LParen}) {
		value := p.Parse_binop()
		p.consume(lexer.RParen, "You better close the '(' >:(")
		return value
	}
	if p.get_current_token().Meta_data != nil {
		p.Current_Idx += 1
		return Expression{
			Value: p.Tokens[p.Current_Idx-1].Meta_data.Value,
			Line:  p.Tokens[p.Current_Idx-1].Line,
			Type:  "LiteralValue",
		}
	}
	panic("invalid token.")
}

func (p *Parser) match_tokens(tts []lexer.TokenType) bool {
	for _, tt := range tts {
		if p.get_current_token().TokenType == tt {
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
	if p.get_current_token().TokenType == expected_type {
		p.Current_Idx += 1
		return true
	}
	panic(error_message + " " + p.File_Name + ":line:" + strconv.Itoa(p.get_current_token().Line+1))
}

func (p *Parser) get_current_token() lexer.Token {
	return p.Tokens[p.Current_Idx]
}
