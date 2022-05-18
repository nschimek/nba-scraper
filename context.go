package main

import (
	"reflect"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/scraper"
)

const (
	GameScraper = "GameScraper"
)

type Context struct {
	colly    *colly.Collector
	scrapers map[reflect.Type]scraper.Scraper
}

func SetupContext() *Context {
	c := colly.NewCollector(colly.AllowedDomains(scraper.AllowedDomain))
	c.Limit(&scraper.LimitRule)

	m := make(map[reflect.Type]scraper.Scraper)

	s := scraper.CreateGameScraper(c)
	m[reflect.TypeOf(s)] = s

	return &Context{
		colly:    c,
		scrapers: m,
	}
}

func ScraperFactory[T scraper.Scraper](c *Context) T {
	return c.scrapers[typeOf[T]()].(T)
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
