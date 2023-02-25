package scraper

import (
	"fmt"
)

type ScraperError struct {
	err error
	url string
}

func (s *ScraperError) Error() string {
	return fmt.Sprintf("error while scraping %s: %v", s.url, s.err)
}

func NewScraperError(err error, url string) *ScraperError {
	return &ScraperError{err: err, url: url}
}
