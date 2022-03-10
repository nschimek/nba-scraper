package scraper

import (
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	BaseHttp = "https://www.basketball-reference.com"
)

var LimitRule = colly.LimitRule{
	Parallelism: 2,
	RandomDelay: 5 * time.Second,
}

var EST, _ = time.LoadLocation("America/New_York")
