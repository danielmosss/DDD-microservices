package models

type KunstwerkType struct {
	ID           int64  `db:"id" json:"id"`
	Naam         string `db:"naam" json:"naam"`
	Beschrijving string `db:"beschrijving" json:"beschrijving"`
}

type Kunstwerk struct {
	ID               int64  `db:"id" json:"id"`
	BeheerIdentifier string `db:"beheer_identifier" json:"beheerIdentifier"`
	Naam             string `db:"naam" json:"naam"`
	Geolocation      string `db:"geolocation" json:"geolocation"`
	KunstwerkTypeID  int    `db:"kunstwerktype_id" json:"kunstwerkTypeId"`
	Beschrijving     string `db:"beschrijving" json:"beschrijving"`
	Deleted          bool   `db:"deleted" json:"deleted"`
}
