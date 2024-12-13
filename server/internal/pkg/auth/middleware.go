package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"net/http"
	"raven-club/internal/pkg/token"
	"strings"
)

// Custom middleware errors
var (
	ErrNoAuthHeader      = "no authorization header"
	ErrInvalidAuthHeader = "invalid authorization header"
	ErrInvalidToken      = errors.New("invalid token")
)

type Middleware struct {
	logger       *zap.Logger
	authService  Service
	tokenService token.Service
}

func NewMiddleware(logger *zap.Logger, authService Service, tokenService token.Service) *Middleware {
	return &Middleware{
		logger:       logger.With(zap.String("middleware", "auth")),
		authService:  authService,
		tokenService: tokenService,
	}
}

func (am *Middleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request
		am.logger.Info("processing request",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		// Get authorization header
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

		// Validate tkn
		tkn, err := am.tokenService.ValidateJwt(tokenParts[1])
		if err != nil || !tkn.Valid {
			am.logger.Warn("invalid tkn",
				zap.Error(err),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := tkn.Claims.(jwt.MapClaims)
		if !ok {
			am.logger.Error("failed to extract tkn claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken})
			c.Abort()
			return
		}

		// Set user context for downstream handlers
		userID := uint16(claims["id"].(float64))
		username := claims["username"].(string)

		c.Set("id", userID)
		c.Set("username", username)

		am.logger.Info("request authenticated",
			zap.Uint16("id", userID),
			zap.String("username", username),
		)

		c.Next()
	}
}
