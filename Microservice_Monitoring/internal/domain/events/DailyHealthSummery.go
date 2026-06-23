package events

import (
	"context"
	"fmt"
	"monitoring/internal/domain/models"
	"monitoring/internal/infra/db"
	"monitoring/server"
	"time"
)

func NewDailyHealthSummary(KunstwerkID int64) (*models.DailyHealthSummary, error) {
	KunstwerkPostgres := db.NewPostgresKunstwerkRepository(server.GetDBPool())
	ctx := context.Background()
	KunstwerkDetail, err := KunstwerkPostgres.GetKunstwerkMetType(ctx, KunstwerkID)
	if err != nil {
		return nil, fmt.Errorf("daily health summary ophalen kunstwerk %d mislukt: %w", KunstwerkID, err)
	}

	AantalSensoren, err := KunstwerkPostgres.GetAantalSensoren(ctx, KunstwerkID)
	if err != nil {
		return nil, fmt.Errorf("daily health summary ophalen aantal sensoren voor kunstwerk %d mislukt: %w", KunstwerkID, err)
	}

	AantalActieveSensoren, err := KunstwerkPostgres.GetAantalActieveSensoren(ctx, KunstwerkID)
	if err != nil {
		return nil, fmt.Errorf("daily health summary ophalen aantal actieve sensoren voor kunstwerk %d mislukt: %w", KunstwerkID, err)
	}

	AantalAfwijkendeSensoren, err := KunstwerkPostgres.GetAantalSensorenMetNAfwijkingen(ctx, KunstwerkID)
	if err != nil {
		return nil, fmt.Errorf("daily health summary ophalen aantal afwijkende sensoren voor kunstwerk %d mislukt: %w", KunstwerkID, err)
	}

	AantalAfwijkingen, err := KunstwerkPostgres.GetAantalAfwijkingen(ctx, KunstwerkID)
	if err != nil {
		return nil, fmt.Errorf("daily health summary ophalen aantal afwijkingen voor kunstwerk %d mislukt: %w", KunstwerkID, err)
	}

	var status = models.StatusHealthy
	if AantalAfwijkendeSensoren == 0 {
		status = models.StatusOffline
	} else {
		var threshold = int(float64(AantalActieveSensoren) * 0.1)
		if threshold > AantalAfwijkendeSensoren {
			status = models.StatusWarning
		} else {
			status = models.StatusCritical
		}
	}

	return &models.DailyHealthSummary{
		KunstwerkID:               KunstwerkID,
		KunstwerkBeheerIdentifier: KunstwerkDetail.Kunstwerk.BeheerIdentifier,
		KunstwerkDetail:           KunstwerkDetail,
		Tijd:                      time.Now(),
		Status:                    status,
		AantalSensoren:            AantalSensoren,
		AantalActieveSensoren:     AantalActieveSensoren,
		AantalAfwijkendeSensoren:  AantalAfwijkendeSensoren,
		AantalAfwijkingen:         AantalAfwijkingen,
	}, nil
}
