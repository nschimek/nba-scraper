package repository

import (
	"github.com/nschimek/nba-scraper/core"
	"gorm.io/gorm/clause"
)

var updateAll = clause.OnConflict{UpdateAll: true}

type GenericRepository struct {
	DB *core.Database `Inject:""`
}

func (r *GenericRepository) Upsert(items any, label string) {
	core.Log.Infof("Create/updating %s...", label)
	result := r.DB.Gorm.Clauses(updateAll).Create(items)

	if result.Error == nil {
		core.Log.WithField("rows", result.RowsAffected).Infof("Successfully added/updated %s to DB", label)
	}
}
