package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	markdown "crg.eti.br/go/exhibit/markdown"
	"golang.org/x/term"
)

type config struct {
	File string `json:"file"`
}

func main() {
	var err error
	cfg := &config{}

	flag.StringVar(&cfg.File, "f", "-", "file to read")
	flag.Parse()

	f := os.Stdin
	if cfg.File != "-" {
		f, err = os.Open(cfg.File)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
	}

	b, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	width := 80
	if term.IsTerminal(0) {
		// in a term
		width, _, err = term.GetSize(0)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	result := markdown.Render(string(b), width, 6)
	os.Stdout.Write(result)
}
