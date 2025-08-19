package gen

import (
	"context"
	"errors"
	"go/ast"
	"go/token"

	"github.com/charmbracelet/log"
	"github.com/pb33f/libopenapi"
	"github.com/spf13/afero"
	openapi2go "github.com/unstoppablemango/openapi2go/pkg"
	"github.com/unstoppablemango/openapi2go/pkg/config"
)

func Execute(ctx context.Context, fset *token.FileSet, opts Options) ([]*ast.File, error) {
	log := log.FromContext(ctx)
	if opts.Fs == nil {
		opts.Fs = afero.NewOsFs()
	}

	log.Debug("Reading specification", "path", opts.Specification)
	spec, err := afero.ReadFile(opts.Fs, opts.Specification)
	if err != nil {
		return nil, err
	}

	log.Debug("Parsing spec model")
	doc, err := libopenapi.NewDocument(spec)
	if err != nil {
		return nil, err
	}

	docModel, errs := doc.BuildV3Model()
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	log.Debug("Generating AST")
	return openapi2go.Generate(fset,
		docModel.Model,
		config.Must(config.Read()),
	)
}
