package main

import (
	"context"
	"errors"
	"github.com/lunn06/library/gateway/internal/api"
	"github.com/lunn06/library/gateway/internal/api/server"
	"github.com/lunn06/library/gateway/internal/config"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"log/slog"
)

var Module = fx.Options(
	config.Module,
	api.Module,

	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	conn *nats.Conn,
	server *server.Server,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Listen(); err != nil {
					slog.Error("Server listen failed", "err", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return errors.Join(
				server.Shutdown(ctx),
				conn.Drain(),
			)
		},
	})
}
