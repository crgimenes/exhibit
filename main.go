package main

import (
	"fmt"
	"os"

	"github.com/crgimenes/exhibit/console"
)

func main() {
	co := console.New()

	err := co.Prepare()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer co.Restore()
	co.InkeyLoop()
}
