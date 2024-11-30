package api

import (
	"fmt"
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
		username, _ := c.Get("username")
		email, _ := c.Get("email")
		password, _ := c.Get("password")

		fmt.Printf("username: %v, email: %v, password: %v\n", username, email, password)

		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"email":    email,
			"password": password,
		})

	}
}

func handleLogin(as auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		password, _ := c.Get("password")

		c.JSON(http.StatusOK, gin.H{
			"user_id":  userID,
			"username": username,
			"password": password,
		})
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
