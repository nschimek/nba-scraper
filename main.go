package main

import (
	"fmt"
)

type Test struct {
	Child *Child `Inject:"self"`
	name  string
}

type Child struct {
	age int
}

func (t *Test) setName(name string) {
	t.name = name
}

func (t *Test) getName() string {
	return t.name
}

func main() {
	// c := SetupContext()
	i := CreateInjector()

	child := InjectorFactory[Child](i)
	child.age = 5

	test := InjectorFactory[Test](i)
	test.setName("Nick")

	fmt.Println(test.getName())

	test2 := InjectorFactory[Test](i)
	fmt.Println(test2.getName())
	fmt.Println(test2.Child.age)

	// startDate, _ := time.Parse("2006-01-02", "2021-10-20")
	// endDate, _ := time.Parse("2006-01-02", "2021-10-25")

	// scheduleScraper, _ := scraper.CreateScheduleScraperWithDates(c, "2022", startDate, endDate)

	// scheduleScraper.Scrape()
	// fmt.Println(scheduleScraper.GetData())
	// fmt.Println(scheduleScraper.GetChildUrls())

	// gameScraper := scraper.CreateGameScraper(c)
	// gameScraper.Scrape("https://www.basketball-reference.com/boxscores/202110300WAS.html", "https://www.basketball-reference.com/boxscores/202204180GSW.html")

	// teamScraper := scraper.CreateTeamScraper(c)
	// teamScraper.Scrape("https://www.basketball-reference.com/teams/TOR/2022.html")

	// playerScraper := scraper.CreatePlayerScraper(c)
	// playerScraper.Scrape("https://www.basketball-reference.com/players/v/vandeja01.html")

	// standingsScraper := scraper.CreateStandingsScraper(c, 2022)
	// standingsScraper.Scrape()

	// injuriesScraper := scraper.CreateInjuriesScraper(c, 2022)
	// injuriesScraper.Scrape()
}
