package config

import (
	"github.com/charmbracelet/log"
	"github.com/goccy/go-yaml"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func ReadFile(fs afero.Fs, file string) (*Config, error) {
	log.Debug("Reading config", "file", file)
	data, err := afero.ReadFile(fs, file)
	if err != nil {
		log.Debug("ReadFile", "file", file, "err", err)
		return Default, nil
	}

	var c Config
	if err = Unmarshal(data, &c); err != nil {
		return nil, err
	} else {
		return &c, nil
	}
}

func Unmarshal(data []byte, config *Config) error {
	return yaml.Unmarshal(data, config)
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
