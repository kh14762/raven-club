package user

import (
	"github.com/google/uuid"
	"go.uber.org/fx"
	"raven-club/internal/pkg/types"
)

var Module = fx.Module("Service",
	fx.Provide(
		fx.Annotate(
			// MakeList, // TODO: some sort of init of user data structure here
			fx.As(new(Service)),
		),
	),
)

type Service interface {
	Get(id uuid.UUID) (types.User, bool)
	GetByEmail(email string) (types.User, bool)
	Add(user types.User)
	List() []types.User
	//Update(User) (User, error)
	//Delete(id uint16)
}
