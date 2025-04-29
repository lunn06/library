package v1

import (
	"github.com/exepirit/go-template/pkg/routing"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewGreeterEndpoints),
	fx.Provide(NewAPI),
)

type API routing.Bindable

func NewAPI(greeter *GreeterEndpoints) API {
	return routing.Union(
		routing.Route("/greeter", greeter),
	)
}
