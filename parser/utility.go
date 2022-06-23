package parser

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// reusable parser utilities

var EST, _ = time.LoadLocation("America/New_York")
var CST, _ = time.LoadLocation("America/Chicago")

func parseLink(e *colly.HTMLElement) string {
	return e.ChildAttr("a", "href")
}

// for URLs where the last part of the URL (*.html)
func ParseLastId(link string) string {
	s := strings.Split(link, "/")
	return strings.Replace(s[len(s)-1], ".html", "", 1)
}

func ParseTeamId(link string) string {
	s := strings.Split(link, "/")
	return s[len(s)-2]
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

// simply remove the $ and , from the currency string and parse as an int
func parseCurrency(s string) (int, error) {
	return strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(s, ",", ""), "$", ""))
}

// removes newlines and strips extra whitespace from a string
func removeNewlines(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\n", "")
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

func RegexParamMap(regEx, target string) (rpm map[string]string) {
	r := regexp.MustCompile(regEx)
	m := r.FindStringSubmatch(target)

	rpm = make(map[string]string)
	for i, n := range r.SubexpNames() {
		if i > 0 && i < len(m) {
			rpm[n] = m[i]
		}
	}

	return
}
