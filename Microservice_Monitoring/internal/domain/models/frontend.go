package models

type TreeOnderdeel struct {
	ID         int64            `json:"id"`
	Naam       string           `json:"naam"`
	ParentId   int64            `json:"parent_id"`
	Sensoren   []int64          `json:"sensoren"`
	Onderdelen []*TreeOnderdeel `json:"onderdelen"`
}

type KunstwerkTreeResponse struct {
	Kunstwerk     KunstwerkDetail  `json:"kunstwerkdetail"`
	LosseSensoren []int64          `json:"losseSensoren"`
	Onderdelen    []*TreeOnderdeel `json:"onderdelen"`
}

type SensorDetailResponse struct {
	ID                  int64              `json:"id"`
	SensorTypeID        int                `json:"sensorTypeId"`
	SensorConfiguratie  SensorConfiguratie `json:"sensorConfiguratie"`
	LaatsteMetingWaarde *float64           `json:"laatsteMeting"`
	Status              string             `json:"status"`
}
