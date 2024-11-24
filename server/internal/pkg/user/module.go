package user

import (
	"go.uber.org/fx"
	"raven-club/internal/pkg/types"
)

var Module = fx.Module("Service",
	fx.Provide(
		fx.Annotate(
			MakeList,
			fx.As(new(Service)),
		),
	),
)

type Service interface {
	Get(id uint16) (types.User, bool)
	GetByEmail(email string) (types.User, bool)
	Add(user types.User)
	List() []types.User
	//Update(User) (User, error)
	//Delete(id uint16)
}
