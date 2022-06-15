package model

import "time"

type Player struct {
	ID, Name             string
	Shoots               string `gorm:"type:enum('L', 'R');default:'R'"`
	BirthPlace           string
	BirthCountryCode     string
	BirthDate            time.Time
	Height, Weight       int
	CreatedAt, UpdatedAt time.Time
}
