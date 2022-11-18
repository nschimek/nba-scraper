package model

import "time"

type Player struct {
	ID, Name         string
	Shoots           string `gorm:"type:enum('L', 'R');default:'R'"`
	BirthPlace       string
	BirthCountryCode string
	BirthDate        time.Time
	Height, Weight   int
	Audit
	ModelError
}

type PlayerInjury struct {
	TeamId           string    `gorm:"primaryKey"`
	PlayerId         string    `gorm:"primaryKey"`
	Season           int       `gorm:"primaryKey"`
	SourceUpdateDate time.Time `gorm:"primaryKey"`
	Description      string
	Audit
	ModelError
}
