package repository

import (
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"gorm.io/gorm/clause"
)

var updateAll = clause.OnConflict{UpdateAll: true}

type SimpleRepository[T any] struct {
	DB *core.Database `Inject:""`
}

func (r *SimpleRepository[T]) Upsert(items []T, label string) {
	core.Log.Infof("Create/updating %s...", label)
	result := r.DB.Gorm.Clauses(updateAll).Create(items)

	if result.Error == nil {
		core.Log.WithField("rows", result.RowsAffected).Infof("Successfully added/updated %s to DB", label)
	}
}

func (r *SimpleRepository[T]) GetRecentlyUpdated(days int, ids []string, label string) ([]string, error) {
	var m T
	var matched []model.IdOnly

	updated := time.Now().AddDate(0, 0, -1*days)
	result := r.DB.Gorm.Model(&m).Where("updated_at > ? AND id IN ?", updated, ids).Find(&matched)

	if result.Error == nil {
		core.Log.Infof("Got %d %s ID(s) updated in past %d days...", result.RowsAffected, label, days)
		m_ids := []string{}

		for _, id := range matched {
			m_ids = append(m_ids, id.ID)
		}

		return m_ids, nil
	} else {
		return nil, result.Error
	}
}
