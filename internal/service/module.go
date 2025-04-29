package service

import (
	"github.com/exepirit/go-template/internal/service/greeter"
	"go.uber.org/fx"
)

var Module = fx.Options(
	greeter.Module,
)
