package parser

import "time"

type Team struct {
	Id, Name string
}

type TeamRoster struct {
	TeamId, PlayerId, Name, Position string
	Season, Number                   int
}

type TeamInjuryReport struct {
	TeamId, PlayerId, Description string
	UpdateDate                    time.Time
}

type TeamPlayerSalaries struct {
	TeamId, PlayerId string
	Season, Salary   int
}
