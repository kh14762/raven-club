package user

import (
	"go.uber.org/fx"
)

var Module = fx.Module("User",
	fx.Provide(
		NewService,
		NewRepository,
	),
	fx.Invoke(func(lc fx.Lifecycle, r Repository) {
		r.Hook(lc)
	}),
)
