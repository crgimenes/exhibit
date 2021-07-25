package main

import (
	"fmt"

	"github.com/crgimenes/exhibit/console"
)

func main() {

	co := console.New()
	err := co.Prepare()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer co.Restore()
	co.InkeyLoop()

}
