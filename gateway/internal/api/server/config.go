package server

import "time"

type Config struct {
	Address      string        `default:":8080"`
	Debug        bool          `default:"true"`
	Prefork      bool          `default:"false"`
	ReadTimeout  time.Duration `default:"5s" split_words:"true"`
	WriteTimeout time.Duration `default:"5s" split_words:"true"`
	IdleTimeout  time.Duration `default:"15s" split_words:"true"`
}
