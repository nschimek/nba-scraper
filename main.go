package main

import (
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/parser"
	"github.com/nschimek/nba-scraper/scraper"
)

func main() {
	c := core.Setup()

	// standingsScraper := core.Factory[scraper.StandingScraper](c.Injector())
	// standingsScraper.Scrape()

	// injuriesScraper := core.Factory[scraper.InjuryScraper](c.Injector())
	// injuriesScraper.Scrape()

	startDate, _ := time.ParseInLocation("2006-01-02", "2021-10-20", parser.EST)
	endDate, _ := time.ParseInLocation("2006-01-02", "2021-10-20", parser.EST)

	scheduleScraper := core.Factory[scraper.ScheduleScraper](c.Injector())
	scheduleScraper.ScrapeDateRange(startDate, endDate)

	gameScraper := core.Factory[scraper.GameScraper](c.Injector())
	gameScraper.Scrape(scheduleScraper.GameIds)

	teamScraper := core.Factory[scraper.TeamScraper](c.Injector())
	teamScraper.Scrape(gameScraper.TeamIds)

	playerScraper := core.Factory[scraper.PlayerScraper](c.Injector())
	playerScraper.Scrape(gameScraper.PlayerIds)

	playerScraper.Persist()
	teamScraper.Persist()
	gameScraper.Persist()
}
