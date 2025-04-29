package api

import (
	v1 "github.com/exepirit/go-template/internal/api/v1"
	"github.com/exepirit/go-template/pkg/routing"
	"go.uber.org/fx"
)

var Module = fx.Options(
	v1.Module,
	fx.Provide(NewAPI),
)

type API routing.Bindable

func NewAPI(apiv1 v1.API) API {
	return routing.Route("/api",
		routing.Route("/v1", apiv1),
	)
}
