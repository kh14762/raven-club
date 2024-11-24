package auth

import (
	"errors"
	"github.com/golang-jwt/jwt"
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
	Register(u types.User) (*TokenResponse, error)
	Login(req LoginRequest) (*TokenResponse, error)
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

func (s *service) Register(u types.User) (*TokenResponse, error) {
	// Add request logging
	s.logger.Info("processing registration request",
		zap.String("username", u.Username),
		zap.String("email", u.Email),
	)

	// Validate and normalize email
	if err := s.emailValidator.Validate(u.Email); err != nil {
		s.logger.Warn("invalid email format",
			zap.String("email", u.Email),
			zap.Error(err),
		)
		return nil, err
	}
	u.Email = s.emailValidator.Normalize(u.Email)

	// Check if email exists
	if err := s.emailLookup.CheckEmailExists(u.Email); err != nil {
		s.logger.Info("email already exists",
			zap.String("email", u.Email),
		)
		return nil, err
	}

	// Check required fields
	if u.Username == "" || u.Email == "" || u.Password == "" {
		s.logger.Warn("missing required fields",
			zap.String("username", u.Username),
			zap.String("email", u.Email),
		)
		return nil, ErrMissingFields
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password",
			zap.Error(err),
		)
		return nil, ErrPasswordProcess
	}
	u.Password = string(hashedPassword)

	// Generate tokens using TokenManager
	accessToken, err := s.tokenManager.GenerateAccessToken(u)
	if err != nil {
		s.logger.Error("failed to generate access token",
			zap.Error(err),
		)
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(u)
	if err != nil {
		s.logger.Error("failed to generate refresh token",
			zap.Error(err),
		)
		return nil, err
	}

	expiresAt := time.Now().Add(AccessTokenDuration).Unix()

	// Store tokens
	u.Token = types.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Type:         "jwt",
	}

	// Add the user
	s.userService.Add(u)

	s.logger.Info("user registered successfully",
		zap.String("username", u.Username),
		zap.String("email", u.Email),
		zap.Uint16("user_id", u.Id),
	)

	return &TokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *service) Login(req LoginRequest) (*TokenResponse, error) {
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

	expiresAt := time.Now().Add(AccessTokenDuration).Unix()

	// Update user's token
	u.Token = types.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Type:         "jwt",
	}

	return &TokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
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
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}
