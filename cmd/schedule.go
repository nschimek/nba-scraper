package cmd

import (
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/parser"
	"github.com/nschimek/nba-scraper/scraper"
	"github.com/spf13/cobra"
)

var (
	startDateString, endDateString  string
	scrapeStandings, scrapeInjuries bool
	scheduleCmd                     = &cobra.Command{
		Use:   "schedule",
		Short: "Scrape NBA game(s) data by the schedule",
		Long: `Scrape game, game stats, and player game stats data via the schedule, using an optional date range.
If no date range is specified, it uses the date of the last game loaded for the current season until today.
This will also potentially cause scrapes of the corresponding Team and Player pages, depending on Suppression settings.  
Injuries and standings can also be optionally scraped; however, they have limited historical support (more info in help).
NOTE: Check that the Season parameter matches the season you are scraping games from as scraping across seasons is not supported.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if startDate, endDate, err := dateRangeFromStrings(startDateString, endDateString); err != nil {
				return err
			} else {
				r.startDate = startDate
				r.endDate = endDate
			}
			s = &scrapers{schedule: true, game: true, team: true, player: true, standing: scrapeStandings, injury: scrapeInjuries}
			return nil
		},
	}
)

func init() {
	scheduleCmd.Flags().StringVarP(&startDateString, "start-date", "s", "", "start date for date range, use YYYY-MM-DD format")
	scheduleCmd.Flags().StringVarP(&endDateString, "end-date", "e", "", "end date for date range, use YYYY-MM-DD format")
	scheduleCmd.Flags().BoolVarP(&scrapeStandings, "standings", "t", false, "scrape current standings for the given season (historical will be last day)")
	scheduleCmd.Flags().BoolVarP(&scrapeInjuries, "injuries", "j", false, "scrape current injuries (warning: not available for historical seasons)")

	scheduleCmd.MarkFlagsRequiredTogether("start-date", "end-date")

	rootCmd.AddCommand(scheduleCmd)
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

// Gets conditionally called by the rootCmd PersistentPostRun
func runScheduleScraper() {
	scheduleScraper := core.Factory[scraper.ScheduleScraper](core.GetInjector())
	scheduleScraper.ScrapeDateRange(r.startDate, r.endDate)
	r.gameIds = appendIds(r.gameIds, scheduleScraper.GameIds)
}
