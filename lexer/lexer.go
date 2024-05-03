package lexer

import (
	"strconv"

	"github.com/ohayouarmaan/proton/helpers"
)

type TokenType string

type Literal_Value struct {
	Value interface{}
	Type  string
}

const (
	// Reserved Keywords
	Process TokenType = "proc"
	Return  TokenType = "return"
	Dec     TokenType = "dec"
	Print   TokenType = "print"

	Identifier TokenType = "identifier"
	// Scope Identifiers
	LParen    TokenType = "("
	RParen    TokenType = ")"
	SemiColon TokenType = ";"
	LBrace    TokenType = "{"
	RBrace    TokenType = "}"

	// Datatypes
	Str     TokenType = "str"
	Num     TokenType = "number"
	Integer TokenType = "int"
	String  TokenType = "string"
	Boolean TokenType = "bool"
	True    TokenType = "true"
	False   TokenType = "false"
	Nil     TokenType = "nil"

	// Binary Operators
	Plus     TokenType = "+"
	Minus    TokenType = "-"
	Multiply TokenType = "*"
	Divide   TokenType = "/"
	Modulo   TokenType = "%"

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

	Dot TokenType = "."

	EOF TokenType = "EOF"
)

var KEYWORDS map[string]TokenType = map[string]TokenType{
	"proc":   Process,
	"int":    Integer,
	"string": String,
	"str":    Str,
	"return": Return,
	"dec":    Dec,
	"print":  Print,
	"nil":    Nil,
	"true":   True,
	"false":  False,
}

type Token struct {
	TokenType TokenType
	Lexeme    string
	Meta_data *struct {
		Value Literal_Value
	}
	Line  int
	start int
	end   int
}

type Lexer struct {
	Source_Code  string
	current_line int
	Tokens       []Token
	Current_Idx  int
}

func New(source_code string) Lexer {
	return Lexer{
		Source_Code: source_code,
		Tokens:      []Token{},
		Current_Idx: -1,
	}
}

func getKeywordToken(lexeme string) (TokenType, bool) {
	if KEYWORDS[lexeme] != "" {
		return KEYWORDS[lexeme], false

	}
	return Identifier, true
}

func (l *Lexer) build_keyword() Token {
	built_string := ""
	start := l.Current_Idx
	for helpers.IsAlphaNumeric(string(l.Source_Code[l.Current_Idx])) {
		built_string += string(l.Source_Code[l.Current_Idx])
		l.Current_Idx += 1
	}
	tt, is_identifier := getKeywordToken(built_string)
	if !is_identifier {
		return Token{
			TokenType: tt,
			start:     start,
			Meta_data: nil,
			Lexeme:    built_string,
			Line:      l.current_line,
			end:       l.Current_Idx,
		}
	}
	return Token{
		TokenType: tt,
		start:     start,
		Meta_data: &struct{ Value Literal_Value }{
			Value: Literal_Value{
				Value: built_string,
				Type:  string(Identifier),
			},
		},
		Line:   l.current_line,
		Lexeme: built_string,
		end:    l.Current_Idx,
	}
}

func (l *Lexer) build_string() Token {
	l.Current_Idx += 1
	built_string := ""
	start := l.Current_Idx
	for {
		if l.Current_Idx >= len(l.Source_Code) {
			break
		}
		if string(l.Source_Code[l.Current_Idx]) == "\\" {
			built_string += string(l.Source_Code[l.Current_Idx+1])
			l.Current_Idx += 1
		}
		built_string += string(l.Source_Code[l.Current_Idx])
		l.Current_Idx += 1
		if l.Current_Idx >= len(l.Source_Code) {
			break
		}
		if string(l.Source_Code[l.Current_Idx]) == string('"') {
			break
		}
	}
	return Token{
		TokenType: Str,
		start:     start,
		Lexeme:    built_string,
		Meta_data: &struct{ Value Literal_Value }{
			Value: Literal_Value{Value: built_string, Type: "string"},
		},
		Line: l.current_line,
		end:  l.Current_Idx,
	}
}

func (l *Lexer) build_token(t TokenType, lexeme string) {
	l.Tokens = append(l.Tokens, Token{
		TokenType: t,
		Lexeme:    lexeme,
		start:     l.Current_Idx,
		Meta_data: nil,
		Line:      l.current_line,
		end:       l.Current_Idx + 1,
	})
	l.Current_Idx += 1
}

func (l *Lexer) build_numbers() Token {
	built_number := ""
	dot_count := 0
	start := l.Current_Idx
	current_character := string(l.Source_Code[l.Current_Idx])
	for ((current_character > "0" && current_character < "9") || (current_character == ".")) && l.Current_Idx < len(l.Source_Code) {
		if current_character == "." && dot_count > 0 {
			break
		}
		if current_character == "." {
			dot_count += 1
		}
		built_number += current_character
		l.Current_Idx += 1
		if l.Current_Idx >= len(l.Source_Code) {
			break
		}
		current_character = string(l.Source_Code[l.Current_Idx])
	}
	if value, err := strconv.ParseFloat(built_number, 64); err == nil {
		return Token{
			TokenType: Num,
			start:     start,
			Lexeme:    built_number,
			Meta_data: &struct{ Value Literal_Value }{
				Value: Literal_Value{Value: value, Type: "number"},
			},
			Line: l.current_line,
			end:  l.Current_Idx,
		}
	}
	panic("illegal float value found.")
}

func (l *Lexer) Generate_Tokens() {
	l.Current_Idx = 0
	l.current_line = 0
	for l.Current_Idx >= 0 && l.Current_Idx < len(l.Source_Code) {
		current_character := string(l.Source_Code[l.Current_Idx])
		// if helpers.IsAlphaNumeric(current_character) {
		// 	str := l.build_string()
		// 	fmt.Println("string: ", str)
		// }
		if current_character == " " || current_character == "\t" {
			l.Current_Idx += 1
			continue
		} else if current_character == "\n" {
			l.current_line += 1
		} else if current_character == "(" {
			l.build_token(LParen, "(")
		} else if current_character == ")" {
			l.build_token(RParen, ")")
		} else if current_character == "{" {
			l.build_token(LBrace, "{")
		} else if current_character == "}" {
			l.build_token(RBrace, "}")
		} else if current_character == ";" {
			l.build_token(SemiColon, ";")
		} else if current_character == "." {
			l.build_token(Dot, ".")
		} else if current_character == "=" {
			l.build_token(EqualTo, ".")

		} else if current_character == "+" {
			l.build_token(Plus, "+")
		} else if current_character == "-" {
			l.build_token(Minus, "-")
		} else if current_character == "*" {
			l.build_token(Multiply, "*")
		} else if current_character == "/" {
			l.build_token(Divide, "/")
		} else if current_character == "!" {
			l.build_token(NotUnary, "!")
		} else if current_character == "~" {
			l.build_token(BitWiseNotUnary, "~")
		} else if current_character == "%" {
			l.build_token(Modulo, "%")

		} else if current_character == "\"" {
			l.Tokens = append(l.Tokens, l.build_string())
			l.Current_Idx += 1
		} else if current_character > "0" && current_character < "9" {
			l.Tokens = append(l.Tokens, l.build_numbers())
		} else {
			l.Tokens = append(l.Tokens, l.build_keyword())
		}
	}
	l.Tokens = append(l.Tokens, Token{
		TokenType: EOF,
		start:     len(l.Source_Code) - 1,
		end:       len(l.Source_Code),
		Meta_data: nil,
		Lexeme:    "EOF",
	})
}
