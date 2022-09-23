package cmd

import (
	"strings"

	"github.com/nschimek/nba-scraper/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			// temporary empty run command to get the onInit to fire
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is config/default.yaml)")
	rootCmd.PersistentFlags().IntP("season", "s", 2022, "season, use starting year (2022-23 season would be 2022)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode - enable for more detailed logging")

	viper.BindPFlag("season", rootCmd.PersistentFlags().Lookup("season"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.SetDefault("debug", false)
}

func initConfig() {
	viper.SetDefault("use-config-file", true) // overrideable with environment variables only
	viper.SetEnvPrefix("nba")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if viper.GetBool("use-config-file") == true {
		if configFile != "" {
			viper.SetConfigFile(configFile)
		} else {
			viper.AddConfigPath("./config")
			viper.SetConfigType("yaml")
			viper.SetConfigName("default")
			configFile = "./config/default.yaml"
		}
		if err := viper.ReadInConfig(); err == nil {
			core.Log.Infof("Loaded config file: %s", viper.ConfigFileUsed())
		} else {
			core.Log.Fatalf("Could not load config file: %s!", configFile)
		}
	} else {
		core.Log.Info("Config file NOT being used...requiring NBA_ENVIRONMENT_VARIABLES")
	}
}
