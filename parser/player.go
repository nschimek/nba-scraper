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
	bornDateRegex       = `Born:\s+(?P<birthMonth>\w+)\s+(?P<birthDay>\d{1,2}),\s+(?P<birthYear>\d{4})\s+in.(?P<birthPlace>.+)\s\s(?P<birthCountry>[a-z]{2})$`
)

type regexParser struct {
	regex  string
	parser func(m *model.Player, rm map[string]string)
}

func parseShootsPosition(p *model.Player, rm map[string]string) {
	p.Shoots = strings.ToUpper(string(rm["shoots"][0]))
}

func parseHeightWeight(m *model.Player, rm map[string]string) {
	ft, _ := strconv.Atoi(rm["ft"])
	in, _ := strconv.Atoi(rm["in"])
	m.Height = (ft * 12) + in
	m.Weight, _ = strconv.Atoi(rm["lb"])
}

func parseBirthInfo(p *model.Player, rm map[string]string) {
	p.BirthDate, _ = time.ParseInLocation("January 2 2006", rm["birthMonth"]+" "+rm["birthDay"]+" "+rm["birthYear"], CST)
	p.BirthPlace = strings.TrimSpace(rm["birthPlace"])
	p.BirthCountryCode = strings.ToUpper(rm["birthCountry"])
}

var regexParsers = [...]regexParser{
	{regex: positionShootsRegex, parser: parseShootsPosition},
	{regex: heightWeightRegex, parser: parseHeightWeight},
	{regex: bornDateRegex, parser: parseBirthInfo},
}

func (p *PlayerParser) PlayerInfoBox(m *model.Player, div *colly.HTMLElement) {
	m.Name = div.ChildText("h1")

	crp := 0 // current regex parser starts at 0

	div.ForEach("p", func(i int, e *colly.HTMLElement) {
		line := strings.TrimSpace(removeNewlines(e.Text))
		if crp < len(regexParsers) {
			rm := RegexParamMap(regexParsers[crp].regex, line)
			// if there's a hit with this regex, we want to run the parser function and increment cra
			if len(rm) > 0 {
				regexParsers[crp].parser(m, rm)
				crp++
			}
			// the regex list is in the same order as the elements are listed on the page... so if we do not get a hit,
			// we simply do nothing and move on to the next line.  the next matching line will come soon enough.
		}
	})
}
