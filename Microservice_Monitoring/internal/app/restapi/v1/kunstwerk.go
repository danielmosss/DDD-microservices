package v1

import "github.com/gin-gonic/gin"

// Get Status
// @Summary Get Status
// @Description Get the status of the API
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/status [get]
func GetStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
