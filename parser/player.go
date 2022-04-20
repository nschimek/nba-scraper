package parser

import "time"

type Player struct {
	Id, Name, Position, Shoots string
	Height, Weight             int
	Birthday                   time.Time
	Active                     bool
	PlayerGameSummaries        []PlayerGameSummary
}

type PlayerGameSummary struct {
	Season, Stat    string
	Min, Max, Count int
}
