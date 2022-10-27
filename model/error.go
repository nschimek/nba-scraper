package model

type ModelError struct {
	Errors []error `gorm:"-"` // ignore this field in persistence
}

func (m *ModelError) CaptureError(err ...error) {
	if err != nil {
		m.Errors = append(m.Errors, err...)
	}
}

func (m *ModelError) HasErrors() bool {
	return len(m.Errors) > 0
}
