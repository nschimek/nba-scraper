package repository

import (
	"github.com/nschimek/nba-scraper/context"
	"github.com/nschimek/nba-scraper/model"
	"gorm.io/gorm/clause"
)

type PlayerRepository struct {
	DB *context.Database `Inject:""`
}

func (r *PlayerRepository) CreateBatch(player []model.Player) {
	result := r.DB.Gorm.Clauses(clause.OnConflict{UpdateAll: true}).Create(player)

	if result.Error != nil {
		context.Log.Error(result.Error)
	} else {
		context.Log.WithField("count", result.RowsAffected).Info("Successfully added players to DB")
	}
}
