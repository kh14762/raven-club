package auth

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"raven-club/internal/pkg/types"
	"time"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidSignMethod = errors.New("unexpected signing method")
	ErrInvalidClaims     = errors.New("invalid token claims")
)

const (
	AccessTokenDuration  = time.Hour
	RefreshTokenDuration = time.Hour * 24 * 30 // 30 days
)

type TokenManager struct {
	secretKey []byte
}

func NewTokenManager(secretKey []byte) *TokenManager {
	return &TokenManager{
		secretKey: secretKey,
	}
}

func (tm *TokenManager) GenerateAccessToken(u types.User) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	expiresAt := time.Now().Add(AccessTokenDuration)
	claims["user_id"] = u.Id
	claims["username"] = u.Username
	claims["exp"] = expiresAt.Unix()

	t, err := token.SignedString(tm.secretKey)
	return t, expiresAt, err
}

func (tm *TokenManager) GenerateRefreshToken(u types.User) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	expiresAt := time.Now().Add(RefreshTokenDuration)
	claims["user_id"] = u.Id
	claims["exp"] = expiresAt.Unix()

	t, err := token.SignedString(tm.secretKey)
	return t, expiresAt, err
}

func (tm *TokenManager) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignMethod
		}
		return tm.secretKey, nil
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

func (tm *TokenManager) ExtractUserID(token *jwt.Token) (uint16, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidClaims
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrInvalidClaims
	}

	return uint16(userID), nil
}
