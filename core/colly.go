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
	RandomDelay: 5 * time.Second,
}

func createColly() *colly.Collector {
	c := colly.NewCollector(colly.AllowedDomains(allowedDomain))
	c.Limit(LimitRule)
	return c
}

func CloneColly(colly *colly.Collector) *colly.Collector {
	c := colly.Clone()
	c.OnRequest(onRequestVisit)
	c.OnError(onError)
	return c
}

func onRequestVisit(r *colly.Request) {
	Log.Infof("Visiting: %s", r.URL.String())
}

func onError(r *colly.Response, err error) {
	Log.Info(string(r.Body))
	Log.Fatalf("Scraping resulted in error: %s", err)
}
