//go:build !arm

package app

import (
	"go.uber.org/fx"
	"raven-club/internal/pkg/api"
	"raven-club/internal/pkg/http"
	"raven-club/internal/pkg/tribute"
	"raven-club/internal/pkg/user"
)

var Module = fx.Module("app",
	api.Module,
	tribute.Module,
	http.Module,
	user.Module)
