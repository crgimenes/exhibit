package console

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/crgimenes/exhibit/config"
	"github.com/crgimenes/exhibit/files"
	"github.com/crgimenes/exhibit/markdown"
	"github.com/pelletier/go-toml/v2"
	terminal "golang.org/x/term"
)

type Console struct {
	cfg       *config.Config
	reader    *bufio.Reader
	term      *terminal.Terminal
	oldState  *terminal.State
	files     []string
	filesRaw  [][]byte
	pageID    int
	totPages  int
	width     int
	height    int
	startLine int
	maxLine   int
}

func ShowFile(
	file string,
	fileRaw []byte,
	w io.Writer,
	startLine, width, height int) (maxLine int, err error) {
	buf := bytes.Buffer{}

	result := markdown.Render(string(fileRaw), width, 6)

	m := strings.Split(strings.ReplaceAll(string(result), "\r\n", "\n"), "\n")
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

	if startLine > (maxLine - h) {
		startLine = maxLine - h
	}

	for _, v := range m[startLine:ln] {
		_, _ = buf.WriteString(v + "\n")
	}

	_, err = buf.WriteTo(w)

	return maxLine, err
}

func (co *Console) update() {
	var err error
	co.width, co.height, err = terminal.GetSize(syscall.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	co.Print("\033[H\033[2J\033[?25l") // clear screen, set cursor position, hide cursor
	co.maxLine, err = ShowFile(
		co.files[co.pageID],
		co.filesRaw[co.pageID],
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
						co.up()

						csi = false
						break loopCSI
					case 'B': // down
						co.down()

						csi = false
						break loopCSI
					case 'C': // right
						co.right()
						csi = false
						break loopCSI

					case 'D': // left
						co.left()
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
		case 'h':
			co.left()
		case 'j':
			co.down()
		case 'k':
			co.up()
		case 'l':
			co.right()
		case 'q', 3: // quit
			return nil
		case 'e': // edit
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vim"
			}
			cmd := exec.Command(editor, co.files[co.pageID])
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			co.filesRaw[co.pageID], err = prepareFile(co.files[co.pageID])
			co.update()
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
			if cmd == "!" {
				// Hello GTFOBins hackers!
				cmd = os.Getenv("SHELL")
				if cmd == "" {
					cmd = "/bin/sh"
				}
				co.Restore()
				cmd := exec.Command(cmd)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				err = cmd.Run()
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
				terminal.MakeRaw(syscall.Stdin)
				co.filesRaw[co.pageID], err = prepareFile(co.files[co.pageID])
				co.update()
			}

			continue
		}
	}
}

func (co *Console) left() {
	co.startLine = 0
	if co.pageID > 0 {
		co.pageID--
	}
}

func (co *Console) right() {
	co.startLine = 0
	if co.pageID < co.totPages-1 {
		co.pageID++
	}
}

func (co *Console) down() {
	co.startLine++
	if co.startLine > co.maxLine {
		co.startLine = co.maxLine
	}

	if co.maxLine-co.height > 0 &&
		co.startLine > co.maxLine-co.height {
		co.startLine = co.maxLine - co.height
	}
}

func (co *Console) up() {
	co.startLine--
	if co.startLine < 0 {
		co.startLine = 0
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
	co.filesRaw = make([][]byte, co.totPages)

	for i, f := range co.files {
		b, err := prepareFile(f)
		if err != nil {
			return err
		}
		co.filesRaw[i] = b
	}

	return nil
}

func prepareFile(f string) ([]byte, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	body, title, _, err := parseHeader(b)
	if err != nil {
		return nil, err
	}

	if title != "" {
		b = []byte(fmt.Sprintf("# %s\n\n%s", title, body))
	}
	return b, nil
}

func parseHeader(b []byte) (body []byte, title string, draft bool, err error) {
	aux := strings.Split(string(b), "+++")

	if len(aux) > 1 {
		d := make(map[string]any)
		err = toml.Unmarshal([]byte(aux[1]), &d)
		if err != nil {
			log.Fatal(err)
		}

		title = getTitle(d)
		draft = getDraft(d)
	}
	if len(aux) > 2 {
		body = []byte(aux[2])
	}

	return
}

func getTitle(d map[string]any) string {
	title := ""
	titleAux, ok := d["title"]
	if !ok {
		titleAux = ""
	}
	switch titleAux.(type) {
	case string:
		title = titleAux.(string)
	}
	return title
}

func getDraft(d map[string]any) bool {
	draft := false
	draftAux, ok := d["draft"]
	if !ok {
		draftAux = false
	}
	switch draftAux.(type) {
	case bool:
		draft = draftAux.(bool)
	default:
		draft = false
	}
	return draft
}
