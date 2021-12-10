package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type config struct {
	File string `json:"file"`
}

func main() {
	cfg := &config{}

	flag.StringVar(&cfg.File, "f", "", "file to read")
	flag.Parse()

	f, err := os.Open(cfg.File)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)

	for {
		c, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}

		fmt.Printf("%q\r\n", c)
	}
}
