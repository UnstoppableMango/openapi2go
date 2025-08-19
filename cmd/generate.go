package cmd

import (
	"go/format"
	"go/token"
	"os"

	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
	"github.com/unstoppablemango/openapi2go/pkg/gen"
)

var generate = NewGenerate()

func init() {
	root.AddCommand(generate)
}

func NewGenerate() *cobra.Command {
	opts := gen.Options{}

	cmd := &cobra.Command{
		Use: "generate",
		Run: func(cmd *cobra.Command, args []string) {
			fset := token.NewFileSet()
			files, err := gen.Execute(cmd.Context(), fset, opts)
			if err != nil {
				cli.Fail(err)
			}

			for _, file := range files {
				if err := format.Node(os.Stdout, fset, file); err != nil {
					cli.Fail(err)
				}
			}
		},
	}

	cmd.Flags().StringVar(&opts.Specification, "specification", "",
		"Path to an OpenAPI specification file",
	)
	cmd.Flags().StringVar(&opts.Output, "output", "",
		"Path to the generated code output",
	)

	return cmd
}
