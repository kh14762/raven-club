//go:build arm

package app

import (
	"go.uber.org/fx"
	"raven-club/internal/pkg/api"
	"raven-club/internal/pkg/http"
)

var Module = fx.Module("app",
	api.Module,
	http.Module,
	user.Module,
	auth.Module)
