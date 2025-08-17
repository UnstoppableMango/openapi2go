package config

import (
	"github.com/spf13/viper"
	"github.com/unmango/go/cli"
	openapi2go "github.com/unstoppablemango/openapi2go/pkg"
)

func Read() (*openapi2go.Config, error) {
	viper := viper.New()

	config := &openapi2go.Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	} else {
		return config, nil
	}
}

func Must(config *openapi2go.Config, err error) *openapi2go.Config {
	if err != nil {
		cli.Fail(err)
	}

	return config
}
