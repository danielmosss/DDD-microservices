package main

type Meting struct {
	SensorID    string  `db:"sensor_id" json:"sensorId"`
	KunstwerkID string  `db:"kunstwerk_id" json:"kunstwerkId"`
	Waarde      float64 `db:"waarde" json:"waarde"`
}
