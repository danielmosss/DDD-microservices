package main

type Meting struct {
	SensorID    int64   `db:"sensor_id" json:"sensorId"`
	KunstwerkID int64   `db:"kunstwerk_id" json:"kunstwerkId"`
	Waarde      float64 `db:"waarde" json:"waarde"`
}
