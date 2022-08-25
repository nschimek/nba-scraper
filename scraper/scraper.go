package scraper

import (
	"bytes"
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

var exists = struct{}{}

func IdMapToArray(idMap map[string]bool) (ids []string) {
	for id, keep := range idMap {
		if keep {
			ids = append(ids, id)
		}
	}
	return
}

func ConsolidateIdMaps(idMaps ...map[string]bool) (idMap map[string]bool) {
	for _, m := range idMaps {
		for k, v := range m {
			idMap[k] = v
		}
	}

	return idMap
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
