package repository

import (
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"github.com/nschimek/nba-scraper/parser"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GameRepository struct {
	DB     *core.Database `Inject:""`
	Config *core.Config   `Inject:""`
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
			} else {
				core.Log.Errorf("Error(s) occurred while loading %s", game.ID)
			}
		}

	}
}

func (r *GameRepository) GetMostRecentGame() (time.Time, error) {
	var games []model.Game

	// get the start time of the latest game in the DB for the currently configured season
	result := r.DB.Gorm.Select("id", "start_time").Where("season = ?", r.Config.Season).Order("start_time desc").Limit(1).Find(&games)

	if result.Error == nil {
		if len(games) > 0 {
			st, id := games[0].StartTime, games[0].ID
			core.Log.Infof("Most recent game for season %d in DB was %s at %s", r.Config.Season, id, st.Format(core.DateRangeFormat))
			return st, nil
		} else {
			core.Log.Warnf("No games for season %d, so could not get most recent; defaulting to October 1st...", r.Config.Season)
			return time.Date(r.Config.Season-1, 10, 1, 0, 0, 0, 0, parser.EST), nil
		}
	} else {
		return time.Time{}, result.Error
	}
}
