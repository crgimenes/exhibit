package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/crgimenes/exhibit/compiler"
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

	content := []rune(string(b))
	t := compiler.NewTokenizer(content)

	tokens, err := t.Tokenize()
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range tokens {
		fmt.Printf("%02d %v\n", k, v)
	}
}
