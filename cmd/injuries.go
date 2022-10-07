package cmd

import (
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/scraper"
	"github.com/spf13/cobra"
)

var (
	injuriesCmd = &cobra.Command{
		Use:   "injuries",
		Short: "Scrape injuries",
		Long: `Scrape current injuries by player IDs.  This will potentially trigger a scrape of the corresponding
	team and player pages (for each injured player).  
	WARNING: BR does not designate a season for the injuries, and there are no historical injuries available.  
	Therefore, all injuries are always loaded under the currently configured season.`,
		Run: func(cmd *cobra.Command, args []string) {
			s = &scrapers{injury: true, team: true, player: true}
		},
	}
)

func init() {
	rootCmd.AddCommand(injuriesCmd)
}

// Gets conditionally called by the rootCmd PersistentPostRun
func runInjuryScraper() {
	injuryScraper := core.Factory[scraper.InjuryScraper](core.GetInjector())
	injuryScraper.Scrape()
	r.teamIds = appendIds(r.teamIds, injuryScraper.TeamIds)
	r.playerIds = appendIds(r.playerIds, injuryScraper.PlayerIds)
	s.persist = append(s.persist, injuryScraper)
}
