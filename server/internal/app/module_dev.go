//go:build !arm

package app

import (
	"go.uber.org/fx"
	"ravens-club/internal/pkg/api"
	"ravens-club/internal/pkg/http"
	"ravens-club/internal/pkg/tribute"
)

var Module = fx.Module("app",
	api.Module,
	tribute.Module,
	http.Module)
