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
