package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "openapi2go",
	Short: "Generate Go code from OpenAPI specifications",
}

func Execute() error {
	log.SetLevel(log.DebugLevel)
	return root.Execute()
}
