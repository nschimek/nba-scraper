package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata" // resolve timezone database issues for windows binary on systems without Go installed

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
)

// reusable parser utilities

var EST, _ = time.LoadLocation("America/New_York")
var CST, _ = time.LoadLocation("America/Chicago")

func parseLink(e *colly.HTMLElement) string {
	if e != nil {
		return e.ChildAttr("a", "href")
	} else {
		core.Log.Error("Could not parse link!")
		return ""
	}
}

// for URLs where we want last part of the URL (*.html)
func ParseLastId(link string) (string, error) {
	s := strings.Split(link, "/")

	if len(s) == 0 {
		return "", errors.New("link not in expected format for getting last ID")
	}

	return strings.Replace(s[len(s)-1], ".html", "", 1), nil
}

func ParseTeamId(link string) (string, error) {
	s := strings.Split(link, "/")

	if len(s) != 4 {
		return "", errors.New("team link not in expected format")
	}

	if id := s[len(s)-2]; len(id) != 3 {
		return "", errors.New("team ID not in expected format")
	} else {
		return id, nil
	}
}

func parseDuration(duration string) (int, error) {
	// durations are in string format of m:s, so convert them into #m#s format for time.ParseDuration()
	d, err := time.ParseDuration(strings.Replace(duration, ":", "m", 1) + "s")
	return int(d.Seconds()), err
}

// when players do not attempt the underlying stat that generates a float, the site returns a blank - convert that to a zero
func parseFloatStat(s string) (float64, error) {
	if s != "" {
		return strconv.ParseFloat(s, 64)
	} else {
		return 0, nil
	}
}

// removes newlines and strips extra whitespace from a string
func removeNewlines(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\n", "")
}

// Given a regular expression with named capture group(s) [P<name> in Golang],
// this will store the match results in a map where the key is equal to the capture name.
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

func getColumnText(rowMap map[string]*colly.HTMLElement, column string) string {
	if col, ok := rowMap[column]; ok {
		return col.Text
	} else {
		core.Log.Warnf("Could not get expected column '%s'! Will be default value.", column)
		return ""
	}
}
