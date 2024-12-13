package session

import (
	"go.uber.org/fx"
)

var Module = fx.Module("Session",
	fx.Provide(
		NewRepository,
		NewService,
	),
	fx.Invoke(func(lc fx.Lifecycle, r *Repository) {
		r.Hook(lc)
	}),
)
