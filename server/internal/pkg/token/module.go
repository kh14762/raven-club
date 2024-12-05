package token

import "go.uber.org/fx"

var Module = fx.Provide("Token",
	fx.Provide(
		NewService,
		NewConfig,
	),
)
