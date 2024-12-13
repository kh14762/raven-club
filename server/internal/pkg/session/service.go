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

	_ "github.com/lib/pq"
)

type Service interface {
	GenerateSessionToken() string
	CreateSession(token string, userId uint32) (*types.Session, error)
	ValidateSessionToken(token string) (*types.Session, *types.User, error)
	InvalidateSession(sessionId string) error
}

type service struct {
	logger     *zap.Logger
	repository *Repository
}

func NewService(logger *zap.Logger, r *Repository) Service {
	return &service{
		logger:     logger,
		repository: r,
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

	_, err := s.repository.db.Exec("INSERT INTO user_session (id, user_id, expires_at) VALUES ($1, $2, $3)",
		sessionID, userId, session.ExpiresAt)
	if err != nil {
		s.logger.Error("Failed to insert session data into db", zap.Error(err))
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
