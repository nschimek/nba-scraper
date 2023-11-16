package parser

type ParserError struct {
	Id     string
	Errors []error
}

func (pe *ParserError) Init(Id string) {
	pe.Id = Id
	pe.Errors = make([]error, 0)
}
