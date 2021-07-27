package config

import (
	"github.com/gosidekick/goconfig"
	_ "github.com/gosidekick/goconfig/json"
)

type Config struct {
	Root string `json:"root" cfg:"root" cfgDefault:"./"`
}

func Load() (*Config, error) {

	cfg := &Config{}

	goconfig.File = "config.json"
	err := goconfig.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
