package cmd

import (
	"errors"

	"github.com/nschimek/nba-scraper/core"
	"github.com/spf13/cobra"
)

var (
	teamsCmd = &cobra.Command{
		Use:   "teams",
		Short: "Scrape NBA team(s) by ID",
		Long: `Scrape teams, team rosters, and team salary data by team IDs.  Team IDs must match the Basketball Reference
	team ID.  Separate with spaces such as: CHI GSW BRK.  This will also cause a scrape of all associated players.
	Note: suppression settings are IGNORED for this command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				s.Team = true
				r.teamIds = appendIds(r.teamIds, core.IdArrayToMap(args))
			} else {
				return errors.New("No team IDs specified!  Please enter team IDs, separated by spaces.  EX: CHI GSW BRK")
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(teamsCmd)
}
