package v1

import (
	"errors"
	"net/http"
	"strconv"

	"monitoring/internal/domain/models"
	"monitoring/internal/infra/db"
	"monitoring/server"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetSensorTypes(c *gin.Context) {
	repo := db.NewPostgresSensorRepository(server.GetDBPool())
	sensorTypes, err := repo.GetSensorTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sensorTypes)
}

func GetSensorConfiguratieBronnen(c *gin.Context) {
	kunstwerkID, err := strconv.ParseInt(c.Param("kunstwerkId"), 10, 64)
	if err != nil || kunstwerkID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kunstwerkId"})
		return
	}

	var sensorTypeID *int64
	if rawSensorTypeID := c.Query("sensorTypeId"); rawSensorTypeID != "" {
		parsedSensorTypeID, err := strconv.ParseInt(rawSensorTypeID, 10, 64)
		if err != nil || parsedSensorTypeID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensorTypeId"})
			return
		}
		sensorTypeID = &parsedSensorTypeID
	}

	repo := db.NewPostgresSensorRepository(server.GetDBPool())
	bronnen, err := repo.GetConfiguratieBronnen(c.Request.Context(), kunstwerkID, sensorTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bronnen)
}

func CreateSensorForOnderdeel(c *gin.Context) {
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

	var request models.CreateSensorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if request.SensorTypeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sensorTypeId is verplicht"})
		return
	}
	if request.Configuratie == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sensorconfiguratie is verplicht bij aanmaken van een sensor"})
		return
	}

	kunstwerkRepo := db.NewPostgresKunstwerkRepository(server.GetDBPool())
	exists, err := kunstwerkRepo.KunstwerkExists(c.Request.Context(), kunstwerkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "kunstwerk niet gevonden"})
		return
	}

	onderdeelExists, err := kunstwerkRepo.OnderdeelBelongsToKunstwerk(c.Request.Context(), onderdeelID, kunstwerkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !onderdeelExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "onderdeel hoort niet bij dit kunstwerk"})
		return
	}

	sensorRepo := db.NewPostgresSensorRepository(server.GetDBPool())
	sensorType, err := sensorRepo.GetSensorTypeByID(c.Request.Context(), request.SensorTypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "sensortype bestaat niet"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if request.ConfiguratieBronSensorID != nil {
		belongs, err := sensorRepo.SensorBelongsToKunstwerk(c.Request.Context(), *request.ConfiguratieBronSensorID, kunstwerkID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !belongs {
			c.JSON(http.StatusBadRequest, gin.H{"error": "configuratiebron hoort niet bij dit kunstwerk"})
			return
		}
	}

	config := *request.Configuratie
	if validationError := validateConfiguratieForSensorType(sensorType, &config); validationError != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationError})
		return
	}

	sensor, err := sensorRepo.CreateSensorWithConfiguratie(c.Request.Context(), kunstwerkID, onderdeelID, request, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detail, err := sensorRepo.GetSensorAndLastMetingAfwijkingFromSensorId(c.Request.Context(), kunstwerkID, sensor.ID)
	if err != nil {
		c.JSON(http.StatusCreated, sensor)
		return
	}

	c.JSON(http.StatusCreated, detail)
}

func UpdateSensorConfiguratie(c *gin.Context) {
	sensorID, err := strconv.ParseInt(c.Param("sensorId"), 10, 64)
	if err != nil || sensorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensorId"})
		return
	}

	var request models.UpdateSensorConfiguratieRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	sensorRepo := db.NewPostgresSensorRepository(server.GetDBPool())
	config, err := sensorRepo.UpdateSensorConfiguratie(c.Request.Context(), sensorID, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

func DeleteSensor(c *gin.Context) {
	sensorID, err := strconv.ParseInt(c.Param("sensorId"), 10, 64)
	if err != nil || sensorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensorId"})
		return
	}

	var request models.DeleteDataRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	repo := db.NewPostgresSensorRepository(server.GetDBPool())
	if err := repo.DeleteSensor(c.Request.Context(), sensorID, request.PreserveSensorData); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "sensor niet gevonden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true, "preserveSensorData": request.PreserveSensorData})
}

func validateConfiguratieForSensorType(sensorType models.SensorType, config *models.UpdateSensorConfiguratieRequest) string {
	if config.MargePercentage != nil && *config.MargePercentage < 0 {
		return "margePercentage mag niet negatief zijn"
	}

	if sensorType.DrempelIsRange {
		if config.MinValue == nil || config.MaxValue == nil {
			return "minValue en maxValue zijn verplicht voor range-sensoren"
		}
		if *config.MinValue >= *config.MaxValue {
			return "minValue moet lager zijn dan maxValue"
		}
		return ""
	}

	if config.MinValue == nil {
		return "minValue is verplicht als normwaarde voor niet-range sensoren"
	}
	config.MaxValue = nil
	return ""
}