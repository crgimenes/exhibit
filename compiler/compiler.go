package compiler

import (
	"fmt"
	"io"
)

type Compiler struct {
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) CompileFile(file string, w io.Writer) error {
	fmt.Println(file)
	w.Write([]byte("teste compile writer\r\n"))
	return nil
}
