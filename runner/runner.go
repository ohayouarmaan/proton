package runner

import (
	"errors"
	"fmt"
	"os"

	"github.com/ohayouarmaan/proton/compiler/positron"
	"github.com/ohayouarmaan/proton/lexer"
	"github.com/ohayouarmaan/proton/parser"
)

type Runner struct {
	File_name string
	Code      string
}

func (r *Runner) Load_program(file_name string) error {
	r.File_name = file_name
	file, err := os.ReadFile(file_name)
	if err != nil {
		return errors.New("error while opening the file")
	}

	r.Code = string(file)
	lexer := lexer.New(r.Code)
	lexer.Generate_Tokens()
	p := parser.Parser{
		Tokens:      lexer.Tokens,
		Current_Idx: 0,
	}
	p.Load_File_Name(r.File_name)
	fmt.Println(lexer.Tokens)
	first_statement := (p.ParseProgram())[0]
	po := positron.Positron{}
	po.Visit_statement(first_statement)
	return nil
}
