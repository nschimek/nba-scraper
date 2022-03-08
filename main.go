package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	BaseUrl       = "https://www.basketball-reference.com"
	AllowedDomain = "www.basketball-reference.com"
	Season        = "2022"
	Month         = "october"
)

var LimitRule = colly.LimitRule{
	Parallelism: 2,
	RandomDelay: 5 * time.Second,
}

type Schedule struct {
	StartTime time.Time
	VisitorId string
	HomeId    string
	GameId    string
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains(AllowedDomain),
	)

	c.Limit(&LimitRule)

	c.OnHTML("body #wrap #content #all_schedule #div_schedule table tbody", func(tbl *colly.HTMLElement) {
		tbl.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
			if tr.Attr("class") != "thead" {
				schedule := Schedule{}
				tr.ForEach("td", func(_ int, td *colly.HTMLElement) {
					switch td.Attr("data-stat") {
					case "home_team_name":
						schedule.HomeId = strings.Split(td.ChildAttr("a", "href"), "/")[2]
					case "visitor_team_name":
						schedule.VisitorId = strings.Split(td.ChildAttr("a", "href"), "/")[2]
					case "box_score_text":
						schedule.GameId = strings.Replace(strings.Split(td.ChildAttr("a", "href"), "/")[2], ".html", "", 1)
					}
				})
				fmt.Println(schedule)
			}
		})
	})

	c.Visit(BaseUrl + "/leagues/NBA_" + Season + "_games-" + Month + ".html")
}
