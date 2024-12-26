package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"go.uber.org/zap"
	"raven-club/internal/pkg/types"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

type Service interface {
	GenerateSessionToken() string
	CreateSession(token string, userId uint32) (*types.Session, error)
	ValidateSessionToken(token string) (*types.Session, *types.User, error)
	InvalidateSession(sessionId string) error
}

type service struct {
	logger *zap.Logger
	store  *sessions.CookieStore
}

func NewService(logger *zap.Logger) Service {
	return &service{
		logger: logger,
		store:  sessions.NewCookieStore([]byte("mx1HxnLBtXDWTgJaI13zAN3tG8lQS6jJnwreIHJhXa7/OqP8U62TZfST4Ao/xasZ")),
	}
}

func (s *service) GenerateSessionToken() string {
	var bytes = make([]byte, 20)
	_, err := rand.Read(bytes) // TODO use UUID here
	if err != nil {
		s.logger.Error("Failed to generate session token", zap.Error(err))
	}

	token := strings.ToLower(base32.StdEncoding.EncodeToString(bytes))
	token = strings.TrimRight(token, "=")

	return token
}

func (s *service) CreateSession(token string, userId uint32) (*types.Session, error) {
	sessionID := encodeHexLowerCase(sha256.Sum256([]byte(token)))

	session := &types.Session{
		ID:        sessionID,
		UserID:    userId,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}

	return session, nil
}

func encodeHexLowerCase(bytes [32]byte) string {
	return strings.ToLower(hex.EncodeToString(bytes[:]))
}

func (s *service) ValidateSessionToken(token string) (*types.Session, *types.User, error) {
	// TODO
	return &types.Session{}, &types.User{}, nil
}

func (s *service) InvalidateSession(sessionId string) error {
	// TODO
	return nil
}
