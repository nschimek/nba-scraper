package main

import (
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/scraper"
)

func main() {
	c := core.Setup()

	gameScraper := core.Factory[scraper.GameScraper](c.Injector())
	gameScraper.Scrape("202110300WAS", "202204180GSW")

	// startDate, _ := time.Parse("2006-01-02", "2021-10-20")
	// endDate, _ := time.Parse("2006-01-02", "2021-10-25")

	// scheduleScraper, _ := scraper.CreateScheduleScraperWithDates(c, "2022", startDate, endDate)

	// scheduleScraper.Scrape()
	// fmt.Println(scheduleScraper.GetData())
	// fmt.Println(scheduleScraper.GetChildUrls())

	// teamScraper := core.Factory[scraper.TeamScraper](c.Injector())
	// teamScraper.Scrape("TOR", "CHI")

	// playerScraper := core.Factory[scraper.PlayerScraper](c.Injector())
	// playerScraper.Scrape("vandeja01", "curryst01")

	// standingsScraper := core.Factory[scraper.StandingScraper](c.Injector())
	// standingsScraper.Scrape()

	// injuriesScraper := core.Factory[scraper.InjuryScraper](c.Injector())
	// injuriesScraper.Scrape()
}
