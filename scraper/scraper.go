package scraper

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Scraper interface {
	Scrape(urls ...string)
	GetData() interface{}
	AttachChild(scraper *Scraper)
	GetChild() Scraper
	GetChildUrls() []string
}

func scrapeChild(s Scraper) {
	if s.GetChild() != nil && len(s.GetChildUrls()) > 0 {
		s.GetChild().Scrape(s.GetChildUrls()...)
	}
}

func urlsMapToArray(urlsMap map[string]string) (urlsArray []string) {
	for _, url := range urlsMap {
		urlsArray = append(urlsArray, url)
	}
	return
}

func parseLink(e *colly.HTMLElement) string {
	return e.ChildAttr("a", "href")
}

func parseGameId(link string) string {
	return strings.Replace(strings.Split(link, "/")[2], ".html", "", 1)
}

func parseTeamId(link string) string {
	return strings.Split(link, "/")[2]
}

func parsePlayerId(link string) string {
	return strings.Replace(strings.Split(link, "/")[3], ".html", "", 1)
}

func parseDuration(duration string) (time.Duration, error) {
	// durations are in string format of m:s, so convert them into #m#s format for time.ParseDuration()
	return time.ParseDuration(strings.Replace(duration, ":", "m", 1) + "s")
}

// when players do not attempt the underlying stat that generates a float, the site returns a blank - convert that to a zero
func parseFloatStat(s string) (float64, error) {
	if s != "" {
		return strconv.ParseFloat(s, 64)
	} else {
		return 0, nil
	}
}

func transformHtmlElement(element *colly.HTMLElement, query string, transform func(html string) string) (*colly.HTMLElement, error) {
	html, _ := element.DOM.Html()
	doc, _ := goquery.NewDocumentFromReader(bytes.NewBufferString(transform(html)))
	sel := doc.Find(query)

	if len(sel.Nodes) == 0 {
		return nil, errors.New("could not find any search elements in transformed table")
	}

	return colly.NewHTMLElementFromSelectionNode(element.Response, sel, sel.Get(0), 0), nil
}

func removeCommentsSyntax(html string) string {
	return strings.ReplaceAll(strings.ReplaceAll(html, "<!--", ""), "-->", "")
}
