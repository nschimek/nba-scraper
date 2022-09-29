package cmd

import (
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/parser"
	"github.com/nschimek/nba-scraper/scraper"
	"github.com/spf13/cobra"
)

var (
	ids                             []string
	startDateString, endDateString  string
	scrapeStandings, scrapeInjuries bool
	gamesCmd                        = &cobra.Command{
		Use:   "games",
		Short: "Scrape games by date range or ID(s)",
		Long: `Scrape game, game stats, and player game stats data by either game IDs or a date range.  
This will also potentially cause scrapes of the corresponding Team and Player pages, depending on Suppression settings.  
Ensure that the Season parameter matches the season you are scraping games from (scraping across seasons is not supported).  
NOTE: if no game IDs or dates are provided, it will default to scraping yesterday.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(ids) > 0 {
				runGameScraperFromIds(core.IdArrayToMap(ids), scrapeStandings, scrapeInjuries)
			} else {
				startDate, endDate, err := dateRangeFromStrings(startDateString, endDateString)
				if err != nil {
					return err
				}
				runGameScraperFromRange(startDate, endDate, scrapeStandings, scrapeInjuries)
			}
			return nil
		},
	}
)

func init() {
	gamesCmd.Flags().StringArrayVarP(&ids, "ids", "i", []string{}, "game ID (specify more than once for multiple) - optional, instead of date range")
	gamesCmd.Flags().StringVarP(&startDateString, "start-date", "s", "", "start date for date range, use YYYY-MM-DD format")
	gamesCmd.Flags().StringVarP(&endDateString, "end-date", "e", "", "end date for date range, use YYYY-MM-DD format")
	gamesCmd.Flags().BoolVarP(&scrapeStandings, "standings", "t", false, "scrape current standings for the given season")
	gamesCmd.Flags().BoolVarP(&scrapeInjuries, "injuries", "j", false, "scrape current injuries (warning: not available for historical seasons)")

	gamesCmd.MarkFlagsRequiredTogether("start-date", "end-date")
	gamesCmd.MarkFlagsMutuallyExclusive("start-date", "ids")
	gamesCmd.MarkFlagsMutuallyExclusive("end-date", "ids")

	rootCmd.AddCommand(gamesCmd)
}

func dateRangeFromStrings(startDateString, endDateString string) (startDate, endDate time.Time, err error) {
	startDate, err = stringToDate(startDateString)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err = stringToDate(endDateString)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return
}

func stringToDate(dateString string) (time.Time, error) {
	if dateString != "" {
		return time.ParseInLocation("2006-01-02", dateString, parser.EST)
	} else {
		return time.Time{}, nil // zero time will be defaulted to yesterday by the schedule scraper
	}
}

func runGameScraperFromRange(startDate, endDate time.Time, scrapeStandings, scrapeInjuries bool) {
	c := core.GetContext()
	scheduleScraper := core.Factory[scraper.ScheduleScraper](c.Injector())
	scheduleScraper.ScrapeDateRange(startDate, endDate)
	runGameScraperFromIds(scheduleScraper.GameIds, scrapeStandings, scrapeInjuries)
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
