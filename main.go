package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.basketball-reference.com"),
	)

	c.OnHTML("body #wrap #content .game_summaries", func(e *colly.HTMLElement) {
		e.ForEach(".game_summary", func(_ int, g *colly.HTMLElement) {
			gameId := g.ChildAttr("table tbody tr:first-child td.gamelink a", "href")
			fmt.Println(gameId)
		})
	})

	c.Visit("https://www.basketball-reference.com/boxscores/?month=03&day=1&year=2022")
}
