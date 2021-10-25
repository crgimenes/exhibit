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

func (c *Compiler) CompileFile(
	file string,
	w io.Writer,
	startLine, width, height int) (maxLine int, err error) {
	buf := bytes.Buffer{}

	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	in, err := io.ReadAll(f)
	if err != nil {
		return 0, err
	}

	m := strings.Split(strings.ReplaceAll(string(in), "\r\n", "\n"), "\n")
	maxLine = len(m)
	h := maxLine
	if h > height-1 {
		h = height - 1
	}
	if startLine > maxLine {
		startLine = maxLine
	}

	ln := h + startLine
	if ln > maxLine {
		ln = maxLine
	}
	/*
		if startLine < 0 {
			startLine = 0
		}
		if ln < 0 {
			ln = maxLine
		}
	*/
	s := ""
	for k, v := range m[startLine:ln] {
		s += fmt.Sprintf("%2d s:%d, ln:%d m:%d %q\r\n", k+startLine+1, startLine, ln, maxLine, v)
	}

	// s := strings.Join(m[:h], "\n")

	_, err = buf.WriteString(s)
	if err != nil {
		return 0, err
	}

	_, err = buf.WriteTo(w)

	return maxLine, err
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
