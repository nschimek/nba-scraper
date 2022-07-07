package repository

import (
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"gorm.io/gorm/clause"
)

type PlayerRepository struct {
	DB *core.Database `Inject:""`
}

func (r *PlayerRepository) CreateBatch(players []model.Player) {
	result := r.DB.Gorm.Clauses(clause.OnConflict{UpdateAll: true}).Create(&players)

	if result.Error != nil {
		core.Log.Error(result.Error)
	} else {
		core.Log.WithField("count", result.RowsAffected).Info("Successfully added/updated players to DB")
	}
}
