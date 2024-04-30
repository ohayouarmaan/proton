package main

import (
	"fmt"

	"github.com/ohayouarmaan/proton/runner"
)

func main() {
	r := runner.Runner{}
	r.Load_program("./examples/test.pr")
	fmt.Println(r.Code)
}
