package cmd

import (
	"github.com/nschimek/nba-scraper/core"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version of the NBA Scraper",
	Long:  `Display the current version of the NBA Scraper`,
	Run: func(cmd *cobra.Command, args []string) {
		core.Log.Infof("Current NBA Scraper version: %s", core.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
