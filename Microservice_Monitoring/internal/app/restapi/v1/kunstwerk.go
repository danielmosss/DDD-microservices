package v1

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"monitoring/internal/domain/models"
	"monitoring/internal/infra/db"
	"monitoring/server"

	"github.com/gin-gonic/gin"
)

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

// GetMetingen
// @Tags Metingen
// @Accept json
// @Produce json
// @Param kunstwerkId path int true "Kunstwerk ID"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset of number of results"
// @Success 200 {array} PaginatedResponse[models.Meting]
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/metingen/{kunstwerkId} [get]
func GetMetingen(c *gin.Context) {
	idStr := c.Param("kunstwerkId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kunstwerkId"})
		return
	}

	pagination, err := ParsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo := db.NewPostgresMetingRepository(server.GetDBPool())
	metingen, total, err := repo.GetByKunstwerkID(c.Request.Context(), id, pagination.Limit, pagination.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse[models.Meting]{
		Data: metingen,
		Pagination: PaginationMeta{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
			Total:  total,
		},
	})
}

// GetMetingenRecent
// @Tags Metingen
// @Accept json
// @Produce json
// @Param kunstwerkId path int true "Kunstwerk ID"
// @Success 200 {array} models.Meting
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/metingen/{kunstwerkId}/recent [get]
func GetMetingenRecent(c *gin.Context) {
	idStr := c.Param("kunstwerkId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kunstwerkId"})
		return
	}

	repo := db.NewPostgresMetingRepository(server.GetDBPool())
	metingen, err := repo.GetRecentPerSensorByKunstwerkID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metingen)
}

// GetAfwijkingen
// @Tags Afwijkingen
// @Accept json
// @Produce json
// @Param kunstwerkId path int true "Kunstwerk ID"
// @Success 200 {array} PaginatedResponse[models.Afwijking]
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/afwijkingen/{kunstwerkId} [get]
func GetAfwijkingen(c *gin.Context) {
	idStr := c.Param("kunstwerkId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kunstwerkId"})
		return
	}

	pagination, err := ParsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo := db.NewPostgresAfwijkingRepository(server.GetDBPool())
	afwijkingen, total, err := repo.GetByKunstwerkID(c.Request.Context(), id, pagination.Limit, pagination.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse[models.Afwijking]{
		Data: afwijkingen,
		Pagination: PaginationMeta{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
			Total:  total,
		},
	})
}

// PostMeting
// @Tags Metingen
// @Accept json
// @Produce json
// @Param meting body models.IncMeting true "IncMeting"
// @Success 201 {object} models.Meting
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/meting [post]
// Handmatige meting invoer door inspecteur. Vereist sensorId.
func PostMeting(c *gin.Context) {
	var inc models.IncMeting
	if err := c.ShouldBindJSON(&inc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if inc.SensorID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sensorId is required for manual meting"})
		return
	}

	meting := models.Meting{
		Time:        time.Now().UTC(),
		SensorID:    inc.SensorID,
		KunstwerkID: inc.KunstwerkID,
		Waarde:      inc.Waarde,
		IsHandmatig: true,
	}

	repo := db.NewPostgresMetingRepository(server.GetDBPool())
	saved, err := repo.Save(context.Background(), meting, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, saved)
}

// GetKunstwerken
// @Tags Kunstwerken
// @Accept json
// @Produce json
// @Success 200 {array} models.Kunstwerk
// @Failure 500 {object} map[string]string
// @Router /api/v1/kunstwerken [get]
func GetKunstwerken(c *gin.Context) {
	repo := db.NewPostgresKunstwerkRepository(server.GetDBPool())
	kunstwerken, err := repo.GetActieveKunstwerken(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kunstwerken)
}

// GetSensorenByKunstwerk
// @Tags Sensoren
// @Accept json
// @Produce json
// @Param kunstwerkId path int true "Kunstwerk ID"
// @Success 200 {array} models.Sensor
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/kunstwerken/{kunstwerkId}/sensoren [get]
func GetSensorenByKunstwerk(c *gin.Context) {
	idStr := c.Param("kunstwerkId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kunstwerkId"})
		return
	}

	repo := db.NewPostgresKunstwerkRepository(server.GetDBPool())
	sensoren, err := repo.GetSensorenByKunstwerkID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sensoren)
}
