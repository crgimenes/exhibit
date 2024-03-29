package main

import (
	"fmt"
	"os"

	"crg.eti.br/go/exhibit/config"
	"crg.eti.br/go/exhibit/console"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	co := console.New(cfg)

	err = co.Prepare()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer co.Restore()
	err = co.Loop()
	if err != nil {
		fmt.Println(err)
	}
}
