package scraper

import (
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	BaseHttp        = "https://www.basketball-reference.com"
	AllowedDomain   = "www.basketball-reference.com"
	baseLeaguesPath = "leagues"
)

var LimitRule = colly.LimitRule{
	Parallelism: 2,
	RandomDelay: 5 * time.Second,
}
