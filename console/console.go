package console

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"unicode"

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
}

func (co *Console) update() {
	co.Print("test print string\r\n")

	w, h, err := terminal.GetSize(syscall.Stdin)
	if err != nil {
		fmt.Println(err)
	}
	co.Printf("%dx%d\r\n", w, h)
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

	co.update()

	for {
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
						fmt.Print("up\r\n")
						csi = false
						break loopCSI
					case 'B': // down
						fmt.Print("down\r\n")
						csi = false
						break loopCSI
					case 'C': // left
						fmt.Print("left\r\n")
						csi = false
						break loopCSI
					case 'D': // right
						fmt.Print("right\r\n")
						csi = false
						break loopCSI

					default:

						if c >= 'a' && c <= 'z' ||
							c >= 'A' && c <= 'Z' ||
							c == '~' {

							fmt.Printf("ESC[%s%c\r\n", s, c)
							csi = false
							break loopCSI
						}
						s += string(c)
					}
				}
			default:
				fmt.Printf("ESC%c", c)
			}
		case 'q': // quit
			return nil
		case ':': // command mode

			fmt.Printf(":")
			cmd, err = co.term.ReadLine()
			if err != nil {
				return err
			}
			fmt.Printf("\r\ncmd line: %s\r\n", cmd)
			continue

		default:
			if unicode.IsControl(c) {
				if c == 3 { // ^c terminate
					return nil
				}
				fmt.Printf("contol %d\r\n", c)
				continue
			}
			fmt.Printf("%c\r\n", c)
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
