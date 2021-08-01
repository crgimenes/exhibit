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

	"github.com/charmbracelet/glamour"
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

	in, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	r, _ := glamour.NewTermRenderer(
		// detect background color and pick either the default dark or light theme
		glamour.WithAutoStyle(),
		// wrap output at specific width
		glamour.WithWordWrap(40),
	)

	out, err := r.RenderBytes(in)
	if err != nil {
		return err
	}
	_, err = buf.Write(out)
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
