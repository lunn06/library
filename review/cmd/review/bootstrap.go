package main

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.uber.org/fx"

	"github.com/lunn06/library/review/internal/api"
	"github.com/lunn06/library/review/internal/app"
	"github.com/lunn06/library/review/internal/config"
	"github.com/lunn06/library/review/internal/infrastructure"
)

var Module = fx.Options(
	config.Module,
	app.Module,
	api.Module,
	infrastructure.Module,

	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	conn *nats.Conn,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return conn.Drain()
		},
	})
}
