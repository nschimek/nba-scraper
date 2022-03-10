package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/scraper"
)

const (
	AllowedDomain = "www.basketball-reference.com"
	Season        = "2022"
	Month         = "october"
)

var LimitRule = colly.LimitRule{
	Parallelism: 2,
	RandomDelay: 5 * time.Second,
}

func main() {
	startDate, _ := time.Parse("2006-01-02", "2021-10-20")
	endDate, _ := time.Parse("2006-01-02", "2021-10-25")
	fmt.Println(scraper.Schedule("2022", startDate, endDate))
}
