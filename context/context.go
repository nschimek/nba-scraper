package context

import "github.com/sirupsen/logrus"

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

	ctx.initialize()

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

	i.AddInjectable(createConfig())
	i.AddInjectable(createColly())
	i.AddInjectable(createDatabase())

	return i
}

func (c *context) Injector() *Injector {
	return c.injector
}

func (c *context) initialize() {
	// connect to database
	db := Factory[Database](c.injector)
	db.Connect()

	// set log level if Debug mode is enabled
	cfg := Factory[Config](c.injector)
	if cfg.Debug {
		Log.SetLevel(logrus.DebugLevel)
		Log.Info("Debug logging enabled!")
	}
}
