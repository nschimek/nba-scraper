package scraper

import "github.com/gocolly/colly/v2"

type TableParser struct {
	columnMaps []map[string]*colly.HTMLElement
}

func ParseTable(tbl *colly.HTMLElement) []map[string]*colly.HTMLElement {
	p := TableParser{
		columnMaps: make([]map[string]*colly.HTMLElement, 0),
	}

	tbl.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
		p.parseRow(tr)
	})

	return p.columnMaps
}

func (t *TableParser) parseRow(tr *colly.HTMLElement) {
	if tr.Attr("class") != "thead" { // exclude table headers (these are someties in the middle of the table)
		columnMap := make(map[string]*colly.HTMLElement)
		tr.ForEach("th", func(_ int, th *colly.HTMLElement) {
			t.parseColumn(columnMap, th)
		})

		tr.ForEach("td", func(_ int, td *colly.HTMLElement) {
			t.parseColumn(columnMap, td)
		})

		t.columnMaps = append(t.columnMaps, columnMap)
	}
}

func (t *TableParser) parseColumn(columnMap map[string]*colly.HTMLElement, td *colly.HTMLElement) {
	if dataStat := td.Attr("data-stat"); dataStat != "" {
		columnMap[dataStat] = td
	}
}
