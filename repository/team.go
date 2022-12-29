package repository

import (
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TeamRepository struct {
	DB *core.Database `Inject:""`
}

func (tr *TeamRepository) UpsertTeams(teams []model.Team) {
	core.Log.WithField("teams", len(teams)).Infof("Create/Updating Teams, Team Players, and Team Player Salaries...")
	for _, team := range teams {
		r1 := tr.DB.Gorm.Clauses(updateAll).Omit("TeamPlayers", "TeamPlayerSalaries").Create(&team)

		if r1.Error == nil {
			var results [2]*gorm.DB

			tr.DB.Gorm.Delete(&model.TeamPlayer{}, &model.TeamPlayer{TeamId: team.ID, Season: team.Season})
			results[0] = tr.DB.Gorm.Create(&team.TeamPlayers)
			tr.DB.Gorm.Delete(&model.TeamPlayerSalary{}, &model.TeamPlayerSalary{TeamId: team.ID, Season: team.Season})
			results[1] = tr.DB.Gorm.Create(&team.TeamPlayerSalaries)

			if results[0].Error == nil && results[1].Error == nil {
				core.Log.WithFields(logrus.Fields{
					"players":  results[0].RowsAffected,
					"salaries": results[1].RowsAffected,
				}).Infof("Successfully create/updated team %s along with players and salaries", team.ID)
			} else {
				core.Log.Errorf("Error(s) occurred while loading %s", team.ID)
			}
		}
	}
}
