package ingest

import (
	"context"
	"fmt"
	"time"

	"monitoring/internal/domain/models"

	"github.com/google/uuid"
)

// 1. Interface: Zo weet de service wel wát hij moet opslaan, maar niet hóe (Postgres/Timescale)
type MetingRepository interface {
	Save(ctx context.Context, m models.Meting) error
}

type IngestService struct {
	repo MetingRepository
}

// Constructor
func NewIngestService(repo MetingRepository) *IngestService {
	return &IngestService{repo: repo}
}

// 2. De Verwerkingslogica
func (s *IngestService) VerwerkMeting(ctx context.Context, inc models.IncMeting) (models.Meting, error) {
	kwID, err := uuid.Parse(inc.KunstwerkID)
	if err != nil {
		return models.Meting{}, fmt.Errorf("ongeldig kunstwerk_id (%s): %w", inc.KunstwerkID, err)
	}

	var sensorIDPtr *uuid.UUID
	if inc.SensorID != "" {
		sID, err := uuid.Parse(inc.SensorID)
		if err != nil {
			return models.Meting{}, fmt.Errorf("ongeldig sensor_id (%s): %w", inc.SensorID, err)
		}
		sensorIDPtr = &sID
	}

	meting := models.Meting{
		Time:        time.Now(),
		SensorID:    sensorIDPtr,
		KunstwerkID: kwID,
		Waarde:      inc.Waarde,
		IsAfwijking: false,
		IsHandmatig: false,
		InspectieID: nil,
	}

	if err := s.repo.Save(ctx, meting); err != nil {
		return models.Meting{}, fmt.Errorf("repo save faalde: %w", err)
	}

	return meting, nil
}
