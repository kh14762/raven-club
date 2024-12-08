package auth

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"raven-club/internal/pkg/token"
	"raven-club/internal/pkg/types"
	"raven-club/internal/pkg/user"
)

var (
	ErrMissingFields      = errors.New("missing required fields")
	ErrPasswordProcess    = errors.New("failed to process password")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service interface {
	Register(req types.RegisterRequest) (*TokenResponse, error)
	Login(req types.LoginRequest) (*TokenResponse, error)
}

type service struct {
	logger         *zap.Logger
	userService    user.Service
	tokenService   token.Service
	emailLookup    *EmailLookup
	emailValidator *EmailValidator
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewService(logger *zap.Logger, us user.Service, ts token.Service) Service {
	return &service{
		logger:         logger, // Add context to all logs
		userService:    us,
		tokenService:   ts,
		emailLookup:    NewEmailLookup(us),
		emailValidator: NewEmailValidator(),
	}
}

// Register accepts RegisterRequest, creates a User, adds User to UserService, returns a TokenResponse
func (s *service) Register(req types.RegisterRequest) (*TokenResponse, error) {
	// Add request logging
	s.logger.Info("processing registration request",
		zap.String("username", req.Username),
		zap.String("email", req.Email),
	)

	// Generate user Id
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// Validate and normalize email
	if err := s.emailValidator.Validate(req.Email); err != nil {
		s.logger.Warn("invalid email format",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, err
	}
	email := s.emailValidator.Normalize(req.Email)

	// Check if email exists
	if err := s.emailLookup.CheckEmailExists(email); err != nil {
		s.logger.Info("email already exists",
			zap.String("email", req.Email),
		)
		return nil, err
	}

	// Check required fields
	if req.Username == "" || req.Email == "" || req.Password == "" {
		s.logger.Warn("missing required fields",
			zap.String("username", req.Username),
			zap.String("email", req.Email),
		)
		return nil, ErrMissingFields
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password",
			zap.Error(err),
		)
		return nil, ErrPasswordProcess
	}

	u := &types.User{
		Id:       id.ID(),
		Username: req.Username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Generate tokens using TokenManager
	accessToken, err := s.tokenService.GenerateAccessToken(*u)
	if err != nil {
		s.logger.Error("failed to generate access token",
			zap.Error(err),
		)
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(*u)
	if err != nil {
		s.logger.Error("failed to generate refresh token",
			zap.Error(err),
		)
		return nil, err
	}

	// Add the user
	s.userService.Add(*u) // TODO: write logic that adds user to database

	s.logger.Info("user registered successfully",
		zap.String("username", u.Username),
		zap.String("email", u.Email),
		zap.Uint32("id", u.Id),
	)

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Login(req types.LoginRequest) (*TokenResponse, error) {
	s.logger.Info("processing login request",
		zap.String("email", req.Email),
	)

	// Find and validate user
	u, err := s.emailLookup.FindUser(req.Email)
	if err != nil {
		s.logger.Info("login failed: user not found",
			zap.String("email", req.Email),
		)
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		s.logger.Info("login failed: invalid password",
			zap.String("email", req.Email),
		)
		return nil, ErrInvalidCredentials
	}

	// Generate new tokens
	accessToken, err := s.tokenService.GenerateAccessToken(u)
	if err != nil {
		s.logger.Error("failed to generate access token",
			zap.Error(err),
			zap.Uint32("id", u.Id),
		)
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(u)
	if err != nil {
		s.logger.Error("failed to generate refresh token",
			zap.Error(err),
			zap.Uint32("id", u.Id),
		)
		return nil, err
	}

	s.logger.Info("user logged in successfully",
		zap.String("email", req.Email),
		zap.Uint32("id", u.Id),
	)

	// Update user's token
	// TODO call RefreshAccessToken?

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
