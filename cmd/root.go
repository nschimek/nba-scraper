package cmd

import (
	"strings"

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
		// Run: func(cmd *cobra.Command, args []string) {
		// 	// temporary empty run command to get the onInit to fire
		// },
	}
)

type TestConfig struct {
	Season        int
	Debug         bool
	UseConfigFile bool `mapstructure:"use-config-file"`
	Suppression   struct {
		Team   int
		Player int
	}
	Database struct {
		User, Password, Location, Port, Name string
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", defaultConfig, "config file (default is "+defaultConfig+")")
	rootCmd.PersistentFlags().IntP("season", "s", 2022, "season, use finishing year (2021-22 season would be 2022)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode - enable for more detailed logging")

	viper.BindPFlag("season", rootCmd.PersistentFlags().Lookup("season"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	core.Log.Infof("%+v\n", getConmfig())
}

func getConmfig() *TestConfig {
	viper.SetDefault("use-config-file", true) // overrideable with environment variables only
	viper.SetEnvPrefix("ns")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if viper.GetBool("use-config-file") == true {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err == nil {
			core.Log.Infof("Loaded config file: %s", viper.ConfigFileUsed())
		} else {
			core.Log.Fatalf("Could not load config file: %s!", configFile)
		}
	} else {
		core.Log.Info("Config file NOT being used...requiring NS_ENVIRONMENT_VARIABLES")
	}

	config := &TestConfig{}
	if err := viper.Unmarshal(config); err != nil {
		core.Log.Fatalf("Error decoding Config struct: %v", err)
	}
	return config
}
