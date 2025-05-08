package app

import (
	"go.uber.org/fx"

	"github.com/lunn06/library/bookinfo/internal/app/repository"
	"github.com/lunn06/library/bookinfo/internal/app/service"
)

var Module = fx.Module("app",
	repository.Module,
	service.Module,
)
