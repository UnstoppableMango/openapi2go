package cmd

import (
	"fmt"
	"go/types"
	"os"

	"github.com/pb33f/libopenapi"
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
)

type GenerateOptions struct {
	Output        string
	Specification string
}

var generate = NewGenerate()

func init() {
	root.AddCommand(generate)
}

func NewGenerate() *cobra.Command {
	opts := &GenerateOptions{}

	cmd := &cobra.Command{
		Use: "generate",
		Run: func(cmd *cobra.Command, args []string) {
			// config := config.Must(config.Read())
			spec, err := os.ReadFile(opts.Specification)
			if err != nil {
				cli.Fail(err)
			}

			doc, err := libopenapi.NewDocument(spec)
			if err != nil {
				cli.Fail(err)
			}

			model, errors := doc.BuildV3Model()
			if len(errors) > 0 {
				cli.Fail(errors)
			}

			// fset := token.NewFileSet()
			// f := fset.AddFile("order.go", fset.Base(), 0)
			// format.Node(os.Stdout, fset, f)
			pkg := types.NewPackage("order.go", "order")
			s := types.NewStruct([]*types.Var{
				types.NewField(0, pkg, "Test", types.Typ[types.String], false),
			}, nil)
			n := types.NewTypeName(0, pkg, "Order", s)
			fmt.Println(n.String())

			for name, value := range model.Model.Components.Schemas.FromOldest() {
				schema := value.Schema()
				fmt.Printf("Found: %s = %v\n", name, schema)
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
