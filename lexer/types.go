package lexer

type TokenType string

const (
	Process   TokenType = "proc"
	Integer   TokenType = "int"
	LParen    TokenType = "("
	RParen    TokenType = ")"
	Return    TokenType = "return"
	SemiColon TokenType = ";"
	LBrace    TokenType = "{"
	RBrace    TokenType = "}"
)

type Token struct {
	TokenType TokenType
}

type Lexer struct {
	Source_Code string
	Tokens      []Token
	Current_Idx int
}

func (l *Lexer) Generate_Tokens() {

}
