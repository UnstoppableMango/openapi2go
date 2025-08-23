package cmd

import (
	"go/format"
	"go/token"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
	openapi2go "github.com/unstoppablemango/openapi2go/pkg"
	"github.com/unstoppablemango/openapi2go/pkg/config"
	"github.com/unstoppablemango/openapi2go/pkg/gen"
)

var (
	configFile string
	opts       gen.Options

	root = &cobra.Command{
		Use:   "openapi2go",
		Short: "Generate Go code from OpenAPI specifications",
		Run: func(cmd *cobra.Command, args []string) {
			fsys := afero.NewOsFs()
			model, err := opts.ReadSpec(fsys)
			if err != nil {
				cli.Fail(err)
			}

			conf, err := config.Read(fsys, configFile)
			if err != nil {
				cli.Fail(err)
			}

			opts.Apply(&conf)

			fset := token.NewFileSet()
			files, err := openapi2go.Generate(fset, model, conf)
			if err != nil {
				cli.Fail(err)
			}

			for _, f := range files {
				w, err := opts.OutputWriter(fsys)
				if err != nil {
					cli.Fail(err)
				}

				if err := format.Node(w, fset, f); err != nil {
					cli.Fail(err)
				}
			}
		},
	}
)

func Execute() error {
	return root.Execute()
}

func init() {
	root.PersistentFlags().StringVar(&configFile, "config", "",
		"Path to a configuration file",
	)

	root.Flags().StringVar(&opts.Specification, "specification", "",
		"Path to an OpenAPI specification file",
	)
	root.Flags().StringVar(&opts.Output, "output", "",
		"Path to the generated code output",
	)
	root.Flags().StringVar(&opts.PackageName, "package-name", "",
		"Name of the output package",
	)
}
