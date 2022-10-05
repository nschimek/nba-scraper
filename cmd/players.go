package cmd

import (
	"errors"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/scraper"
	"github.com/spf13/cobra"
)

var (
	playersCmd = &cobra.Command{
		Use:   "teams",
		Short: "Scrape NBA players(s) by ID",
		Long: `Scrape players by player IDs.  Player IDs must match the Basketball Reference player ID.  
	Separate with spaces such as: curryst01 lavinza01 jamesle01. 
	NOTE: suppression settings are ignored for this command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				config.Suppression.Player = 0
				s = &scrapers{player: true}
				r.teamIds = appendIds(r.playerIds, core.IdArrayToMap(args))
			} else {
				return errors.New("No player IDs specified!  Please enter player IDs, separated by spaces.  EX: curryst01 lavinza01 jamesle01")
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(playersCmd)
}

// Gets conndtionally called by the rootCmd PersistentPostRun
func runPlayerScraper() {
	playerScraper := core.Factory[scraper.PlayerScraper](core.GetInjector())
	playerScraper.Scrape(r.playerIds)
}
