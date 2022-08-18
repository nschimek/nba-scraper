package main

import (
	"time"

	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/scraper"
)

func main() {
	c := core.Setup()

	startDate, _ := time.Parse("2006-01-02", "2021-10-20")
	endDate, _ := time.Parse("2006-01-02", "2021-10-25")

	scheduleScraper := core.Factory[scraper.ScheduleScraper](c.Injector())
	scheduleScraper.ScrapeDateRange(startDate, endDate)

	gameScraper := core.Factory[scraper.GameScraper](c.Injector())
	gameScraper.Scrape(scheduleScraper.GetIds()...)

	// teamScraper := core.Factory[scraper.TeamScraper](c.Injector())
	// teamScraper.Scrape("TOR", "CHI", "BRK", "GSW")

	// playerScraper := core.Factory[scraper.PlayerScraper](c.Injector())
	// playerScraper.Scrape("vandeja01", "curryst01")

	// standingsScraper := core.Factory[scraper.StandingScraper](c.Injector())
	// standingsScraper.Scrape()

	// injuriesScraper := core.Factory[scraper.InjuryScraper](c.Injector())
	// injuriesScraper.Scrape()
}
