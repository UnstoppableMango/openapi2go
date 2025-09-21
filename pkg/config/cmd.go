package config

import "github.com/spf13/pflag"

type Flags struct {
	Output string
}

func CmdFlags(name string, opts *Flags) *pflag.FlagSet {
	flags := pflag.NewFlagSet(name, pflag.ContinueOnError)
	OutputFlag(flags, opts)
	return flags
}

func OutputFlag(flags *pflag.FlagSet, opts *Flags) {
	flags.StringVarP(&opts.Output, "output", "o", "", "Path to write output data")
}
