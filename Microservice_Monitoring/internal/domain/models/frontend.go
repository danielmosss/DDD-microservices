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
	ID                 int64              `json:"id"`
	SensorType         SensorType         `json:"sensorType"`
	SensorConfiguratie SensorConfiguratie `json:"sensorConfiguratie"`
	LaatsteMeting      *Meting            `json:"laatsteMeting,omitempty"`
	Afwijking          *Afwijking         `json:"afwijking,omitempty"`
	Status             Status             `json:"status,omitempty"`
}
