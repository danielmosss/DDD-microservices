package ingest

import (
	"context"
	"fmt"
	"time"

	"monitoring/internal/domain/models"
)

type MetingRepository interface {
	Save(ctx context.Context, m models.Meting, returnObject bool) (models.Meting, error)
}

type IngestService struct {
	repo MetingRepository
}

func NewIngestService(repo MetingRepository) *IngestService {
	return &IngestService{repo: repo}
}

func (s *IngestService) VerwerkMeting(ctx context.Context, inc models.IncMeting) (models.Meting, error) {
	meting := models.Meting{
		Time:        time.Now().UTC(),
		SensorID:    inc.SensorID,
		KunstwerkID: inc.KunstwerkID,
		Waarde:      inc.Waarde,
		IsHandmatig: false,
	}

	saved, err := s.repo.Save(ctx, meting, false)
	if err != nil {
		return models.Meting{}, fmt.Errorf("repo save faalde: %w", err)
	}

	return saved, nil
}
