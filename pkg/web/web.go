package web

import (
	"net/http"
	"simplecrud/pkg/handlers"
	"simplecrud/pkg/user"
	"simplecrud/utils"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func StartServer(userRepo user.Repository) {
	userService := user.NewService(userRepo) // pass userRepo to NewService
	userHandler := handlers.NewUserHandler(userService)

	gin.SetMode(utils.GetEnv("GIN_MODE", gin.DebugMode))

	r := gin.Default()

	// Create a limiter struct.
	limiter := tollbooth.NewLimiter(1, nil) // this will limit 1 request/second

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// User routes with rate limiter middleware
	r.GET("/users", tollbooth_gin.LimitHandler(limiter), userHandler.GetAllUsers)
	r.GET("/users/:id", tollbooth_gin.LimitHandler(limiter), userHandler.GetUser)
	r.POST("/users", tollbooth_gin.LimitHandler(limiter), userHandler.CreateUser)
	r.PUT("/users/:id", tollbooth_gin.LimitHandler(limiter), userHandler.UpdateUser)
	r.DELETE("/users/:id", tollbooth_gin.LimitHandler(limiter), userHandler.DeleteUser)

	port := utils.GetEnv("PORT", "8080")
	utils.HandleError("PORT not set", nil)

	err := r.Run(":" + port)
	utils.HandleError("Failed to run server", err)
}
