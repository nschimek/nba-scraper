package repository

import (
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	suppressedLogMessage = "%d (of %d) %s ID(s) were updated in last %d days and suppressed from scraping"
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
		core.Log.Infof(suppressedLogMessage, result.RowsAffected, len(ids), label, days)
		m_ids := []string{}

		for _, id := range matched {
			m_ids = append(m_ids, id.ID)
		}

		return m_ids, nil
	} else {
		return nil, result.Error
	}
}

func runQueryIfNotEmpty(length int, query *gorm.DB) *gorm.DB {
	if length > 0 {
		return query
	} else {
		core.Log.Warn("there were no records to add, so the query was skipped (check stats)...")
		return new(gorm.DB)
	}
}
