package core

import (
	"strings"

	"github.com/spf13/viper"
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

func SetupViper() {
	viper.SetDefault("use-config-file", true) // overrideable with environment variables only
	viper.SetEnvPrefix("ns")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
}

func createConfig(configFile string) *Config {
	if viper.GetBool("use-config-file") == true {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err == nil {
			Log.Infof("Loaded config file: %s", viper.ConfigFileUsed())
		} else {
			Log.Fatalf("Could not load config file: %s!", configFile)
		}
	} else {
		Log.Info("Config file NOT being used...requiring NS_ENVIRONMENT_VARIABLES")
		bindViperEnvVars()
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		Log.Fatalf("Error decoding Config struct: %v", err)
	}

	if config.Season < 1947 { // unspecified is 0 and BR goes back to 1946-47 season!
		Log.Fatalf("No or invalid Season specified!  Please specify a valid Season in YYYY format.  Use the finishing year (2021-22 would be 2022)")
	} else {
		Log.Infof("Season set to: %d", config.Season)
	}

	return config
}

// viper needs a little help with these nested variables...
func bindViperEnvVars() {
	viper.BindEnv("suppression.team")
	viper.BindEnv("suppression.player")
	viper.BindEnv("database.user")
	viper.BindEnv("database.password")
	viper.BindEnv("database.location")
	viper.BindEnv("database.port")
	viper.BindEnv("database.name")
}
