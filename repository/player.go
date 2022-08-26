package repository

import (
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
)

type PlayerRepository struct {
	DB *core.Database `Inject:""`
}

func (r *PlayerRepository) SuppressRecentlyUpdated(days int, idMap map[string]bool) {
	var players []model.Player

	updated := time.Now().AddDate(0, 0, -1*days)
	result := r.DB.Gorm.Where("updated_at > ? AND id IN ?", updated, core.IdMapToArray(idMap)).Find(&players)

	core.Log.Infof("Suppressed %d Player IDs from scraping!", result.RowsAffected)

	for _, player := range players {
		idMap[player.ID] = false
	}

	return
}
