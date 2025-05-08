package author

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewEntRepo,
			fx.As(new(Repo)),
		),
	),
)
