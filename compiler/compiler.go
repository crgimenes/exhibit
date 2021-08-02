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
	"strings"
)

type Compiler struct {
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) CompileFile(file string, w io.Writer, width, height int) error {
	buf := bytes.Buffer{}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	in, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	m := strings.Split(strings.ReplaceAll(string(in), "\r\n", "\n"), "\n")
	h := len(m)
	if h > height-1 {
		h = height - 1
	}

	s := ""
	for k, v := range m[:h] {
		s += fmt.Sprintf("%d %q\r\n", k, v)
	}

	//s := strings.Join(m[:h], "\n")

	_, err = buf.WriteString(s)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)

	return err
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
