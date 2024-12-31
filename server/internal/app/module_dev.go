//go:build !arm

package app

import (
	"go.uber.org/fx"
	"raven-club/internal/pkg/api"
	"raven-club/internal/pkg/auth"
	"raven-club/internal/pkg/http"
	"raven-club/internal/pkg/session"
	"raven-club/internal/pkg/user"
)

var Module = fx.Module("app",
	api.Module,
	http.Module,
	user.Module,
	auth.Module,
	session.Module)
