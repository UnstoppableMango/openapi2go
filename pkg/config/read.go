package config

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	openapi2go "github.com/unstoppablemango/openapi2go/pkg"
)

func Read() (openapi2go.Config, error) {
	v := Viper("")

	if err := v.ReadInConfig(); err == nil {
		log.Debug("Using config file", "path", v.ConfigFileUsed())
	}

	return Parse(v)
}

func Parse(v *viper.Viper) (openapi2go.Config, error) {
	config := openapi2go.Config{}
	if err := Unmarshal(v, &config); err != nil {
		return openapi2go.Config{}, err
	} else {
		return config, nil
	}
}

func Unmarshal(v *viper.Viper, config *openapi2go.Config) error {
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
