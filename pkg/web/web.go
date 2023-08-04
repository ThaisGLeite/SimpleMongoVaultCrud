package web

import (
	"net/http"

	"simplecrud/pkg/handlers"
	"simplecrud/utils"

	"simplecrud/pkg/user"

	"github.com/gin-gonic/gin"
)

func StartServer(userRepo user.Repository) {
	userService := user.NewService(userRepo) // pass userRepo to NewService
	userHandler := handlers.NewUserHandler(userService)

	gin.SetMode(utils.GetEnv("GIN_MODE", gin.DebugMode))

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// User routes
	r.GET("/users", userHandler.GetAllUsers)
	r.GET("/users/:id", userHandler.GetUser)
	r.POST("/users", userHandler.CreateUser)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	port := utils.GetEnv("PORT", "8080")
	utils.HandleError("PORT not set", nil)

	err := r.Run(":" + port)
	utils.HandleError("Failed to run server", err)
}
