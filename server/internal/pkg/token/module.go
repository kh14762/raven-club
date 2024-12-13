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
	GenerateAccessJwt(user types.User) (string, error)
	GenerateRefreshJwt(user types.User) (string, error)
	RefreshAccessJwt(token string) (string, error)
	ValidateJwt(token string) (*jwt.Token, error)
}
