package openapi2go

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/unstoppablemango/openapi2go/pkg/config"
	"github.com/unstoppablemango/openapi2go/pkg/openapi"
)

type Generator struct {
	config.Config
	doc v3.Document
}

func NewGenerator(doc v3.Document, config config.Config) *Generator {
	return &Generator{config, doc}
}

func (g *Generator) Execute(fset *token.FileSet) ([]*ast.File, error) {
	if c := g.doc.Components; c == nil || c.Schemas == nil {
		return nil, nil
	}

	decls := map[string]ast.Decl{}
	for name, proxy := range g.doc.Components.Schemas.FromOldest() {
		log.Info("Generating types", "name", name)
		if decl, err := openapi.Type(name, proxy.Schema(), &g.Config); err != nil {
			return nil, err
		} else {
			decls[name] = decl
		}
	}

	if f, err := g.parseFile(fset); err != nil {
		return nil, err
	} else {
		f.Decls = slices.Collect(maps.Values(decls))
		return []*ast.File{f}, nil
	}
}

func (g *Generator) packageName() string {
	if n := openapi.PackageName(g.doc); validPackageName(n) {
		return g.PackageName
	} else {
		return n
	}
}

func (g *Generator) filename() string {
	return g.packageName() + g.FileNameSuffix
}

func (g *Generator) parseFile(fset *token.FileSet) (*ast.File, error) {
	return parser.ParseFile(fset,
		g.filename(),
		fmt.Sprintf("package %s", g.packageName()),
		parser.SkipObjectResolution,
	)
}

func Generate(fset *token.FileSet, doc v3.Document, config config.Config) ([]*ast.File, error) {
	return NewGenerator(doc, config).Execute(fset)
}

func validPackageName(name string) bool {
	return strings.ContainsAny(name, " \t\n")
}
