package infrastructure

import (
	"go.uber.org/fx"

	"github.com/lunn06/library/bookinfo/internal/infrastructure/db"
)

var Module = fx.Module("infrastructure",
	db.Module,
)
