package user

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"raven-club/internal/pkg/types"
)

type service struct {
	logger *zap.Logger
}

type List struct {
	Id    uuid.UUID    `json:"id"`
	Users []types.User `json:"users"`
}

func (l *List) List() []types.User {
	return l.Users
}

func (l *List) Get(id uuid.UUID) (types.User, bool) {
	for _, user := range l.Users {
		if user.Id == id {
			return user, true
		}
	}
	return types.User{}, false
}

func (l *List) GetByEmail(email string) (types.User, bool) {
	for _, user := range l.Users {
		if user.Email == email {
			return user, true
		}
	}
	return types.User{}, false
}

func (l *List) Add(user types.User) {
	user.Id = l.Id

	l.Users = append(l.Users, user)
}

//func MakeList() *List {
//	return &List{
//		Id:    0,
//		Users: make([]types.User, 0),
//	}
//}
