package model

type Team struct {
	Id, Name           string
	Season             int
	TeamPlayers        []TeamPlayer
	TeamPlayerSalaries []TeamPlayerSalary
	Audit
}

type TeamPlayer struct {
	PlayerId, Position string
	Number             int
	Audit
}

type TeamPlayerSalary struct {
	PlayerId     string
	Salary, Rank int
	Audit
}
