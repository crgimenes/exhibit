package config

import (
	"flag"
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

	return cfg, nil
}
