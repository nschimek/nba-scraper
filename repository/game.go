package repository

import (
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GameRepository struct {
	DB *core.Database `Inject:""`
}

func (r *GameRepository) UpsertGames(games []model.Game) {
	core.Log.WithField("games", len(games)).Infof("Create/Updating Games, Game Stats, and Game Player Stats...")
	for _, game := range games {
		r1 := r.DB.Gorm.Clauses(updateAll).Omit(
			"HomeLineScores",
			"AwayLineScores",
			"HomeFourFactors",
			"AwayFourFactors",
			"GamePlayers",
			"GamePlayersBasicStats",
			"GamePlayersAdvancedStats").Create(&game)

		if r1.Error == nil {
			var results [7]*gorm.DB
			var errors int

			results[0] = r.DB.Gorm.Clauses(updateAll).Create(&game.HomeLineScores)
			results[1] = r.DB.Gorm.Clauses(updateAll).Create(&game.AwayLineScores)
			results[2] = r.DB.Gorm.Clauses(updateAll).Create(&game.HomeFourFactors)
			results[3] = r.DB.Gorm.Clauses(updateAll).Create(&game.AwayFourFactors)
			results[4] = r.DB.Gorm.Clauses(updateAll).Create(&game.GamePlayers)
			results[5] = r.DB.Gorm.Clauses(updateAll).Create(&game.GamePlayersBasicStats)
			results[6] = r.DB.Gorm.Clauses(updateAll).Create(&game.GamePlayersAdvancedStats)

			for _, r := range results {
				if r.Error != nil {
					errors++
				}
			}

			if errors == 0 {
				core.Log.WithFields(logrus.Fields{
					"line scores":  results[0].RowsAffected + results[1].RowsAffected,
					"four factors": results[2].RowsAffected + results[3].RowsAffected,
					"game players": results[4].RowsAffected,
					"basic stats":  results[5].RowsAffected,
					"adv stats":    results[6].RowsAffected,
				}).Infof("Succesfully create/updated game %s along with all scores and stats", game.ID)
			}
		}

	}
}
