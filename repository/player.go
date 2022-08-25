package repository

import (
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
)

type PlayerRepository struct {
	SimpleRepository
	DB *core.Database `Inject:""`
}

func (r *PlayerRepository) SuppressRecentlyUpdated(days int, idMap map[string]bool) {
	var players []model.Player
	ids := []string{}

	for id := range idMap {
		ids = append(ids, id)
	}

	updated := time.Now().AddDate(0, 0, -1*days)

	result := r.DB.Gorm.Where("updated_at > ? AND id IN ?", updated, ids).Find(&players)

	core.Log.Infof("Suppressed %d Player IDs from scraping!", result.RowsAffected)

	for _, player := range players {
		idMap[player.ID] = false
	}

	return
}
