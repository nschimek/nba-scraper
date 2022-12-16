package model

import "time"

type Schedule struct {
	StartTime                         time.Time
	GameId, VisitorTeamId, HomeTeamId string
	Played                            bool
	ModelError
}
