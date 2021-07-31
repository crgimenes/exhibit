package console

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/crgimenes/exhibit/compiler"
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
	width    int
	height   int
}

func (co *Console) update() {
	co.Print("\033[H\033[2J")
	co.Print("test print string\r\n")
	var err error
	co.width, co.height, err = terminal.GetSize(syscall.Stdin)
	if err != nil {
		fmt.Println(err)
	}
	co.Printf("%dx%d\r\n", co.height, co.width)
	c := compiler.New()
	c.CompileFile(co.files[0], co.term)
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

func (co *Console) Loop() error {
	var (
		c   rune
		err error
		cmd string
	)

	for {
		co.update()
		c, _, err = co.reader.ReadRune()
		if err != nil {
			return err
		}

		switch c {
		case 27: // ESC, try to parse control sequence
			c, _, err = co.reader.ReadRune()
			if err != nil {
				return err
			}

			switch c {

			case '[': // `ESC[` CSI, Control Sequence Introducer
				csi := true
				s := ""

			loopCSI:
				for csi {
					c, _, err = co.reader.ReadRune()
					if err != nil {
						return err
					}
					switch c {
					case 'A': // up
						csi = false
						break loopCSI
					case 'B': // down
						csi = false
						break loopCSI
					case 'C': // left
						csi = false
						break loopCSI
					case 'D': // right
						csi = false
						break loopCSI

					default:

						if c >= 'a' && c <= 'z' ||
							c >= 'A' && c <= 'Z' ||
							c == '~' {
							csi = false
							break loopCSI
						}
						s += string(c)
					}
				}
			}
		case 'q': // quit
			fallthrough
		case 3: // ^c
			return nil

		case ':': // command mode

			co.Printf("\033[%d;0H\033[2K", co.height) // set position and clear nine

			co.term.SetPrompt(":")
			cmd, err = co.term.ReadLine()
			if err != nil {
				return err
			}
			if cmd == "q" {
				return nil
			}
			co.Printf("\r\ncmd line: %s\r\n", cmd)
			continue
		}
	}
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
