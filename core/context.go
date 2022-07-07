package core

import "github.com/sirupsen/logrus"

type coreContext struct {
	injector *Injector
}

var ctx *coreContext

func Setup() *coreContext {
	if ctx != nil {
		Log.Fatal("Core Context already setup")
	}

	ctx = &coreContext{
		injector: setupInjector(),
	}

	ctx.initialize()

	return ctx
}

func Get() *coreContext {
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

func (c *coreContext) Injector() *Injector {
	return c.injector
}

func (c *coreContext) initialize() {
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
