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

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/unix"
)

type Console struct {
	reader   *bufio.Reader
	term     *terminal.Terminal
	oldState *terminal.State
}

func getSize(fd int) (width, height int, err error) {
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return -1, -1, err
	}
	return int(ws.Col), int(ws.Row), nil
}

func update(term *terminal.Terminal) {
	term.Write([]byte("teste write"))

	w, h, err := getSize(syscall.Stdin)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%dx%d\r\n", w, h)

}

func New() *Console {
	return &Console{}
}

func (co *Console) Restore() {
	terminal.Restore(syscall.Stdin, co.oldState)
}

func (co *Console) InkeyLoop() {
	var (
		c   rune
		err error
	)

	for {
		c, _, err = co.reader.ReadRune()
		if c == 'q' {
			return
		}
		if c == 'i' {
			err = co.inlineImagesProtocol(co.term, "./nonfree/crg.png")
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
		if unicode.IsControl(c) {
			//fmt.Printf("contol %d\r\n", c)
			fmt.Printf("%c", c)
			continue
		}
		fmt.Printf("%c", c)
		//fmt.Printf("%d ('%c')\r\n", c, c)
	}
}

func (co *Console) inlineImagesProtocol(term *terminal.Terminal, file string) error {
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

	term.Write([]byte(fmt.Sprintf("\033]1337;File=name=%s;inline=1preserveAspectRatio=1;size=%d:", nb, len(encoded))))

	term.Write([]byte(encoded))
	term.Write([]byte("\a"))
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
			update(co.term)
		}
	}()
	signal.Notify(resize, syscall.SIGWINCH)
	return nil
}
