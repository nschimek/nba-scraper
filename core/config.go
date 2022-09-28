package core

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	dir  = "conf"
	ext  = "ini"
	core = "core"
)

type Config struct {
	Season        int
	Debug         bool
	UseConfigFile bool `mapstructure:"use-config-file"`
	Suppression   struct {
		Team   int
		Player int
	}
	Database struct {
		User, Password, Location, Name string
		Port                           int
	}
}

func createConfig(configFile string) *Config {
	viper.SetDefault("use-config-file", true) // overrideable with environment variables only
	viper.SetEnvPrefix("ns")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if viper.GetBool("use-config-file") == true {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err == nil {
			Log.Infof("Loaded config file: %s", viper.ConfigFileUsed())
		} else {
			Log.Fatalf("Could not load config file: %s!", configFile)
		}
	} else {
		Log.Info("Config file NOT being used...requiring NS_ENVIRONMENT_VARIABLES")
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		Log.Fatalf("Error decoding Config struct: %v", err)
	}
	return config
}
