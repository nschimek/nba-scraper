package context

import (
	"gopkg.in/ini.v1"
)

var dir = "env"

func createConfig(env string) *ini.File {
	f := dir + "/" + env + ".ini"
	ini, err := ini.Load(f)

	if err != nil {
		Log.WithField("file", f).Fatal("Could not find config INI file")
	}

	return ini
}
