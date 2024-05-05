package runner

import (
	"errors"
	"fmt"
	"os"

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
	l := lexer.New(r.Code)
	l.Generate_Tokens()
	p := parser.Parser{
		Tokens:      l.Tokens,
		Current_Idx: 0,
	}
	p.Load_File_Name(r.File_name)
	// fmt.Println(l.Tokens)

	program := p.ParseProgram()
	fmt.Println(program)
	// po := positron.Positron{}
	// po.Memory = prototype.New()
	// po.Visit_Program(program)
	return nil
}
