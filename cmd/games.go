package cmd

import (
	"errors"

	"github.com/nschimek/nba-scraper/core"
	"github.com/spf13/cobra"
)

var (
	ids                []string
	startDate, endDate string
	gamesCmd           = &cobra.Command{
		Use:   "games",
		Short: "Scrape games by date range or ID(s)",
		Long: `Scrape game, game stats, and player game stats data.  This will also potentially cause scrapes of
	the corresponding Team and Player pages, depending on Suppression settings.  Ensure that the Season parameter matches
	the season you are scraping games from.  Note: scraping across seasons in a single run is not supported.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateParams(); err != nil {
				return err
			}

			core.Log.Info(ids)

			return nil
		},
	}
)

func init() {
	gamesCmd.Flags().StringArrayVarP(&ids, "ids", "i", []string{}, "list of game IDs (ex: 202202280BRK,202202280CLE) - specify instead of date range")
	gamesCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "start date for date range, use YYYY-MM-DD format")
	gamesCmd.Flags().StringVarP(&endDate, "end-date", "e", "", "end date for date range, use YYYY-MM-DD format")

	gamesCmd.MarkFlagsRequiredTogether("start-date", "end-date")
	gamesCmd.MarkFlagsMutuallyExclusive("start-date", "ids")
	gamesCmd.MarkFlagsMutuallyExclusive("end-date", "ids")

	rootCmd.AddCommand(gamesCmd)
}

func validateParams() error {
	if startDate == "" && endDate == "" && len(ids) == 0 {
		return errors.New("No parameters provided.  Please provide either start-date AND end-date, or ids.")
	}
	return nil
}
