package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	openapi2go "github.com/unstoppablemango/openapi2go/pkg"
	"github.com/unstoppablemango/openapi2go/pkg/config"
)

var (
	conf       openapi2go.Config
	configFile string

	root = &cobra.Command{
		Use:   "openapi2go",
		Short: "Generate Go code from OpenAPI specifications",
	}
)

func Execute() error {
	return root.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	root.PersistentFlags().StringVar(&configFile, "config", "",
		"Path to a configuration file",
	)
}

func initConfig() {
	v := config.Viper(configFile)
	if err := v.ReadInConfig(); err == nil {
		log.Debug("Using config file", "path", v.ConfigFileUsed())
	}

	if c, err := config.Parse(v); err != nil {
		log.Debug("Using default config")
		conf = config.Default
	} else {
		conf = c
	}
}
