package user

import (
	"go.uber.org/zap"
	"raven-club/internal/pkg/types"
)

type service struct {
	logger *zap.Logger
	Users  []types.User `json:"users"`
}

func NewService(logger *zap.Logger) Service {
	return &service{
		logger: logger,
		Users:  make([]types.User, 0),
	}
}

func (s *service) Get(id uint32) (types.User, bool) {
	for _, user := range s.Users {
		if user.ID == id {
			return user, true
		}
	}
	return types.User{}, false
}

func (s *service) GetByEmail(email string) (types.User, bool) {
	for _, user := range s.Users {
		if user.Email == email {
			return user, true
		}
	}
	return types.User{}, false
}

func (s *service) Add(user types.User) {
	s.Users = append(s.Users, user)
}

func (s *service) List() []types.User {
	return s.Users
}
