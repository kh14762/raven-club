package token

import (
	"github.com/golang-jwt/jwt"
	"go.uber.org/fx"
	"raven-club/internal/pkg/types"
)

var Module = fx.Module("Token",
	fx.Provide(
		NewService,
		NewConfig,
	),
)

type Service interface {
	GenerateAccessToken(user types.User) (string, error)
	GenerateRefreshToken(user types.User) (string, error)
	RefreshAccessToken(token string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
