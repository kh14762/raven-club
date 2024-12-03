package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
	"raven-club/internal/pkg/auth"
	"raven-club/internal/pkg/types"
)

var AuthController = fx.Module("AuthController",
	fx.Invoke(InitAuthController),
	fx.Invoke(InitProtectedRoutes), // Add protected routes initialization
)

// InitAuthController Public routes
func InitAuthController(engine *gin.Engine, as auth.Service) {
	engine.POST("/api/auth/register", handleRegister(as))
	engine.POST("/api/auth/login", handleLogin(as))
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
		var req types.RegisterRequest
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// offer to auth service

		c.JSON(http.StatusOK, gin.H{
			"message": "Registration successful", // TODO generate and pass JWT back to client
		})
	}
}

func handleLogin(as auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO handle login
	}
}

func handleProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user info from context
	}
}

func handleRefreshToken(as auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Your refresh token logic
	}
}
