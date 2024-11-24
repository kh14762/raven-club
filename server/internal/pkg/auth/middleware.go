package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Custom middleware errors
var (
	ErrNoAuthHeader      = "no authorization header"
	ErrInvalidAuthHeader = "invalid authorization header"
)

type Middleware struct {
	authService Service
	logger      *zap.Logger
}

func NewMiddleware(authService Service, logger *zap.Logger) *Middleware {
	return &Middleware{
		authService: authService,
		logger:      logger.With(zap.String("middleware", "auth")),
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

		// Validate token
		token, err := am.authService.ValidateToken(tokenParts[1])
		if err != nil || !token.Valid {
			am.logger.Warn("invalid token",
				zap.Error(err),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			am.logger.Error("failed to extract token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken})
			c.Abort()
			return
		}

		// Set user context for downstream handlers
		userID := uint16(claims["user_id"].(float64))
		username := claims["username"].(string)

		c.Set("user_id", userID)
		c.Set("username", username)

		am.logger.Info("request authenticated",
			zap.Uint16("user_id", userID),
			zap.String("username", username),
		)

		c.Next()
	}
}
