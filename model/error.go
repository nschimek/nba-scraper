package model

import "github.com/nschimek/nba-scraper/core"

type ModelError struct {
	Errors []error `gorm:"-"` // ignore this field in persistence
}

func (m *ModelError) CaptureError(errs ...error) {
	for _, err := range errs {
		if err != nil {
			m.Errors = append(m.Errors, err)
		}
	}
}

func (m *ModelError) HasErrors() bool {
	return len(m.Errors) > 0
}

func (m *ModelError) LogErrors(prefix string) {
	core.Log.Errorf("%s has the following critical parsing errors and will not be persisted:", prefix)
	for _, err := range m.Errors {
		core.Log.Errorf(" - %s", err.Error())
	}
}
