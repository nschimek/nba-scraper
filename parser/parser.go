package parser

import "github.com/nschimek/nba-scraper/core"

type ParserError struct {
	Id     string
	Errors []error
}

func (pe *ParserError) Init(Id string) {
	pe.Id = Id
	pe.Errors = make([]error, 0)
}

func (pe *ParserError) capture(err error) {
	if err != nil {
		pe.Errors = append(pe.Errors, err)
		core.Log.WithField("id", pe.Id).Errorf("parsing error: %s", err.Error())
	}
}
