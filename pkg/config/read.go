package config

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func Read() (Config, error) {
	v := Viper("")

	if err := v.ReadInConfig(); err == nil {
		log.Debug("Using config file", "path", v.ConfigFileUsed())
	}

	return Parse(v)
}

func Parse(v *viper.Viper) (Config, error) {
	config := Config{}
	if err := Unmarshal(v, &config); err != nil {
		return Config{}, err
	} else {
		return config, nil
	}
}

func Unmarshal(v *viper.Viper, config *Config) error {
	return viper.Unmarshal(config)
}

func Viper(configFile string) *viper.Viper {
	v := viper.New()
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.AddConfigPath(".")
		v.SetConfigName(".openapi2go")
		v.SetConfigType("yaml")
	}

	return v
}
