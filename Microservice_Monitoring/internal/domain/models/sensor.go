package models

type Sensor struct {
	ID                   int64              `db:"id" json:"id"`
	KunstwerkID          int64              `db:"kunstwerk_id" json:"kunstwerkId"`
	OnderdeelID          *int64             `db:"onderdeel_id" json:"onderdeelId"`
	Geolocation          *string            `db:"geolocation" json:"geolocation"`
	SensorTypeID         int                `db:"sensortype_id" json:"sensortype_id"`
	LastAnalyzedMetingID *int64             `db:"last_analyzed_meting_id" json:"last_analyzed_meting_id"`
	Deleted              bool               `db:"deleted" json:"deleted"`
	SensorConfiguratie   SensorConfiguratie `db:"sensor_configuratie" json:"sensorConfiguratie"`
}

type SensorConfiguratie struct {
	ID              int      `db:"id" json:"id"`
	SensorID        int64    `db:"sensor_id" json:"sensor_id"`
	MinValue        *float64 `db:"min_value" json:"min_value"`
	MaxValue        *float64 `db:"max_value" json:"max_value"`
	MargePercentage *float64 `db:"marge_percentage" json:"marge_percentage"`
}

type SensorType struct {
	ID             int64  `db:"id" json:"id"`
	Naam           string `db:"naam" json:"naam"`
	Eenheid        string `db:"eenheid" json:"eenheid"`
	DrempelIsRange bool   `db:"drempel_is_range" json:"drempel_is_range"`
}

type UpdateSensorConfiguratieRequest struct {
	MinValue        *float64 `json:"minValue"`
	MaxValue        *float64 `json:"maxValue"`
	MargePercentage *float64 `json:"margePercentage"`
}

type CreateSensorRequest struct {
	SensorTypeID             int64                            `json:"sensorTypeId"`
	Geolocation              *string                          `json:"geolocation"`
	ConfiguratieBronSensorID *int64                           `json:"configuratieBronSensorId"`
	Configuratie             *UpdateSensorConfiguratieRequest `json:"configuratie"`
}

type SensorConfiguratieBron struct {
	SensorID           int64              `json:"sensorId"`
	OnderdeelID        *int64             `json:"onderdeelId"`
	OnderdeelNaam      *string            `json:"onderdeelNaam"`
	SensorType         SensorType         `json:"sensorType"`
	SensorConfiguratie SensorConfiguratie `json:"sensorConfiguratie"`
}
