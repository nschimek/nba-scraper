package context

import (
	"gopkg.in/ini.v1"
)

const (
	dir  = "conf"
	ext  = "ini"
	core = "core"
)

type Config struct {
	Environment, BaseHttp string
	Database              database
}

type database struct {
	User, Password, Host, Name string
}

func CreateConfig() *Config {
	cfg := new(Config)

	coreName := getFullName(core)
	cfg.loadAndMap(coreName)

	if cfg.Environment == "" {
		Log.Fatal("Could not determine Environment from core INI")
	}

	cfg.loadAndMap(coreName, getFullName(cfg.Environment))

	return cfg
}

func (c *Config) loadAndMap(core string, others ...string) {
	var i *ini.File
	var err error

	if len(others) > 0 {
		i, err = ini.Load(core, others)
	} else {
		i, err = ini.Load(core)
	}

	if err != nil {
		Log.Fatal(err)
	}

	err = i.MapTo(c)

	if err != nil {
		Log.Fatal(err)
	}
}

func getFullName(n string) string {
	return dir + "/" + n + "." + ext
}
