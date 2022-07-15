package main

import (
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/scraper"
)

func main() {
	c := core.Setup()

	// gameScraper := context.Factory[scraper.GameScraper](c.Injector())
	// gameScraper.Scrape("https://www.basketball-reference.com/boxscores/202110300WAS.html", "https://www.basketball-reference.com/boxscores/202204180GSW.html")

	// startDate, _ := time.Parse("2006-01-02", "2021-10-20")
	// endDate, _ := time.Parse("2006-01-02", "2021-10-25")

	// scheduleScraper, _ := scraper.CreateScheduleScraperWithDates(c, "2022", startDate, endDate)

	// scheduleScraper.Scrape()
	// fmt.Println(scheduleScraper.GetData())
	// fmt.Println(scheduleScraper.GetChildUrls())

	teamScraper := core.Factory[scraper.TeamScraper](c.Injector())
	teamScraper.Scrape("TOR", "CHI")
	// teamScraper.Scrape("https://www.basketball-reference.com/teams/TOR/2022.html")

	// playerScraper := core.Factory[scraper.PlayerScraper](c.Injector())
	// playerScraper.Scrape("vandeja01", "curryst01")

	// standingsScraper := scraper.CreateStandingsScraper(c, 2022)
	// standingsScraper.Scrape()

	// injuriesScraper := scraper.CreateInjuriesScraper(c, 2022)
	// injuriesScraper.Scrape()
}
