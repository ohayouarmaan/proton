package lexer

import (
	"github.com/ohayouarmaan/proton/helpers"
)

type TokenType string

const (
	// Reserved Keywords
	Process TokenType = "proc"
	Integer TokenType = "int"
	Return  TokenType = "return"

	// Scope Identifiers
	LParen    TokenType = "("
	RParen    TokenType = ")"
	SemiColon TokenType = ";"
	LBrace    TokenType = "{"
	RBrace    TokenType = "}"

	// Datatypes
	Str TokenType = "string"
	Num TokenType = "number"

	// Binary Operators
	Plus     TokenType = "+"
	Minus    TokenType = "-"
	Multiply TokenType = "*"
	Divide   TokenType = "/"

	// Comparison Operator
	EqualTo          TokenType = "=="
	LesserThan       TokenType = "<"
	LesserThanEqual  TokenType = "<="
	GreaterThan      TokenType = ">"
	GreaterThanEqual TokenType = ">="

	// Unary
	NotUnary        TokenType = "!"
	BitWiseNotUnary TokenType = "~"
	Pointer         TokenType = "&"

	EOF TokenType = "EOF"
)

type Token struct {
	TokenType TokenType
	Lexeme    string
	start     int
	end       int
}

type Lexer struct {
	Source_Code string
	Tokens      []Token
	Current_Idx int
}

func New(source_code string) Lexer {
	return Lexer{
		Source_Code: source_code,
		Tokens:      []Token{},
		Current_Idx: -1,
	}
}

func (l *Lexer) build_string() Token {
	built_string := ""
	start := l.Current_Idx
	for helpers.IsAlphaNumeric(string(l.Source_Code[l.Current_Idx])) {
		built_string += string(l.Source_Code[l.Current_Idx])
		l.Current_Idx += 1
	}
	return Token{
		TokenType: Str,
		start:     start,
		Lexeme:    built_string,
		end:       l.Current_Idx,
	}
}

func (l *Lexer) build_token(t TokenType, lexeme string) {
	l.Tokens = append(l.Tokens, Token{
		TokenType: t,
		Lexeme:    lexeme,
		start:     l.Current_Idx,
		end:       l.Current_Idx + 1,
	})
	l.Current_Idx += 1
}

func (l *Lexer) build_numbers() Token {
	built_number := ""
	dot_count := 0
	start := l.Current_Idx
	current_character := string(l.Source_Code[l.Current_Idx])
	for (current_character > "0" && current_character < "9") || (current_character == ".") {
		if dot_count > 0 {
			if current_character == "." {
				break
			}
		}
		if current_character == "." {
			dot_count += 1
		}
		built_number += current_character
		l.Current_Idx += 1
		current_character = string(l.Source_Code[l.Current_Idx])
	}
	return Token{
		TokenType: Num,
		start:     start,
		Lexeme:    built_number,
		end:       l.Current_Idx,
	}
}

func (l *Lexer) Generate_Tokens() {
	l.Current_Idx = 0
	for l.Current_Idx >= 0 && l.Current_Idx < len(l.Source_Code) {
		current_character := string(l.Source_Code[l.Current_Idx])
		// if helpers.IsAlphaNumeric(current_character) {
		// 	str := l.build_string()
		// 	fmt.Println("string: ", str)
		// }
		if current_character == " " || current_character == "\n" || current_character == "\t" {
			l.Current_Idx += 1
			continue
		} else if current_character == "(" {
			l.build_token(LParen, "(")
		} else if current_character == ")" {
			l.build_token(RParen, ")")
		} else if current_character == "{" {
			l.build_token(LBrace, "{")
		} else if current_character == "}" {
			l.build_token(RBrace, "}")
		} else if current_character > "0" && current_character < "9" {
			l.Tokens = append(l.Tokens, l.build_numbers())
		} else if current_character == ";" {
			l.build_token(SemiColon, ";")
		} else {
			l.Tokens = append(l.Tokens, l.build_string())
		}
	}
	l.Tokens = append(l.Tokens, Token{
		TokenType: EOF,
		start:     len(l.Source_Code) - 1,
		end:       len(l.Source_Code),
		Lexeme:    "EOF",
	})
}
