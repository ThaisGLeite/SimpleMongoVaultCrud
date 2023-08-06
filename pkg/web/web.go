package web

import (
	"net/http"

	"simplecrud/pkg/handlers"
	"simplecrud/pkg/user"
	"simplecrud/utils"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

// Define constants for logging levels, Gin mode, and server port.
const (
	InfoLogLevel   = "I"
	ErrorLogLevel  = "E"
	GinModeKey     = "GIN_MODE"
	DefaultGinMode = gin.DebugMode
	PortKey        = "PORT"
	DefaultPort    = "8080"
)

// StartServer function initializes and starts the web server.
func StartServer(userRepo user.Repository) {
	// Create a new user service with the provided repository.
	userService := user.NewService(userRepo)
	// Create a new user handler with the created user service.
	userHandler := handlers.NewUserHandler(userService)

	// Set the gin mode. This can be either debug or release.
	gin.SetMode(utils.GetEnv(GinModeKey, DefaultGinMode))

	// Create a new gin engine. Use gin.New() to have more control over the middleware.
	r := gin.New()

	// Use the Recovery middleware to recover from any panics and write a 500 if it happens.
	r.Use(gin.Recovery())

	// Secure headers middleware
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'")
		// Other headers as needed
		c.Next()
	})

	// Create a new rate limiter. This will limit to 1 request/second.
	limiter := tollbooth.NewLimiter(1, nil)

	// Setup the routes for the server.
	setupRoutes(r, limiter, userHandler)

	// Get the port from environment variables or use the default port.
	port := utils.GetEnv(PortKey, DefaultPort)

	// Run the server on the defined port. If there is an error, log it.
	if err := r.Run(":" + port); err != nil {
		utils.HandleError(ErrorLogLevel, "Failed to run server", err)
	}
}

// setupRoutes function sets up all the routes for the server.
func setupRoutes(router *gin.Engine, limiter *limiter.Limiter, userHandler *handlers.UserHandler) {
	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// User routes. These routes are wrapped with a rate limiter middleware.
	router.GET("/users", tollbooth_gin.LimitHandler(limiter), userHandler.GetAllUsers)       // Get all users
	router.GET("/users/:id", tollbooth_gin.LimitHandler(limiter), userHandler.GetUser)       // Get a single user by ID
	router.POST("/users", tollbooth_gin.LimitHandler(limiter), userHandler.CreateUser)       // Create a new user
	router.PUT("/users/:id", tollbooth_gin.LimitHandler(limiter), userHandler.UpdateUser)    // Update a user by ID
	router.DELETE("/users/:id", tollbooth_gin.LimitHandler(limiter), userHandler.DeleteUser) // Delete a user by ID
}
