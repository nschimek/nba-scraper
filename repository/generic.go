package repository

import (
	"github.com/nschimek/nba-scraper/core"
	"gorm.io/gorm/clause"
)

type GenericRepository struct {
	DB *core.Database `Inject:""`
}

func (r *GenericRepository) Upsert(items any, label string) {
	result := r.DB.Gorm.Clauses(clause.OnConflict{UpdateAll: true}).Create(items)

	if result.Error == nil {
		core.Log.WithField("rows", result.RowsAffected).Infof("Successfully added/updated %s to DB", label)
	}
}
