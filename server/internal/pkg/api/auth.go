package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
	"raven-club/internal/pkg/auth"
)

var AuthController = fx.Module("AuthController",
	fx.Invoke(InitAuthController),
	fx.Invoke(InitProtectedRoutes), // Add protected routes initialization
)

// InitAuthController Public routes
func InitAuthController(engine *gin.Engine, as auth.Service) {
	engine.POST("/auth/register", handleRegister(as))
	engine.POST("/auth/login", handleLogin(as))
}

// InitProtectedRoutes Protected routes
func InitProtectedRoutes(engine *gin.Engine, as auth.Service, am *auth.Middleware) {
	// Create protected group
	protected := engine.Group("/api")
	protected.Use(am.Handler())
	{
		protected.GET("/profile", handleProfile())
		protected.POST("/refresh-token", handleRefreshToken(as))
		// Add more protected routes here
	}
}

// Handler functions
func handleRegister(as auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Your registration logic
	}
}

func handleLogin(as auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Your login logic
	}
}

func handleProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user info from context
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		c.JSON(http.StatusOK, gin.H{
			"user_id":  userID,
			"username": username,
		})
	}
}

func handleRefreshToken(as auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Your refresh token logic
	}
}
