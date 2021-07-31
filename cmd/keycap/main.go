package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	terminal "golang.org/x/term"
)

func update() {}

func main() {
	oldState, err := terminal.MakeRaw(syscall.Stdin)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer terminal.Restore(syscall.Stdin, oldState)

	reader := bufio.NewReader(os.Stdin)

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := terminal.NewTerminal(screen, "")

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		<-sigTerm
		terminal.Restore(syscall.Stdin, oldState)
		os.Exit(0)
	}()

	resize := make(chan os.Signal)
	go func() {
		for range resize {
			update()
		}
	}()
	signal.Notify(resize, syscall.SIGWINCH)

	var c rune
	for {
		update()
		c, _, err = reader.ReadRune()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(term, "%d -> %c\r\n", c, c)

		if c == 'q' || c == 3 {
			return
		}
	}
}
