package cmd

import (
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfig = "./config/default.yaml"
)

var (
	configFile string
	useConfig  bool

	rootCmd = &cobra.Command{
		Use:   "nba-scraper",
		Short: "Scrape and store NBA data from Basketball Reference",
		Long: `A complete NBA data acquision solution capable of scraping games, 
game team stats, game player stats, teams, team rosters, standings, injuries, and more.`,
		Run: func(cmd *cobra.Command, args []string) {
			core.Log.Info("Started without params, defaulting to Schedule with Injuries and Standings")
			runGameScraperFromRange(time.Time{}, time.Time{}, true, true)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", defaultConfig, "config file")
	rootCmd.PersistentFlags().IntP("season", "n", 0, "optional override of season set in config file; use finishing year (2021-22 season would be 2022)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode - use for more detailed logging")

	viper.BindPFlag("season", rootCmd.PersistentFlags().Lookup("season"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	cobra.OnInitialize(setup)
}

func setup() {
	core.SetupContext(configFile)
}
