package models

type Sensor struct {
	ID           int64  `db:"id" json:"id"`
	KunstwerkID  int64  `db:"kunstwerk_id" json:"kunstwerkId"`
	OnderdeelID  int64  `db:"onderdeel_id" json:"onderdeelId"`
	Geolocation  string `db:"geolocation" json:"geolocation"`
	SensorTypeID int    `db:"sensortype_id" json:"sensorTypeId"`
}

type SensorConfiguratie struct {
	ID              int      `db:"id" json:"id"`
	SensorID        int64    `db:"sensor_id" json:"sensorId"`
	MinValue        *float64 `db:"min_value" json:"minValue"`
	MaxValue        *float64 `db:"max_value" json:"maxValue"`
	MargePercentage *float64 `db:"marge_percentage" json:"margePercentage"`
}
