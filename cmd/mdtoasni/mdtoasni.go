package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/gosidekick/goconfig"
)

type config struct {
	File string `json:"file" cfg:"f" cfgRequired:"true"`
}

func main() {
	cfg := &config{}

	err := goconfig.Parse(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

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
