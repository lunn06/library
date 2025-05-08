package app

import (
	"go.uber.org/fx"

	"github.com/lunn06/library/review/internal/app/repository"
	"github.com/lunn06/library/review/internal/app/service"
)

var Module = fx.Module("app",
	repository.Module,
	service.Module,
)
