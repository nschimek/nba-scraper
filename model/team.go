package model

type Team struct {
	ID, Name           string
	Season             int
	TeamPlayers        []TeamPlayer
	TeamPlayerSalaries []TeamPlayerSalary
	Audit
}

type TeamPlayer struct {
	TeamId   string `gorm:"primaryKey"`
	PlayerId string `gorm:"primaryKey"`
	Position string
	Number   int
	Audit
}

type TeamPlayerSalary struct {
	TeamId       string `gorm:"primaryKey"`
	PlayerId     string `gorm:"primaryKey"`
	Salary, Rank int
	Audit
}
