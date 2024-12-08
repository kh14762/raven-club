package user

import (
	"go.uber.org/fx"
	"raven-club/internal/pkg/types"
)

var Module = fx.Module("User",
	fx.Provide(
		NewService,
	),
)

type Service interface {
	Get(id uint32) (types.User, bool)
	GetByEmail(email string) (types.User, bool)
	Add(user types.User)
	List() []types.User
	//Update(User) (User, error)
	//Delete(id uint16)
}
