package nats

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewConnection,
		NewReviewConsumer,
	),
	fx.Invoke(
		RegisterReviewConsumer,
	),
)
