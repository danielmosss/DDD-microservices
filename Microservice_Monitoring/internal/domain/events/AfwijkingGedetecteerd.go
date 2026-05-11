package events

import (
	"time"
	"monitoring/internal/domain/models"
)

type AfwijkingGedetecteerd struct {
	ID                  int64      `json:"id"`
	MetingID            int64      `json:"metingId"`
	MetingTime          time.Time  `json:"metingTime"`
	KunstwerkID         int64      `json:"kunstwerkId"`
	SensorID            *int64     `json:"sensorId"`
	AfwijkingsTime      time.Time  `json:"afwijkingsTime"`
	NormMinWaarde       float64    `json:"normMinWaarde"`
	NormMaxWaarde       float64    `json:"normMaxWaarde"`
	NormMargePercentage *float64   `json:"normMargePercentage"`
	GemetenWaarde       float64    `json:"gemetenWaarde"`
	IsWarning           bool       `json:"isWarning"`
}

func NieuweAfwijkingGedetecteerd(internAfwijking models.Afwijking) AfwijkingGedetecteerd {
	return AfwijkingGedetecteerd{
		ID:                  internAfwijking.ID,
		MetingID:            internAfwijking.MetingID,
		MetingTime:          internAfwijking.MetingTime,
		KunstwerkID:         internAfwijking.KunstwerkID,
		SensorID:            internAfwijking.SensorID,
		AfwijkingsTime:      internAfwijking.Time,
		NormMinWaarde:       internAfwijking.NormMinWaarde,
		NormMaxWaarde:       internAfwijking.NormMaxWaarde,
		NormMargePercentage: internAfwijking.NormMargePercentage,
		GemetenWaarde:       internAfwijking.GemetenWaarde,
		IsWarning:           internAfwijking.IsWarning,
	}
}
