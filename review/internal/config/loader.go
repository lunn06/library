package config

import (
	"github.com/kelseyhightower/envconfig"
)

const prefix = "REVIEW"

func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
