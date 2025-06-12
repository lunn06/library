package config

import (
	"github.com/lunn06/library/gateway/internal/api/nats"
	"github.com/lunn06/library/gateway/internal/api/server"
	"go.uber.org/fx"
)

type Config struct {
	fx.Out

	Nats   nats.Config
	Server server.Config
}
