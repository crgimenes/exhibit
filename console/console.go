package console

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"unicode"

	"github.com/crgimenes/exhibit/config"
	"github.com/crgimenes/exhibit/files"
	terminal "golang.org/x/term"
)

type Console struct {
	cfg      *config.Config
	reader   *bufio.Reader
	term     *terminal.Terminal
	oldState *terminal.State
	files    []string
	pageID   int
	totPages int
}

func (co *Console) update() {
	co.Print("test print string")

	w, h, err := terminal.GetSize(syscall.Stdin)
	if err != nil {
		fmt.Println(err)
	}
	co.Printf("%dx%d\r\n", w, h)
}

func (co *Console) Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(co.term, a...)
}

func (co *Console) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(co.term, format, a...)
}

func New(cfg *config.Config) *Console {
	return &Console{
		cfg: cfg,
	}
}

func (co *Console) Restore() {
	terminal.Restore(syscall.Stdin, co.oldState)
}

func (co *Console) Loop() {
	var (
		c   rune
		err error
	)

	co.update()

	for {
		c, _, err = co.reader.ReadRune()
		if c == 'q' {
			return
		}
		if c == 'i' {
			err = co.inlineImagesProtocol("./nonfree/crg.png")
			if err != nil {
				fmt.Println(err)
				return
			}
			continue
		}
		if c == ':' {
			fmt.Printf(":")
			line, err := co.term.ReadLine()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("line: %s\r\n", line)
			continue
		}
		if c == 'h' { // left in normal mode
			co.update()
			continue
		}
		if unicode.IsControl(c) {
			//fmt.Printf("contol %d\r\n", c)
			fmt.Printf("%c", c)
			continue
		}
		fmt.Printf("%c", c)
		//fmt.Printf("%d ('%c')\r\n", c, c)
	}
}

func (co *Console) inlineImagesProtocol(file string) error {
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

	co.Printf("\033]1337;File=name=%s;inline=1preserveAspectRatio=1;size=%d:", nb, len(encoded))
	co.Print(encoded)
	co.Print("\a")
	return nil
}

func (co *Console) Prepare() (err error) {
	co.oldState, err = terminal.MakeRaw(syscall.Stdin)
	if err != nil {
		return err
	}

	co.reader = bufio.NewReader(os.Stdin)

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	co.term = terminal.NewTerminal(screen, "")

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		<-sigTerm
		co.Restore()
		os.Exit(0)
	}()

	resize := make(chan os.Signal)
	go func() {
		for range resize {
			co.update()
		}
	}()
	signal.Notify(resize, syscall.SIGWINCH)

	co.files, err = files.Find(co.cfg.Root, ".md")
	return err
}
