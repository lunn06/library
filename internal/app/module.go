package app

import (
	"github.com/exepirit/go-template/internal/api"
	"github.com/exepirit/go-template/internal/config"
	"github.com/exepirit/go-template/internal/infrastructure"
	"github.com/exepirit/go-template/internal/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	infrastructure.Module,
	service.Module,
	api.Module,
	fx.Invoke(bootstrap),
)
