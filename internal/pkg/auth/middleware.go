package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Custom middleware errors
var (
	ErrNoAuthHeader        = "no authorization header"
	ErrInvalidAuthHeader   = "invalid authorization header"
	ErrInvalidSessionToken = errors.New("invalid session token")
)

type Middleware struct {
	logger      *zap.Logger
	authService Service
}

func NewMiddleware(logger *zap.Logger, authService Service) *Middleware {
	return &Middleware{
		logger:      logger.With(zap.String("middleware", "auth")),
		authService: authService,
	}
}

func (am *Middleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request
		am.logger.Info("processing request",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		// GetUserByID authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			am.logger.Warn("missing auth header",
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrNoAuthHeader})
			c.Abort()
			return
		}

		// Validate header format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			am.logger.Warn("invalid auth header format",
				zap.String("auth_header", authHeader),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidAuthHeader})
			c.Abort()
			return
		}

		// TODO: validate session

		c.Next()
	}
}
