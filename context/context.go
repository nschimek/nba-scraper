package context

type context struct {
	injector *Injector
}

var ctx *context

func Setup() *context {
	if ctx != nil {
		return ctx
	}

	ctx = &context{
		injector: setupInjector(),
	}

	return ctx
}

func Get() *context {
	if ctx == nil {
		return Setup()
	} else {
		return ctx
	}
}

func setupInjector() *Injector {
	i := createInjector()

	i.AddInjectable(createColly())

	return i
}

func (c *context) Injector() *Injector {
	return c.injector
}
