package v2

import "github.com/gin-gonic/gin"

func GetStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"version": "v2",
		"message": "Hello World",
	})
}
