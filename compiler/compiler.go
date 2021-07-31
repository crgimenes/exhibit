package compiler

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/crgimenes/exhibit/lex"
)

type Compiler struct {
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) CompileFile(file string, w io.Writer) error {
	buf := bytes.Buffer{}
	buf.WriteString(file)
	buf.WriteString("\r\n")

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := lex.Parse(f)
	if err != nil {
		return err
	}

	for _, v := range t {
		buf.WriteString(fmt.Sprintf("%q", v.Literal))
		buf.WriteString(" -> ")
		buf.WriteString(v.Type)
		buf.WriteString("|")
		buf.WriteString("\r\n")
	}

	/*
		r := bufio.NewReader(f)

		var o rune
		for {
			o, _, err = r.ReadRune()
			if err != nil {
				break
			}

			buf.WriteRune(o)
		}

		if err != io.EOF {
			return err
		}
	*/
	buf.WriteTo(w)
	return nil
}

func (co *Compiler) inlineImagesProtocol(file string, w io.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	nb := base64.StdEncoding.EncodeToString([]byte(filepath.Base(file)))

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(
		"\033]1337;File=name=%s;inline=1preserveAspectRatio=1;size=%d:",
		nb,
		len(encoded)))

	buf.WriteString(encoded)
	buf.WriteString("\a")
	buf.WriteTo(w)

	return nil
}
