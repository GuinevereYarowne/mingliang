package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()
	r.GET("/api/v1/uuid", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"uuid": uuid.New().String(),
		})
	})
	r.Run(":8080")
}