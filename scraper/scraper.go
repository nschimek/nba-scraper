package scraper

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Scraper interface {
	Scrape(urls ...string)
	GetData() interface{}
	AttachChild(scraper *Scraper)
	GetChild() *Scraper
	GetChildUrls() []string
}

func onRequestVisit(r *colly.Request) {
	fmt.Println("Visiting: ", r.URL.String())
}

func scrapeChild(s Scraper) {
	if s.GetChild() != nil && len(s.GetChildUrls()) > 0 {
		c := *s.GetChild()
		c.Scrape(s.GetChildUrls()...)
	}
}

func urlsMapToArray(urlsMap map[string]string) (urlsArray []string) {
	for _, url := range urlsMap {
		urlsArray = append(urlsArray, url)
	}
	return
}

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
