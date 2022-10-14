package cmd

import (
	"os"
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/scraper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfig = "./config/default.yaml"
)

type results struct {
	playerIds, teamIds, gameIds map[string]struct{}
	startDate, endDate          time.Time
}

type scrapers struct {
	schedule, injury, standing, game, team, player bool
	persist                                        []scraper.Scraper
}

var (
	config     *core.Config
	r          *results
	s          *scrapers
	configFile string
	useConfig  bool

	rootCmd = &cobra.Command{
		Use:   "nba-scraper",
		Short: "Scrape and store NBA data from Basketball Reference",
		Long: `A complete NBA data acquision solution capable of scraping games, 
game team stats, game player stats, teams, team rosters, standings, injuries, and more.`,
		Run: func(cmd *cobra.Command, args []string) {
			core.Log.Info("Started without commands or parameters, defaulting to Schedule with Injuries and Standings")
			s = &scrapers{standing: true, injury: true, schedule: true, game: true, team: true, player: true}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if cmd.Use == "version" {
				os.Exit(0) // the version command should do nothing more
			}

			conditionallyRun(runStandingScraper, s.standing)
			conditionallyRun(runInjuryScraper, s.injury)
			conditionallyRun(runScheduleScraper, s.schedule)
			conditionallyRun(runGameScraper, s.game)
			conditionallyRun(runTeamScraper, s.team)
			conditionallyRun(runPlayerScraper, s.player)

			// need to persist in the reverse order the scrapers were added to the list
			for i := len(s.persist) - 1; i >= 0; i-- {
				s.persist[i].Persist()
			}
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
	config = core.Factory[core.Config](core.GetInjector())
	r = &results{
		playerIds: make(map[string]struct{}),
		teamIds:   make(map[string]struct{}),
		gameIds:   make(map[string]struct{}),
	}
}

func appendIds(target, ids map[string]struct{}) map[string]struct{} {
	for k, v := range ids {
		target[k] = v
	}
	return target
}

func conditionallyRun(runScraper func(), condition bool) {
	if condition {
		runScraper()
	}
}
