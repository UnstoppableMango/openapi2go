package config

import (
	"github.com/goccy/go-yaml"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func Read(fs afero.Fs, file string) (c Config, err error) {
	data, err := afero.ReadFile(fs, file)
	if err != nil {
		return c, err
	}

	if err = yaml.Unmarshal(data, &c); err != nil {
		return c, err
	} else {
		return c, nil
	}
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
