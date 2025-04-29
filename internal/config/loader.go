package config

import (
	"errors"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func Load() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg := Config{}

	if err := viper.ReadInConfig(); err != nil {
		switch {
		case errors.As(err, &viper.ConfigFileNotFoundError{}):
			log.Println("[WARN] Configuration file not found")
		default:
			return cfg, err
		}
	}

	return cfg, viper.Unmarshal(&cfg)
}

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
