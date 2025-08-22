package main

import (
	"github.com/charmbracelet/log"
	"github.com/unmango/go/cli"
	"github.com/unstoppablemango/openapi2go/cmd"
)

func main() {
	log.SetLevel(log.DebugLevel)
	if err := cmd.Execute(); err != nil {
		cli.Fail(err)
	}
}
