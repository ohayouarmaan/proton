package runner

import (
	"errors"
	"os"

	"github.com/ohayouarmaan/proton/lexer"
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
	return nil
}
