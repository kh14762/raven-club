package token

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"raven-club/internal/pkg/types"
	"raven-club/internal/pkg/user"
	"time"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidSignMethod = errors.New("unexpected signing method")
	ErrInvalidClaims     = errors.New("invalid token claims")
	ErrUserNotFound      = errors.New("user not found")
)

type service struct {
	logger      *zap.Logger
	config      *Config
	userService user.Service
}

func NewService(logger *zap.Logger, cfg *Config, us user.Service) Service {
	return &service{
		logger:      logger.With(zap.String("service", "token")),
		config:      cfg,
		userService: us,
	}
}

func (s *service) GenerateAccessToken(u types.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	expiresAt := time.Now().Add(s.config.AccessTokenDuration)
	claims["id"] = u.Id
	claims["username"] = u.Username
	claims["exp"] = expiresAt.Unix()

	t, err := token.SignedString(s.config.SecretKey)
	return t, err
}

func (s *service) GenerateRefreshToken(u types.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	expiresAt := time.Now().Add(s.config.RefreshTokenDuration)
	claims["id"] = u.Id
	claims["exp"] = expiresAt.Unix()

	t, err := token.SignedString(s.config.SecretKey)
	return t, err
}

func (s *service) RefreshAccessToken(refreshToken string) (string, error) {
	s.logger.Info("processing token refresh request")

	// Validate refresh token
	token, err := s.ValidateToken(refreshToken)
	if err != nil {
		s.logger.Warn("invalid refresh token",
			zap.Error(err),
		)
		return "", err
	}

	// Extract user ID from token
	userID, err := s.extractUserID(token)
	if err != nil {
		s.logger.Error("failed to extract user ID from token",
			zap.Error(err),
		)
		return "", err
	}

	// Get user
	u, found := s.userService.Get(userID)
	if !found {
		s.logger.Warn("user not found during token refresh",
			zap.Uint32("id", userID),
		)
		return "", ErrUserNotFound
	}

	// Generate new access token
	accessToken, err := s.GenerateAccessToken(u)
	if err != nil {
		s.logger.Error("failed to generate new access token",
			zap.Error(err),
			zap.Uint32("id", userID),
		)
		return "", err
	}

	s.logger.Info("token refreshed successfully",
		zap.Uint32("id", userID),
	)

	return accessToken, nil
}

func (s *service) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignMethod
		}
		return token, nil
	})

	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrInvalidToken
	}

	return token, nil
}

func (s *service) extractUserID(token *jwt.Token) (uint32, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidClaims
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return 0, ErrInvalidClaims
	}

	return uint32(userID), nil
}
