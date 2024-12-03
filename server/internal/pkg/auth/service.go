package auth

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"raven-club/internal/pkg/types"
	"raven-club/internal/pkg/user"
	"time"
)

// Error definitions
var (
	ErrMissingFields      = errors.New("missing required fields")
	ErrPasswordProcess    = errors.New("failed to process password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service interface {
	Register(req types.RegisterRequest) (*TokenResponse, error)
	Login(req types.LoginRequest) (*TokenResponse, error)
	ValidateToken(token string) (*jwt.Token, error)
	RefreshToken(refreshToken string) (*TokenResponse, error)
}

type service struct {
	userService    user.Service
	emailLookup    *EmailLookup
	emailValidator *EmailValidator
	tokenManager   *TokenManager
	logger         *zap.Logger
}

func NewService(us user.Service, config *Config, logger *zap.Logger) Service {
	return &service{
		userService:    us,
		emailLookup:    NewEmailLookup(us),
		emailValidator: NewEmailValidator(),
		tokenManager:   NewTokenManager(config.SecretKey),
		logger:         logger.With(zap.String("service", "auth")), // Add context to all logs
	}
}

// Register accepts RegisterRequest, creates a User, offers User to UserService, returns an access token
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

	u := &types.User {
		Id: id,
		Username: req.Username,
		Email: email,
		Password: string(hashedPassword),
	}

	// Generate tokens using TokenManager
	accessToken, err := s.tokenManager.GenerateAccessToken(*u)
	if err != nil {
		s.logger.Error("failed to generate access token",
			zap.Error(err),
		)
		return nil, err
	}
	expiresAt := time.Now().Add(AccessTokenDuration).Unix()


	refreshToken, err := s.tokenManager.GenerateRefreshToken(*u)
	if err != nil {
		s.logger.Error("failed to generate refresh token",
			zap.Error(err),
		)
		return nil, err
	}


	// Store tokens
	token := types.Token{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessExpiresAt: expiresAt,
		Type:            "jwt",
	}

	// Add the user
	s.userService.Add(*u) // TODO: write logic that adds user to database

	s.logger.Info("user registered successfully",
		zap.String("username", u.Username),
		zap.String("email", u.Email),
		zap.String("user_id", u.Id.String()),
	)

	return &TokenResponse{
		AccessToken:      accessToken,
		AccessExpiresAt:  expiresAt,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: expiresAt,
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
	accessToken, err := s.tokenManager.GenerateAccessToken(u)
	if err != nil {
		s.logger.Error("failed to generate access token",
			zap.Error(err),
			zap.Uint16("user_id", u.Id),
		)
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(u)
	if err != nil {
		s.logger.Error("failed to generate refresh token",
			zap.Error(err),
			zap.Uint16("user_id", u.Id),
		)
		return nil, err
	}

	s.logger.Info("user logged in successfully",
		zap.String("email", req.Email),
		zap.Uint16("user_id", u.Id),
	)

	accessExpiresAt := time.Now().Add(AccessTokenDuration).Unix()

	// Update user's token
	u.Token = types.Token{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessExpiresAt: accessExpiresAt,
		RefreshExpiresAt:
		Type:            "jwt",
	}

	return &TokenResponse{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessExpiresAt: accessExpiresAt,
	}, nil
}

func (s *service) ValidateToken(token string) (*jwt.Token, error) {
	return s.tokenManager.ValidateToken(token)
}

func (s *service) RefreshToken(refreshToken string) (*TokenResponse, error) {
	s.logger.Info("processing token refresh request")

	// Validate refresh token
	token, err := s.tokenManager.ValidateToken(refreshToken)
	if err != nil {
		s.logger.Warn("invalid refresh token",
			zap.Error(err),
		)
		return nil, err
	}

	// Extract user ID from token
	userID, err := s.tokenManager.ExtractUserID(token)
	if err != nil {
		s.logger.Error("failed to extract user ID from token",
			zap.Error(err),
		)
		return nil, err
	}

	// Get user
	u, found := s.userService.Get(userID)
	if !found {
		s.logger.Warn("user not found during token refresh",
			zap.Uint16("user_id", userID),
		)
		return nil, ErrUserNotFound
	}

	// Generate new access token
	accessToken, err := s.tokenManager.GenerateAccessToken(u)
	if err != nil {
		s.logger.Error("failed to generate new access token",
			zap.Error(err),
			zap.Uint16("user_id", userID),
		)
		return nil, err
	}

	s.logger.Info("token refreshed successfully",
		zap.Uint16("user_id", userID),
	)

	expiresAt := time.Now().Add(AccessTokenDuration).Unix()

	return &TokenResponse{
		AccessToken:     accessToken,
		AccessExpiresAt: expiresAt,
	}, nil
}
