package cmd

import (
	"github.com/spf13/cobra"
)

var (
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
	root.PersistentFlags().StringVar(&configFile, "config", "",
		"Path to a configuration file",
	)
}
