package cmd

import (
	"errors"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/scraper"
	"github.com/spf13/cobra"
)

var (
	gamesCmd = &cobra.Command{
		Use:   "games",
		Short: "Scrape NBA game(s) data by ID(s)",
		Long: `Scrape game, game stats, and player game stats data by game IDs, separated by spaces.  
This will also potentially cause scrapes of the corresponding Team and Player pages, depending on Suppression settings.  
NOTE: Check that the Season parameter matches the season you are scraping games from as scraping across seasons is not supported.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				s = &scrapers{game: true, team: true, player: true}
				r.gameIds = appendIds(r.gameIds, core.IdArrayToMap(args))
			} else {
				return errors.New("No game IDs specified!  Please enter game IDs, separated by spaces.  EX: 202202280BRK 202202280CLE")
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(gamesCmd)
}

func runGameScraper() {
	gameScraper := core.Factory[scraper.GameScraper](core.GetInjector())
	gameScraper.Scrape(r.gameIds)
	r.teamIds = appendIds(r.teamIds, gameScraper.TeamIds)
	r.playerIds = appendIds(r.playerIds, gameScraper.PlayerIds)
}
