package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/parser"
)

const (
	basePlayerBodyElement = "body > div#wrap"
	playerInfoElement     = basePlayerBodyElement + " > div#info > div#meta > div:nth-child(2)"
	positionShootsRegex   = `Position:\s+(?P<position>\w+?\s\w*).+Shoots:\s+(?P<shoots>\w+)`
	heightWeightRegex     = `(?P<ft>\d{1,2})-(?P<in>\d{1,2}),.(?P<lb>\d{2,3})lb`
	birthdayRegex         = `Born:\s+(?P<month>\w+)\s(?P<day>\d{1,2}),\s+(?P<year>\d{4})`
)

type PlayerScraper struct {
	colly       colly.Collector
	season      int
	ScrapedData []parser.Player
	Errors      []error
	child       Scraper
	childUrls   map[string]string
}

func CreatePlayerScraper(c *colly.Collector, season int) PlayerScraper {
	return PlayerScraper{
		colly:     *c,
		season:    season,
		childUrls: make(map[string]string),
	}
}

func (s *PlayerScraper) GetData() interface{} {
	return s.ScrapedData
}

func (s *PlayerScraper) AttachChild(scraper *Scraper) {
	s.child = *scraper
}

func (s *PlayerScraper) GetChild() Scraper {
	return s.child
}

func (s *PlayerScraper) GetChildUrls() []string {
	return urlsMapToArray(s.childUrls)
}

func (s *PlayerScraper) Scrape(urls ...string) {

	for _, url := range urls {
		player := s.parsePlayerPage(url)
		s.ScrapedData = append(s.ScrapedData, player)
	}

	// fmt.Printf("%+v\n", s.ScrapedData)

	scrapeChild(s)
}

func (s *PlayerScraper) parsePlayerPage(url string) (player parser.Player) {
	c := s.colly.Clone()

	c.OnRequest(onRequestVisit)

	player.Id = parser.ParseLastId(url)

	c.OnHTML(playerInfoElement, func(div *colly.HTMLElement) {
		player.Name = strings.Split(div.ChildText("h1"), " 20")[0]

		div.ForEach("p", func(i int, p *colly.HTMLElement) {
			line := strings.ReplaceAll(strings.TrimSpace(p.Text), "\n", "")

			if i == 2 {
				regexMap := parser.RegexParamMap(positionShootsRegex, line)
				player.Shoots = regexMap["shoots"]
				player.Position = regexMap["position"]
			} else if i == 3 {
				regexMap := parser.RegexParamMap(heightWeightRegex, line)
				ft, _ := strconv.Atoi(regexMap["ft"])
				in, _ := strconv.Atoi(regexMap["in"])
				player.Height = (ft * 12) + in
				player.Weight, _ = strconv.Atoi(regexMap["lb"])
			} else if i == 4 {
				p.ForEach("span", func(_ int, s *colly.HTMLElement) {
					if s.Attr("itemprop") == "birthDate" {
						fmt.Println(s.Attr("data-birth"))
					}
					if s.Attr("itemprop") == "birthPlace" {
						fmt.Println(strings.ReplaceAll(strings.TrimSpace(s.Text), "\n", ""))
					}
					if strings.Contains(s.Attr("class"), "f-i") {
						fmt.Println(strings.ToUpper(s.Text))
					}
				})
			} else if i == 8 {
				fmt.Println(p.DOM.Text()) // Draft
			} else if i == 9 {
				fmt.Println(p.DOM.Text()) // NBA Debut
			}
		})

		fmt.Printf("%+v\n", player)
	})

	c.Visit(strings.Replace(url, ".html", "", 1) + "/gamelog/" + strconv.Itoa(s.season))

	return
}
