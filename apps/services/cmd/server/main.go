package main

import (
	"github.com/HarshithRajesh/clarity/internal/config"
	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Health check working"})
}

func main() {
	config.ConnectDB()

	r := gin.Default()
	r.GET("/health", health)
	r.Run(":8080")
}
