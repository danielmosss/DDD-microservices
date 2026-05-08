package events

import (
	"monitoring/internal/domain/models"

	"github.com/google/uuid"
)

type AfwijkingGedetecteerd struct {
	ID             int64     `json:"id"`
	MetingID       uuid.UUID `json:"metingId"`
	MetingTime     string    `json:"metingTime"`
	KunstwerkID    int64     `json:"kunstwerkId"`
	SensorID       int64     `json:"sensorId"`
	AfwijkingsTime string    `json:"afwijkingsTime"`
	NormMinWaarde  float64   `json:"normMinWaarde"`
	NormMaxWaarde  float64   `json:"normMaxWaarde"`
	GemetenWaarde  float64   `json:"gemetenWaarde"`
	IsWarning      bool      `json:"isWarning"`
}

func NieuweAfwijkingGedetecteerd(internAfwijking models.Afwijking) AfwijkingGedetecteerd {
	return AfwijkingGedetecteerd{
		ID:             internAfwijking.ID,
		MetingID:       internAfwijking.MetingID,
		MetingTime:     internAfwijking.MetingTime.Format("DD-MM-YYYY HH:mm:ss"),
		KunstwerkID:    internAfwijking.KunstwerkID,
		SensorID:       internAfwijking.SensorID,
		AfwijkingsTime: internAfwijking.Time.Format("DD-MM-YYYY HH:mm:ss"),
		NormMinWaarde:  internAfwijking.NormMinWaarde,
		NormMaxWaarde:  internAfwijking.NormMaxWaarde,
		GemetenWaarde:  internAfwijking.GemetenWaarde,
		IsWarning:      internAfwijking.IsWarning,
	}
}
