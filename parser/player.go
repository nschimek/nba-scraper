package parser

import (
	"fmt"
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

type regexAssign struct {
	regex    string
	assigner func(m *model.Player, rm map[string]string)
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
	fmt.Println(rm)
	p.BirthDate, _ = time.Parse("January 2 2006", rm["birthMonth"]+" "+rm["birthDay"]+" "+rm["birthYear"])
	p.BirthPlace = strings.TrimSpace(rm["birthPlace"])
	p.BirthCountryCode = strings.ToUpper(rm["birthCountry"])
}

var regexAssigners = [...]regexAssign{
	{regex: positionShootsRegex, assigner: parseShootsPosition},
	{regex: heightWeightRegex, assigner: parseHeightWeight},
	{regex: bornDateRegex, assigner: parseBirthInfo},
}

func (p *PlayerParser) PlayerInfoBox(m *model.Player, div *colly.HTMLElement) {
	m.Name = div.ChildText("h1")

	cra := 0 // current regex assigner starts at 0
	max := len(regexAssigners)

	div.ForEach("p", func(i int, e *colly.HTMLElement) {
		line := strings.TrimSpace(removeNewlines(e.Text))
		if cra < max {
			rm := RegexParamMap(regexAssigners[cra].regex, line)
			// if there's a hit with this regex, we want to run the assigner function and increment cra
			if len(rm) > 0 {
				regexAssigners[cra].assigner(m, rm)
				cra++
			}
			// the regex list is in the same order as the elements are listed on the page... so if we do not get a hit,
			// we simply do nothing and move on to the next line.  the next matching line will come soon enough.
		}
	})
}
