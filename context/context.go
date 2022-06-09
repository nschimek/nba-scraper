package context

type context struct {
	injector *Injector
}

var ctx *context

func Setup() *context {
	if ctx != nil {
		Log.Fatal("Context already setup")
	}

	ctx = &context{
		injector: setupInjector(),
	}

	return ctx
}

func Get() *context {
	if ctx == nil {
		Log.Fatal("Context not setup")
		return nil
	} else {
		return ctx
	}
}

func setupInjector() *Injector {
	i := createInjector()

	i.AddInjectable(createColly())
	// i.AddInjectable(connectToDatabase())

	return i
}

func (c *context) Injector() *Injector {
	return c.injector
}
