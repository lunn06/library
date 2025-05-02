package nats

import "time"

type Config struct {
	Url            string        `default:"nats://127.0.0.1:4222" split_words:"true"`
	RequestTimeout time.Duration `default:"5s" split_words:"true"`
}
