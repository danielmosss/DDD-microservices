package events

import (
	"context"
	"fmt"
	"monitoring/internal/domain/models"
	"monitoring/internal/infra/db"
	"monitoring/server"
	"time"
)

type Status string

const (
	StatusHealthy  Status = "healthy"
	StatusWarning  Status = "warning"
	StatusCritical Status = "critical"
	StatusOffline  Status = "offline"
)

type DailyHealthSummery struct {
	KunstwerkID               int64                  `json:"kunstwerkId"`
	KunstwerkBeheerIdentifier string                 `json:"kunstwerkBeheerIdentifier"`
	KunstwerkDetail           models.KunstwerkDetail `json:"kunstwerkDetail"`
	Tijd                      time.Time              `json:"tijd"`
	Status                    Status                 `json:"status"`
	AantalActieveSensoren     int                    `json:"aantalActieveSensoren"`
	AantalAfwijkendeSensoren  int                    `json:"aantalAfwijkendeSensoren"`
	AantalAfwijkingen         int                    `json:"aantalAfwijkingen"`
}

func NewDailyHealthSummery(KunstwerkID int64) (*DailyHealthSummery, error) {
	KunstwerkPostgres := db.NewPostgresKunstwerkRepository(server.GetDBPool())
	ctx := context.Background()
	KunstwerkDetail, err := KunstwerkPostgres.GetKunstwerkMetType(ctx, KunstwerkID)
	if err != nil {
		return nil, fmt.Errorf("daily health summary ophalen kunstwerk %d mislukt: %w", KunstwerkID, err)
	}

	AantalActieveSensoren, err := KunstwerkPostgres.GetAantalActieveSensoren(ctx, KunstwerkID)
	if err != nil {
		return nil, fmt.Errorf("daily health summary ophalen aantal actieve sensoren voor kunstwerk %d mislukt: %w", KunstwerkID, err)
	}

	AantalAfwijkendeSensoren, err := KunstwerkPostgres.GetAantalSensorenMetNAfwijkingen(ctx, KunstwerkID, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("daily health summary ophalen aantal afwijkende sensoren voor kunstwerk %d mislukt: %w", KunstwerkID, err)
	}

	AantalAfwijkingen, err := KunstwerkPostgres.GetAantalAfwijkingen(ctx, KunstwerkID, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("daily health summary ophalen aantal afwijkingen voor kunstwerk %d mislukt: %w", KunstwerkID, err)
	}

	return &DailyHealthSummery{
		KunstwerkID:               KunstwerkID,
		KunstwerkBeheerIdentifier: KunstwerkDetail.Kunstwerk.BeheerIdentifier,
		KunstwerkDetail:           KunstwerkDetail,
		Tijd:                      time.Now(),
		Status:                    StatusHealthy,
		AantalActieveSensoren:     AantalActieveSensoren,
		AantalAfwijkendeSensoren:  AantalAfwijkendeSensoren,
		AantalAfwijkingen:         AantalAfwijkingen,
	}, nil
}
