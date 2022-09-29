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
			if len(args) == 0 {
				return errors.New("No game IDs specified!  Please enter game IDs, separated by spaces.  EX: 202202280BRK 202202280CLE")
			} else {
				runGameScraperFromIds(core.IdArrayToMap(args), false, false)
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(gamesCmd)
}

func runGameScraperFromIds(idMap map[string]struct{}, scrapeStandings, scrapeInjuries bool) {
	c := core.GetContext()
	standingsScraper := core.Factory[scraper.StandingScraper](c.Injector())
	injuriesScraper := core.Factory[scraper.InjuryScraper](c.Injector())
	gameScraper := core.Factory[scraper.GameScraper](c.Injector())
	teamScraper := core.Factory[scraper.TeamScraper](c.Injector())
	playerScraper := core.Factory[scraper.PlayerScraper](c.Injector())

	if scrapeStandings {
		standingsScraper.Scrape()
	}
	if scrapeInjuries {
		injuriesScraper.Scrape()
	}

	gameScraper.Scrape(idMap)
	teamScraper.Scrape(core.ConsolidateIdMaps(standingsScraper.TeamIds, injuriesScraper.TeamIds, gameScraper.TeamIds))
	playerScraper.Scrape(core.ConsolidateIdMaps(injuriesScraper.PlayerIds, gameScraper.PlayerIds, teamScraper.PlayerIds))

	playerScraper.Persist()
	teamScraper.Persist()
	gameScraper.Persist()
	standingsScraper.Persist()
	injuriesScraper.Persist()
}
