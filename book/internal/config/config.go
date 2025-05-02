package config

import (
	"go.uber.org/fx"

	"github.com/lunn06/library/book/internal/api/nats"
	"github.com/lunn06/library/book/internal/infrastructure/db/postgres"
)

type Config struct {
	fx.Out

	Nats     nats.Config
	Postgres postgres.Config
}
