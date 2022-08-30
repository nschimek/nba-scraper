package repository

import (
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
	"github.com/sirupsen/logrus"
)

type TeamRepository struct {
	DB *core.Database `Inject:""`
}

func (tr *TeamRepository) UpsertTeams(teams []model.Team) {
	core.Log.WithField("teams", len(teams)).Infof("Create/Updating Teams, Team Players, and Team Player Salaries...")
	for _, team := range teams {
		r1 := tr.DB.Gorm.Clauses(updateAll).Omit("TeamPlayers", "TeamPlayerSalaries").Create(&team)

		if r1.Error == nil {
			tr.DB.Gorm.Delete(&model.TeamPlayer{}, &model.TeamPlayer{TeamId: team.ID, Season: team.Season})
			r2 := tr.DB.Gorm.Create(&team.TeamPlayers)
			tr.DB.Gorm.Delete(&model.TeamPlayerSalary{}, &model.TeamPlayerSalary{TeamId: team.ID, Season: team.Season})
			r3 := tr.DB.Gorm.Create(&team.TeamPlayerSalaries)

			if r2.Error == nil && r3.Error == nil {
				core.Log.WithFields(logrus.Fields{
					"players":  r2.RowsAffected,
					"salaries": r3.RowsAffected,
				}).Infof("Successfully create/updated team %s along with players and salaries", team.ID)
			}
		}
	}
}
