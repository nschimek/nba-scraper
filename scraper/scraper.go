package scraper

import (
	"bytes"
	"errors"
	"html"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
)

type PeristableScraper interface {
	Persist()
}

var exists = struct{}{}

func transformHtmlElement(element *colly.HTMLElement, query string, transform func(html string) string) (*colly.HTMLElement, error) {
	h, _ := element.DOM.Html()
	// call the transform function passed in on the html string, but first unescape - the new version of Colly seems to require this
	doc, _ := goquery.NewDocumentFromReader(bytes.NewBufferString(transform(html.UnescapeString(h))))
	sel := doc.Find(query)

	if len(sel.Nodes) == 0 || sel == nil {
		return nil, errors.New("could not find any search elements in transformed table")
	}

	return colly.NewHTMLElementFromSelectionNode(element.Response, sel, sel.Get(0), 0), nil
}

func removeCommentsSyntax(html string) string {
	return strings.ReplaceAll(strings.ReplaceAll(html, "<!--", ""), "-->", "")
}

func persistIfPopulated[T any](persist func(scrapedData []T, label string), scrapedData []T, label string) {
	if len(scrapedData) > 0 {
		persist(scrapedData, label)
	} else {
		core.Log.Warnf("No %s(s) scraped to persist, skipping!", label)
	}
}
