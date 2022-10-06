package cmd

import (
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/scraper"
	"github.com/spf13/cobra"
)

var (
	standingsCmd = &cobra.Command{
		Use:   "standings",
		Short: "Scrape standings",
		Long: `Scrape currently available standings.  This will potentially trigger a scrape of the corresponding
	team, team roster, team salary, and all player pages.  
	NOTE: Historical standings are always as of the last day of the regular season.`,
		Run: func(cmd *cobra.Command, args []string) {
			s = &scrapers{standing: true, team: true, player: true}
		},
	}
)

func init() {
	rootCmd.AddCommand(standingsCmd)
}

// Gets conditionally called by the rootCmd PersistentPostRun
func runStandingScraper() {
	standingsScraper := core.Factory[scraper.StandingScraper](core.GetInjector())
	standingsScraper.Scrape()
	r.teamIds = appendIds(r.teamIds, standingsScraper.TeamIds)
}
