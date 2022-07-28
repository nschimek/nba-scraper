package model

type Team struct {
	ID, Name           string
	Season             int `gorm:"-"`
	TeamPlayers        []TeamPlayer
	TeamPlayerSalaries []TeamPlayerSalary
	Audit
}

type TeamPlayer struct {
	TeamId   string `gorm:"primaryKey"`
	PlayerId string `gorm:"primaryKey"`
	Season   int    `gorm:"primaryKey"`
	Position string
	Number   int
	Audit
}

type TeamPlayerSalary struct {
	TeamId       string `gorm:"primaryKey"`
	PlayerId     string `gorm:"primaryKey"`
	Season       int    `gorm:"primaryKey"`
	Salary, Rank int
	Audit
}

type TeamStanding struct {
	TeamId          string `gorm:"primaryKey"`
	Season          int    `gorm:"primaryKey"`
	Rank            int
	Overall         WinLoss `gorm:"embedded;embeddedPrefix:overall_"`
	Home            WinLoss `gorm:"embedded;embeddedPrefix:home_"`
	Road            WinLoss `gorm:"embedded;embeddedPrefix:road_"`
	East            WinLoss `gorm:"embedded;embeddedPrefix:east_"`
	West            WinLoss `gorm:"embedded;embeddedPrefix:west_"`
	Atlantic        WinLoss `gorm:"embedded;embeddedPrefix:atlantic_"`
	Central         WinLoss `gorm:"embedded;embeddedPrefix:central_"`
	Southeast       WinLoss `gorm:"embedded;embeddedPrefix:southeast_"`
	Northwest       WinLoss `gorm:"embedded;embeddedPrefix:northwest_"`
	Pacific         WinLoss `gorm:"embedded;embeddedPrefix:pacific_"`
	Southwest       WinLoss `gorm:"embedded;embeddedPrefix:southwest_"`
	PreAllStar      WinLoss `gorm:"embedded;embeddedPrefix:pre_all_star_"`
	PostAllStar     WinLoss `gorm:"embedded;embeddedPrefix:post_all_star_"`
	MarginLess3     WinLoss `gorm:"embedded;embeddedPrefix:margin_less_3_"`
	MarginGreater10 WinLoss `gorm:"embedded;embeddedPrefix:margin_greater_10_"`
	October         WinLoss `gorm:"embedded;embeddedPrefix:oct_"`
	November        WinLoss `gorm:"embedded;embeddedPrefix:nov_"`
	December        WinLoss `gorm:"embedded;embeddedPrefix:dec_"`
	January         WinLoss `gorm:"embedded;embeddedPrefix:jan_"`
	February        WinLoss `gorm:"embedded;embeddedPrefix:feb_"`
	March           WinLoss `gorm:"embedded;embeddedPrefix:mar_"`
	April           WinLoss `gorm:"embedded;embeddedPrefix:apr_"`
}

type WinLoss struct {
	Wins, Losses int
}
