package api

import (
	"go.uber.org/fx"

	"github.com/lunn06/library/review/internal/api/nats"
)

var Module = fx.Module("api",
	nats.Module,
)
