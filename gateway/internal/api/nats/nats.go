package nats

import "github.com/nats-io/nats.go"

func NewConnection(cfg Config) (*nats.Conn, error) {
	return nats.Connect(cfg.URL)
}
