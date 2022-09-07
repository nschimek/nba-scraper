package core

import (
	"gopkg.in/ini.v1"
)

const (
	dir  = "conf"
	ext  = "ini"
	core = "core"
)

var iniOptions = ini.LoadOptions{IgnoreInlineComment: true}

type Config struct {
	Season                                     int
	Environment                                string
	Debug                                      bool
	TeamSuppressionDays, PlayerSuppressionDays int
	Database                                   database
}

type database struct {
	User, Password, Location, Name string
}

func createConfig() *Config {
	cfg := new(Config)

	coreName := getFullName(core)
	cfg.loadAndMap(coreName, "")

	if cfg.Environment == "" {
		Log.Fatal("Could not determine Environment from core INI")
	}

	cfg.loadAndMap(coreName, getFullName(cfg.Environment))

	return cfg
}

func (c *Config) loadAndMap(core string, env string) {
	i := c.loadFromIni(core, env)
	c.mapFromIni(i)
}

func (c *Config) loadFromIni(core string, env string) (i *ini.File) {
	var err error

	if env == "" {
		Log.WithField("file", core).Info("Loading core config from INI...")
		i, err = ini.LoadSources(iniOptions, core)
	} else {
		Log.WithField("file", env).Info("Loading environmental config from INI...")
		i, err = ini.LoadSources(iniOptions, core, env)
	}

	if err != nil {
		Log.Fatal(err)
		return nil
	} else {
		return i
	}
}

func (c *Config) mapFromIni(i *ini.File) {
	err := i.MapTo(c)

	if err != nil {
		Log.Fatal(err)
	}
}

func getFullName(n string) string {
	return dir + "/" + n + "." + ext
}
