package scraper

import (
	"bytes"
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/nschimek/nba-scraper/core"
	"github.com/nschimek/nba-scraper/model"
)

func handleModelErrors(label string, id string, model model.ModelError) {
	core.Log.Errorf("%s %s has the following critical parsing errors:", label, id)
	model.LogErrors()
}

type PeristableScraper interface {
	Persist()
}

var exists = struct{}{}

func transformHtmlElement(element *colly.HTMLElement, query string, transform func(html string) string) (*colly.HTMLElement, error) {
	html, _ := element.DOM.Html()
	doc, _ := goquery.NewDocumentFromReader(bytes.NewBufferString(transform(html)))
	sel := doc.Find(query)

	if len(sel.Nodes) == 0 {
		return nil, errors.New("could not find any search elements in transformed table")
	}

	return colly.NewHTMLElementFromSelectionNode(element.Response, sel, sel.Get(0), 0), nil
}

func removeCommentsSyntax(html string) string {
	return strings.ReplaceAll(strings.ReplaceAll(html, "<!--", ""), "-->", "")
}

func persistIfPopulated[T any](persist func(scrapedData []T), scrapedData []T, label string) {
	if len(scrapedData) > 0 {
		persist(scrapedData)
	} else {
		core.Log.Warnf("No %s(s) scraped to persist, skipping!")
	}
}
