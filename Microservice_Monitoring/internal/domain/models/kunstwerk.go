package models

import (
	"github.com/google/uuid"
)

type KunstwerkType struct {
	ID           int64   `db:"id" json:"id"`
	Naam         string  `db:"naam" json:"naam"`
	Beschrijving *string `db:"beschrijving" json:"beschrijving"`
}

type Kunstwerk struct {
	ID              uuid.UUID `db:"id" json:"id"`
	Naam            string    `db:"naam" json:"naam"`
	Geolocation     *string   `db:"geolocation" json:"geolocation"`
	KunstwerkTypeID int       `db:"kunstwerktype_id" json:"kunstwerkTypeId"`
	Beschrijving    *string   `db:"beschrijving" json:"beschrijving"`
	Deleted         bool      `db:"deleted" json:"deleted"`
}
