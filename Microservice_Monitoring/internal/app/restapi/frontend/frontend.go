package frontend

import (
	"context"
	"fmt"
	"monitoring/internal/domain/models"
	"monitoring/internal/infra/db"
	"monitoring/server"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetKunstwerkTree
// @Summary Returns a hierarchical tree of the kunstwerk, its onderdelen, and their sensoren
// @Description Tree of the kunstwerk
// @Tags Frontend
// @Accept json
// @Produce json
// @Param kunstwerkId path int true "Kunstwerk ID"
// @Success 200 {array} models.KunstwerkTreeResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/frontend/kunstwerken/{kunstwerkId}/tree [get]
func GetKunstwerkTree(c *gin.Context) {
	idParam := c.Param("kunstwerkId")
	var KunstwerkId int64
	KunstwerkId, err := strconv.ParseInt(idParam, 10, 64)

	KunstwerkPostgres := db.NewPostgresKunstwerkRepository(server.GetDBPool())
	ctx := context.Background()

	KunstwerkDetails, err := KunstwerkPostgres.GetKunstwerkMetType(ctx, KunstwerkId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Fout bij ophalen kunstwerk: %v", err)})
		return
	}

	KunstwerkOnderdelenMetSensors, err := KunstwerkPostgres.GetKunstwerkOnderdelenWithSensorIDs(ctx, KunstwerkId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Fout bij ophalen onderdelen: %v", err)})
		return
	}

	response := models.KunstwerkTreeResponse{
		Kunstwerk:  KunstwerkDetails,
		Onderdelen: []*models.TreeOnderdeel{},
	}

	var onderdelenMap = make(map[int64]*models.TreeOnderdeel)

	for _, onderdeel := range KunstwerkOnderdelenMetSensors {
		ParentId := onderdeel.ParentId
		if ParentId == nil {
			ParentId = new(int64)
		}
		onderdelenMap[onderdeel.ID] = &models.TreeOnderdeel{
			ID:         onderdeel.ID,
			Naam:       onderdeel.Naam,
			ParentId:   *ParentId,
			Deleted:    onderdeel.Deleted,
			Sensoren:   onderdeel.SensorIds,
			Onderdelen: nil,
		}
	}

	for _, onderdeel := range KunstwerkOnderdelenMetSensors {
		onderdeelInMap := onderdelenMap[onderdeel.ID]
		if onderdeelInMap.ParentId != 0 {
			parentOnderdeel, exists := onderdelenMap[onderdeelInMap.ParentId]
			if exists {
				parentOnderdeel.Onderdelen = append(parentOnderdeel.Onderdelen, onderdeelInMap)
			} else {
				response.Onderdelen = append(response.Onderdelen, onderdeelInMap)
			}
		} else {
			// This is the parent
			response.Onderdelen = append(response.Onderdelen, onderdeelInMap)
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetBulkSensorData
// @Summary Returns for each sensor that is given the latest meting (and afwijking if it is) and the sensor configuration.
// @Description
// @Tags Frontend
// @Accept json
// @Produce json
// @Param kunstwerkId path int true "Kunstwerk ID"
// @Param request body BulkSensorDataRequest true "JSON body met een array van sensor IDs"
// @Success 200 {array} models.SensorDetailResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/frontend/kunstwerken/{kunstwerkId}/sensoren/bulk-actueel [post]
func GetBulkSensorData(c *gin.Context) {
	idParam := c.Param("kunstwerkId")
	var KunstwerkId int64
	KunstwerkId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kunstwerkId parameter"})
		return
	}

	var reqBody BulkSensorDataRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	SensorPostgres := db.NewPostgresSensorRepository(server.GetDBPool())
	ctx := context.Background()

	var response []models.SensorDetailResponse
	var errorOccurred bool
	var errorMessages []string

	for _, sensorId := range reqBody.SensorIds {
		SensorDetailRes, err := SensorPostgres.GetSensorAndLastMetingAfwijkingFromSensorId(ctx, KunstwerkId, sensorId)
		if err != nil {
			errorOccurred = true
			errorMessages = append(errorMessages, err.Error())
		}
		response = append(response, SensorDetailRes)
	}

	if errorOccurred {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessages})
		return
	}

	c.JSON(http.StatusOK, response)
}

type BulkSensorDataRequest struct {
	SensorIds []int64 `json:"sensorIds"`
}

//// 2. GET /onderdelen/:id/sensoren
//func GetSensorenByOnderdeel(c *gin.Context) {
//	onderdeelId := c.Param("id")
//
//	// 1. Haal alle sensoren op die gekoppeld zijn aan dit specifieke onderdeelId
//	// var sensoren []models.Sensor
//	// db.Select(&sensoren, "SELECT * FROM sensoren WHERE onderdeel_id = ?", onderdeelId)
//
//	// 2. Map dit naar je uitgebreide DTO
//	var responseList []models.SensorDetailResponse
//
//	// Voorbeeld loop:
//	// for _, s := range sensoren {
//	//    config := fetchConfig(s.ID)
//	//    laatsteMeting := fetchLastMeting(s.ID)
//	//    status := berekenStatus(laatsteMeting, config)
//	//
//	//    responseList = append(responseList, SensorDetailResponse{
//	//        ID: s.ID,
//	//        SensorTypeID: s.SensorTypeID,
//	//        SensorConfiguratie: config,
//	//        LaatsteMetingWaarde: laatsteMeting,
//	//        Status: status,
//	//    })
//	// }
//
//	// Zelfs als het leeg is, stuur een lege array terug (geen null), dat vindt Angular fijner
//	if responseList == nil {
//		responseList = []models.SensorDetailResponse{}
//	}
//
//	c.JSON(http.StatusOK, responseList)
//}
//
//// 3. GET /sensoren/:id/metingen?range=24h
//func GetMetingenForSensor(c *gin.Context) {
//	sensorId := c.Param("id")
//
//	// Haal de query parameter op, met een fallback als de frontend niets meestuurt
//	timeRange := c.DefaultQuery("range", "24h")
//
//	// Bepaal de starttijd voor je SQL query
//	var startTime time.Time
//	now := time.Now()
//
//	switch timeRange {
//	case "1h":
//		startTime = now.Add(-1 * time.Hour)
//	case "24h":
//		startTime = now.Add(-24 * time.Hour)
//	case "7d":
//		startTime = now.Add(-7 * 24 * time.Hour)
//	default:
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid range parameter. Use 1h, 24h, or 7d"})
//		return
//	}
//
//	// 1. Haal alle metingen op uit de DB voor deze sensor vanaf startTime
//	// var metingen []models.Meting
//	// db.Select(&metingen, "SELECT * FROM metingen WHERE sensor_id = ? AND time >= ? ORDER BY time ASC", sensorId, startTime)
//
//	// mock output
//	var metingen []models.Meting
//
//	if metingen == nil {
//		metingen = []models.Meting{}
//	}
//
//	c.JSON(http.StatusOK, metingen)
