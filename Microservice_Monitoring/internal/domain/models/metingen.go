package models

import (
	"time"

	"github.com/google/uuid"
)

type IncMeting struct {
	SensorID    *int64  `db:"sensor_id" json:"sensorId"`
	KunstwerkID int64   `db:"kunstwerk_id" json:"kunstwerkId"`
	Waarde      float64 `db:"waarde" json:"waarde"`
}

type Meting struct {
	Time        time.Time `db:"time" json:"time"`
	ID          uuid.UUID `db:"id" json:"id"`
	SensorID    *int64    `db:"sensor_id" json:"sensorId"`
	KunstwerkID int64     `db:"kunstwerk_id" json:"kunstwerkId"`
	Waarde      float64   `db:"waarde" json:"waarde"`
	IsHandmatig bool      `db:"is_handmatig" json:"isHandmatig"`
	InspectieID *string   `db:"inspectie_id" json:"inspectieId"`
	Afgehandeld bool      `db:"afgehandeld" json:"afgehandeld"`
}

type Afwijking struct {
	ID            int64     `db:"id" json:"id"`
	MetingID      uuid.UUID `db:"meting_id" json:"metingId"`
	MetingTime    time.Time `db:"meting_time" json:"-"`
	KunstwerkID   int64     `db:"kunstwerk_id" json:"kunstwerkId"`
	SensorID      int64     `db:"sensor_id" json:"sensorId"`
	Time          time.Time `db:"time" json:"time"`
	NormMinWaarde float64   `db:"norm_min_waarde" json:"normMinWaarde"`
	NormMaxWaarde float64   `db:"norm_max_waarde" json:"normMaxWaarde"`
	GemetenWaarde float64   `db:"gemeten_waarde" json:"gemetenWaarde"`
	IsWarning     bool      `db:"is_warning" json:"isWarning"`
}
