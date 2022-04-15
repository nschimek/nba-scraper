package parser

import "time"

type Team struct {
	Id, Name string
	Season   int
}

type TeamRoster struct {
	TeamId, PlayerId, Name, Position string
	Number                           int
}

type TeamInjuryReport struct {
	TeamId, PlayerId, Description string
	UpdateDate                    time.Time
}

type TeamPlayerSalaries struct {
	TeamId, PlayerId string
	Salary           int
}
