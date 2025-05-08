package nats

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewConnection,
		NewAuthorConsumer,
		NewBookConsumer,
		NewGenreConsumer,
	),
	fx.Invoke(
		RegisterAuthorConsumer,
		RegisterBookConsumer,
		RegisterGenreConsumer,
	),
)
