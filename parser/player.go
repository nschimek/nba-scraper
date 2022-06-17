package parser

import (
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/model"
)

type PlayerParser struct{}

const (
	positionShootsRegex = `Position:\s+(?P<position>\w+?\s\w*).+Shoots:\s+(?P<shoots>\w+)`
	heightWeightRegex   = `(?P<ft>\d{1,2})-(?P<in>\d{1,2}),.(?P<lb>\d{2,3})lb`
)

func (p *PlayerParser) PlayerInfoBox(m *model.Player, div *colly.HTMLElement) {
	m.Name = div.ChildText("h1")

	div.ForEach("p", func(i int, e *colly.HTMLElement) {
		line := removeNewlines(e.Text)

		switch i {
		case 2:
			m.Shoots, _ = p.parseShootsPosition(line) // position on this page is the wild west, so we ignore it
		case 3:
			m.Height, m.Weight = p.parseHeightWeight(line)
		case 4:
			m.BirthDate, m.BirthPlace, m.BirthCountryCode = p.parseBirthInfo(e)
		}
	})
}

func (p *PlayerParser) parseShootsPosition(s string) (shoots, position string) {
	regexMap := RegexParamMap(positionShootsRegex, s)
	shoots = strings.ToUpper(string(regexMap["shoots"][0]))
	position = regexMap["position"]
	return
}

func (p *PlayerParser) parseHeightWeight(s string) (height, weight int) {
	regexMap := RegexParamMap(heightWeightRegex, s)
	ft, _ := strconv.Atoi(regexMap["ft"])
	in, _ := strconv.Atoi(regexMap["in"])
	height = (ft * 12) + in
	weight, _ = strconv.Atoi(regexMap["lb"])
	return
}

func (p *PlayerParser) parseBirthInfo(e *colly.HTMLElement) (birthDate time.Time, birthPlace, birthCountryCode string) {
	e.ForEach("span", func(_ int, s *colly.HTMLElement) {
		if s.Attr("itemprop") == "birthDate" {
			birthDate, _ = time.Parse("2006-01-02", s.Attr("data-birth"))
		} else if s.Attr("itemprop") == "birthPlace" {
			birthPlace = strings.TrimSpace(strings.TrimPrefix(removeNewlines(s.Text), "in"))
		} else if strings.Contains(s.Attr("class"), "f-i") {
			birthCountryCode = strings.ToUpper(s.Text)
		}
	})
	return
}
