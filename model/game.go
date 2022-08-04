package model

import "time"

type Game struct {
	ID, Location                     string
	Type                             string `gorm:"type:enum('R', 'P');default:'P'"`
	Season, Quarters                 int
	StartTime                        time.Time
	Home                             GameTeam        `gorm:"embedded;embeddedPrefix:home_"`
	Away                             GameTeam        `gorm:"embedded;embeddedPrefix:away_"`
	HomeLineScore, AwayLineScore     []GameLineScore // these will end up in their own table due to the possiblity of OT
	HomeFourFactors, AwayFourFactors GameFourFactor
	GamePlayers                      []GamePlayer
	GamePlayersBasicStats            []GamePlayerBasicStat
	GamePlayersAdvancedStats         []GamePlayerAdvancedStat
}

type GameTeam struct {
	TeamId              string
	Result              string `gorm:"type:enum('W', 'L');default:'W'"`
	Wins, Losses, Score int
}

type GameFourFactor struct {
	GameId                                                                       string `gorm:"primaryKey"`
	TeamId                                                                       string `gorm:"primaryKey"`
	Pace, EffectiveFgPct, TurnoverPct, OffensiveRbPct, FtPerFga, OffensiveRating float64
}

type GameLineScore struct {
	GameId         string `gorm:"primaryKey"`
	TeamId         string `gorm:"primaryKey"`
	Quarter, Score int
}

type GamePlayer struct {
	GameId   string `gorm:"primaryKey"`
	TeamId   string `gorm:"primaryKey"`
	PlayerId string `gorm:"primaryKey"`
	Status   string `gorm:"type:enum('S', 'R', 'D', 'I');default:'I'"`
}

type GamePlayerBasicStat struct {
	GameId                                                                                                  string `gorm:"primaryKey"`
	TeamId                                                                                                  string `gorm:"primaryKey"`
	PlayerId                                                                                                string `gorm:"primaryKey"`
	Quarter                                                                                                 int    `gorm:"primaryKey"`
	TimePlayed                                                                                              time.Duration
	FieldGoals, FieldGoalsAttempted, ThreePointers, ThreePointersAttempted, FreeThrows, FreeThrowsAttempted int
	FieldGoalPct, ThreePointersPct, FreeThrowsPct                                                           float64
	OffensiveRB, DefensiveRB, TotalRB, Assists, Steals, Blocks, Turnovers, PersonalFouls, Points, PlusMinus int
}

type GamePlayerAdvancedStat struct {
	GameId                                                                                                                string `gorm:"primaryKey"`
	TeamId                                                                                                                string `gorm:"primaryKey"`
	PlayerId                                                                                                              string `gorm:"primaryKey"`
	TrueShootingPct, EffectiveFgPct, ThreePtAttemptRate, FreeThrowAttemptRate, OffensiveRbPct, DefensiveRbPct, TotalRbPct float64
	AssistPct, StealPct, BlockPct, TurnoverPct, UsagePct, BoxPlusMinus                                                    float64
	OffensiveRating, DefensiveRating                                                                                      int
}
