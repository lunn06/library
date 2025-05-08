package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"entgo.io/ent/dialect"

	"github.com/lunn06/library/bookinfo/internal/infrastructure/db/ent"

	_ "github.com/lib/pq"
)

func Connect(cfg Config) (*ent.Client, error) {
	dns := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.URL, cfg.DB, cfg.SslMode,
	)
	client, err := ent.Open(dialect.Postgres, dns)
	if err != nil {
		slog.Error(
			"Failed opening connection to Postgres",
			"connectionString", dns,
			"error", err,
		)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err = client.Schema.Create(ctx); err != nil {
		slog.Error(
			"Failed applying schema to Postgres",
			"error", err,
		)
		return nil, err
	}

	return client, nil
}
