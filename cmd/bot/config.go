package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ttrueno/rl2-final/config"
	"github.com/ttrueno/rl2-final/internal/lib/e"
)

func loadConfig() (*config.Config, error) {
	var (
		errmsg = `main.loadConfig`
		cfg    config.Config
	)

	if err := cleanenv.ReadConfig("config.yaml", &cfg); err != nil {
		return nil, e.Wrap(errmsg, err)
	}

	return &cfg, nil
}
