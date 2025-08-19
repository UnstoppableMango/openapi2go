package cmd

import (
	"errors"
	"go/format"
	"go/token"
	"os"
	"path/filepath"

	"github.com/pb33f/libopenapi"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
	openapi2go "github.com/unstoppablemango/openapi2go/pkg"
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
			fsys := afero.NewOsFs()
			spec, err := afero.ReadFile(fsys, opts.Specification)
			if err != nil {
				cli.Fail(err)
			}

			docModel, err := libopenapi.NewDocument(spec)
			if err != nil {
				cli.Fail(err)
			}

			doc, errs := docModel.BuildV3Model()
			if len(errs) > 0 {
				cli.Fail(errors.Join(errs...))
			}

			// TODO: Be less lazy
			if len(opts.PackageName) > 0 {
				conf.PackageName = opts.PackageName
			}

			fset := token.NewFileSet()
			files, err := openapi2go.Generate(fset, doc.Model, conf)
			if err != nil {
				cli.Fail(err)
			}

			if len(opts.Output) > 0 {
				out := filepath.Join(opts.Output, "petstore.go") // TODO
				f, err := fsys.Create(out)
				if err != nil {
					cli.Fail(err)
				}

				if err := format.Node(f, fset, files[0]); err != nil {
					cli.Fail(err)
				}
			} else {
				for _, f := range files {
					if err := format.Node(os.Stdout, fset, f); err != nil {
						cli.Fail(err)
					}
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
	cmd.Flags().StringVar(&opts.PackageName, "package-name", "",
		"Name of the output package",
	)

	return cmd
}
