package web

import (
	"net/http"

	"simplecrud/utils"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	gin.SetMode(utils.GetEnv("GIN_MODE", gin.DebugMode))

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	port := utils.GetEnv("PORT", "8080")
	utils.HandleError("PORT not set", nil)

	err := r.Run(":" + port)
	utils.HandleError("Failed to run server", err)
}
