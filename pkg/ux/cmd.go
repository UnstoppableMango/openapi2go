package ux

import (
	"fmt"
	"go/format"
	"go/token"
	"os"

	openapi2go "github.com/unstoppablemango/openapi2go/pkg"
	"github.com/unstoppablemango/openapi2go/pkg/config"
	"github.com/unstoppablemango/openapi2go/pkg/openapi"
	"github.com/unstoppablemango/ux/pkg/plugin/skel"
)

func Execute(args *skel.CmdArgs) error {
	if len(args.Args) == 0 {
		return fmt.Errorf("missing specification")
	}

	spec, err := os.ReadFile(args.Args[0])
	if err != nil {
		return fmt.Errorf("reading spec: %w", err)
	}

	doc, err := openapi.ParseDocument(spec)
	if err != nil {
		return fmt.Errorf("parsing openapi doc: %w", err)
	}

	fset := token.NewFileSet()
	files, err := openapi2go.Generate(fset, doc, *config.Default)
	if err != nil {
		return fmt.Errorf("generating code: %w", err)
	}

	for _, f := range files {
		if err := format.Node(os.Stdout, fset, f); err != nil {
			return fmt.Errorf("printing output: %w", err)
		}
	}

	return nil
}

func Generate(args *skel.CmdArgs) error {
	panic("not implemented")
}

func Funcs() skel.UxFuncs {
	return skel.UxFuncs{
		Execute:  Execute,
		Generate: Generate,
	}
}

func Main() {
	skel.PluginMain(Funcs())
}
