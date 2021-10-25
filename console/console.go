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
	cfg       *config.Config
	reader    *bufio.Reader
	term      *terminal.Terminal
	oldState  *terminal.State
	compiler  *compiler.Compiler
	files     []string
	pageID    int
	totPages  int
	width     int
	height    int
	startLine int
	maxLine   int
}

func (co *Console) update() {
	var err error
	co.width, co.height, err = terminal.GetSize(syscall.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	co.Print("\033[H\033[2J\033[?25l") // clear screen, set cursor position, hide cursor
	co.maxLine, err = co.compiler.CompileFile(
		co.files[co.pageID],
		co.term,
		co.startLine,
		co.width,
		co.height)
	co.Printf("\033[%d;0H\033[2K\033[?25h", co.height) // set position, clear line, show cursor
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func (co *Console) Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(co.term, a...)
}

func (co *Console) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(co.term, format, a...)
}

func New(cfg *config.Config) *Console {
	return &Console{
		cfg:      cfg,
		compiler: compiler.New(),
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
	co.startLine = 0
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
						co.startLine--
						if co.startLine < 0 {
							co.startLine = 0
						}
						csi = false
						break loopCSI
					case 'B': // down
						co.startLine++
						if co.startLine > co.maxLine {
							co.startLine = co.maxLine
						}
						/*
							if co.maxLine-co.height > 0 &&
								co.startLine > co.maxLine-co.height {
								co.startLine = co.maxLine - co.height
							}
						*/
						csi = false
						break loopCSI
					case 'C': // right
						co.startLine = 0
						if co.pageID < co.totPages-1 {
							co.pageID++
						}
						csi = false
						break loopCSI

					case 'D': // left
						co.startLine = 0
						if co.pageID > 0 {
							co.pageID--
						}
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
		case 'q', 3: // quit
			return nil

		case ':': // command mode

			co.Printf("\033[%d;0H\033[2K", co.height) // set position and clear line

			co.term.SetPrompt(":")
			cmd, err = co.term.ReadLine()
			if err != nil {
				return err
			}
			if cmd == "q" {
				return nil
			}
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

	resize := make(chan os.Signal, 1)
	go func() {
		for range resize {
			co.update()
		}
	}()
	signal.Notify(resize, syscall.SIGWINCH)

	co.files, err = files.Find(co.cfg.Root, ".md")
	if err != nil {
		return err
	}

	co.totPages = len(co.files)

	return nil
}
