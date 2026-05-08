package models

import (
	"github.com/google/uuid"
)

type Sensor struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	KunstwerkID  uuid.UUID  `db:"kunstwerk_id" json:"kunstwerkId"`
	OnderdeelID  *uuid.UUID `db:"onderdeel_id" json:"onderdeelId"`
	Geolocation  *string    `db:"geolocation" json:"geolocation"`
	SensorTypeID int        `db:"sensortype_id" json:"sensorTypeId"`
}

type SensorConfiguratie struct {
	ID              int       `db:"id" json:"id"`
	SensorID        uuid.UUID `db:"sensor_id" json:"sensorId"`
	MinValue        *float64  `db:"min_value" json:"minValue"`
	MaxValue        *float64  `db:"max_value" json:"maxValue"`
	MargePercentage *float64  `db:"marge_percentage" json:"margePercentage"`
}
