package db

import (
	"go.uber.org/fx"

	"github.com/lunn06/library/bookinfo/internal/infrastructure/db/postgres"
)

var Module = fx.Options(
	postgres.Module,
)
