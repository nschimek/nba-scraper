package repository

import (
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"github.com/sirupsen/logrus"
)

type GameRepository struct {
	DB *core.Database `Inject:""`
}

func (r *GameRepository) UpsertGames(games []model.Game) {
	core.Log.WithField("games", len(games)).Infof("Create/Updating Games, Game Stats, and Game Player Stats...")
	for _, game := range games {
		r.DB.Gorm.Clauses(updateAll).Omit(
			"HomeLineScores",
			"AwayLineScores",
			"HomeFourFactors",
			"AwayFourFactors",
			"GamePlayers",
			"GamePlayersBasicStats",
			"GamePlayersAdvancedStats").Create(&game)

		r1 := r.DB.Gorm.Clauses(updateAll).Create(&game.HomeLineScores)
		r2 := r.DB.Gorm.Clauses(updateAll).Create(&game.AwayLineScores)
		r3 := r.DB.Gorm.Clauses(updateAll).Create(&game.HomeFourFactors)
		r4 := r.DB.Gorm.Clauses(updateAll).Create(&game.AwayFourFactors)
		r5 := r.DB.Gorm.Clauses(updateAll).Create(&game.GamePlayers)
		r6 := r.DB.Gorm.Clauses(updateAll).Create(&game.GamePlayersBasicStats)
		r7 := r.DB.Gorm.Clauses(updateAll).Create(&game.GamePlayersAdvancedStats)

		core.Log.WithFields(logrus.Fields{
			"line scores":  r1.RowsAffected + r2.RowsAffected,
			"four factors": r3.RowsAffected + r4.RowsAffected,
			"game players": r5.RowsAffected,
			"basic stats":  r6.RowsAffected,
			"adv stats":    r7.RowsAffected,
		}).Infof("Succesfully create/updated game %s along with all scores and stats", game.ID)
	}
}
