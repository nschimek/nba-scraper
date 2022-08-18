package core

import (
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	allowedDomain = "www.basketball-reference.com"
	domainGlob    = "*" + allowedDomain + "*"
)

var LimitRule = &colly.LimitRule{
	DomainGlob:  domainGlob,
	Parallelism: 1,
	RandomDelay: 3 * time.Second,
}

func createColly() *colly.Collector {
	c := colly.NewCollector(colly.AllowedDomains(allowedDomain))
	c.Limit(LimitRule)
	return c
}
