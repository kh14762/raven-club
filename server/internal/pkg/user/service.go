package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"raven-club/internal/pkg/types"
)

// Service TODO: change api to match repo interface
type Service interface {
	CreateUser(ctx *gin.Context, user types.User) error
	GetUserByID(ctx *gin.Context, id string) (*types.User, error)
	GetUserByUsername(ctx *gin.Context, username string) (*types.User, error)
	GetUserByEmail(ctx *gin.Context, email string) (*types.User, error)
	UpdateUser(ctx *gin.Context, user types.User) error
	DeleteUser(ctx *gin.Context, id string) error
	ListUsers(ctx *gin.Context) ([]types.User, error)
}

type service struct {
	logger     *zap.Logger
	repository Repository
}

func NewService(logger *zap.Logger, r Repository) Service {
	return &service{
		logger:     logger,
		repository: r,
	}
}

func (s *service) CreateUser(ctx *gin.Context, user types.User) error {
	if err := s.repository.CreateUser(ctx, user); err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return err
	}
	s.logger.Info("user created", zap.String("username", user.Username))
	return nil
}

func (s *service) GetUserByID(ctx *gin.Context, id string) (*types.User, error) {
	user, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user by ID", zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *service) GetUserByUsername(ctx *gin.Context, username string) (*types.User, error) {
	user, err := s.repository.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.Error("failed to get user by username", zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *service) GetUserByEmail(ctx *gin.Context, email string) (*types.User, error) {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to get user by email", zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *service) UpdateUser(ctx *gin.Context, user types.User) error {
	if err := s.repository.UpdateUser(ctx, user); err != nil {
		s.logger.Error("failed to update user", zap.Error(err))
		return err
	}
	s.logger.Info("user updated", zap.String("username", user.Username))
	return nil
}

func (s *service) DeleteUser(ctx *gin.Context, id string) error {
	if err := s.repository.DeleteUser(ctx, id); err != nil {
		s.logger.Error("failed to delete user", zap.Error(err))
		return err
	}
	s.logger.Info("user deleted", zap.String("id", id))
	return nil
}

func (s *service) ListUsers(ctx *gin.Context) ([]types.User, error) {
	users, err := s.repository.ListUsers(ctx)
	if err != nil {
		s.logger.Error("failed to list users", zap.Error(err))
		return nil, err
	}
	return users, nil
}
