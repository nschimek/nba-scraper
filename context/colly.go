package context

import (
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	allowedDomain = "www.basketball-reference.com"
)

var limitRule = &colly.LimitRule{
	Parallelism: 2,
	RandomDelay: 5 * time.Second,
}

func createColly() *colly.Collector {
	c := colly.NewCollector(colly.AllowedDomains(allowedDomain))
	c.Limit(limitRule)
	return c
}
