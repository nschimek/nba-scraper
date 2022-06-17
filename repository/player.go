package repository

import (
	"github.com/nschimek/nba-scraper/context"
	"github.com/nschimek/nba-scraper/model"
)

type PlayerRepository struct {
	DB *context.Database `Inject:""`
}

func (r *PlayerRepository) CreateBatch(players []model.Player) {
	result := r.DB.Gorm.Create(players)

	if result.Error != nil {
		context.Log.Error(result.Error)
	} else {
		context.Log.WithField("count", result.RowsAffected).Info("Successfully added players to DB")
	}
}
