package main

import (
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/parser"
	"github.com/nschimek/nba-scraper/scraper"
)

func main() {
	c := core.Setup()

	scheduleScraper := core.Factory[scraper.ScheduleScraper](c.Injector())
	standingsScraper := core.Factory[scraper.StandingScraper](c.Injector())
	injuriesScraper := core.Factory[scraper.InjuryScraper](c.Injector())
	gameScraper := core.Factory[scraper.GameScraper](c.Injector())
	teamScraper := core.Factory[scraper.TeamScraper](c.Injector())
	playerScraper := core.Factory[scraper.PlayerScraper](c.Injector())

	standingsScraper.Scrape()
	injuriesScraper.Scrape()

	startDate, _ := time.ParseInLocation("2006-01-02", "2022-04-01", parser.EST)
	endDate, _ := time.ParseInLocation("2006-01-02", "2022-06-30", parser.EST)

	scheduleScraper.ScrapeDateRange(startDate, endDate)

	playerScraper.Persist()
	teamScraper.Persist()
	gameScraper.Persist()
	standingsScraper.Persist()
	injuriesScraper.Persist()
}
