package postgres

import (
	"context"
	"log/slog"
	"time"

	"entgo.io/ent/dialect"

	"github.com/lunn06/library/book/internal/infrastructure/db/ent"

	_ "github.com/lib/pq"
)

func Connect(cfg Config) (*ent.Client, error) {
	client, err := ent.Open(dialect.Postgres, cfg.Dns)
	if err != nil {
		slog.Error(
			"Failed opening connection to Postgres",
			"connectionString", cfg.Dns,
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
