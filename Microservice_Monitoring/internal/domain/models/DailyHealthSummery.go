package models

import "time"

type Status string

const (
	StatusHealthy  Status = "healthy"
	StatusWarning  Status = "warning"
	StatusCritical Status = "critical"
	StatusOffline  Status = "offline"
)

type DailyHealthSummary struct {
	KunstwerkID               int64           `json:"kunstwerkId"`
	KunstwerkBeheerIdentifier string          `json:"kunstwerkBeheerIdentifier"`
	KunstwerkDetail           KunstwerkDetail `json:"kunstwerkDetail"`
	Tijd                      time.Time       `json:"tijd"`
	Status                    Status          `json:"status"`
	AantalSensoren            int             `json:"aantalSensoren"`
	AantalActieveSensoren     int             `json:"aantalActieveSensoren"`
	AantalAfwijkendeSensoren  int             `json:"aantalAfwijkendeSensoren"`
	AantalAfwijkingen         int             `json:"aantalAfwijkingen"`
}

type DailyHealthUpdate struct {
	KunstwerkID              int64  `json:"kunstwerkId"`
	Status                   Status `json:"status"`
	AantalSensoren           int    `json:"aantalSensoren"`
	AantalActieveSensoren    int    `json:"aantalActieveSensoren"`
	AantalAfwijkendeSensoren int    `json:"aantalAfwijkendeSensoren"`
	AantalAfwijkingen        int    `json:"aantalAfwijkingen"`
}
