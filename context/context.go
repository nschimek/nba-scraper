package context

import "gopkg.in/ini.v1"

type context struct {
	injector *Injector
}

var ctx *context
var Config *ini.File

func Setup(env string) *context {
	if ctx != nil {
		Log.Fatal("Context already setup")
	}

	ctx = &context{
		injector: setupInjector(),
	}

	// Create global config instance with the environment variable passed in
	Config = createConfig(env)

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
