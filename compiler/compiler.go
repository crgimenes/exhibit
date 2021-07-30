package compiler

import (
	"io"
)

type Compiler struct {
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) CompileFile(file string, w io.Writer) error {
	w.Write([]byte(file))
	w.Write([]byte("\r\n"))
	w.Write([]byte("teste compile writer\r\n"))
	return nil
}
