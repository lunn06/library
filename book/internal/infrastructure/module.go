package infrastructure

import (
	"go.uber.org/fx"

	"github.com/lunn06/library/book/internal/infrastructure/db"
)

var Module = fx.Module("infrastructure",
	db.Module,
)
