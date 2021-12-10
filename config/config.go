package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Root string `json:"root"`
}

func Load() (*Config, error) {
	// load root directory from the command line
	// if not set, use the default value
	root := "."
	flag.StringVar(&root, "root", root, "root directory")
	flag.Parse()

	cfg := &Config{
		Root: root,
	}
	fmt.Println("root:", cfg.Root)
	os.Exit(0)

	return cfg, nil
}
