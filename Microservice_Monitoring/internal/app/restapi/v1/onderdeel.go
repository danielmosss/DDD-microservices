package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"monitoring/internal/domain/models"
	"monitoring/internal/infra/db"
	"monitoring/server"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func CreateOnderdeel(c *gin.Context) {
	kunstwerkID, err := strconv.ParseInt(c.Param("kunstwerkId"), 10, 64)
	if err != nil || kunstwerkID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kunstwerkId"})
		return
	}

	var request models.CreateOnderdeelRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	request.Naam = strings.TrimSpace(request.Naam)
	if request.Naam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "naam is verplicht"})
		return
	}
	if len(request.Naam) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "naam mag maximaal 255 tekens bevatten"})
		return
	}

	repo := db.NewPostgresKunstwerkRepository(server.GetDBPool())
	exists, err := repo.KunstwerkExists(c.Request.Context(), kunstwerkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "kunstwerk niet gevonden"})
		return
	}

	onderdeel, err := repo.CreateOnderdeel(c.Request.Context(), kunstwerkID, request)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parent onderdeel hoort niet bij dit kunstwerk"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, onderdeel)
}

func DeleteOnderdeel(c *gin.Context) {
	kunstwerkID, err := strconv.ParseInt(c.Param("kunstwerkId"), 10, 64)
	if err != nil || kunstwerkID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kunstwerkId"})
		return
	}

	onderdeelID, err := strconv.ParseInt(c.Param("onderdeelId"), 10, 64)
	if err != nil || onderdeelID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid onderdeelId"})
		return
	}

	var request models.DeleteDataRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	repo := db.NewPostgresKunstwerkRepository(server.GetDBPool())
	if err := repo.DeleteOnderdeelTree(c.Request.Context(), kunstwerkID, onderdeelID, request.PreserveSensorData); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "onderdeel niet gevonden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true, "preserveSensorData": request.PreserveSensorData})
}