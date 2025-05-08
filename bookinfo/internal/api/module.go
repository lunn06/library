package api

import (
	"go.uber.org/fx"

	"github.com/lunn06/library/bookinfo/internal/api/nats"
)

var Module = fx.Module("api",
	nats.Module,
)
