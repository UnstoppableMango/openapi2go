package openapi2go

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/unstoppablemango/openapi2go/pkg/config"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCase = cases.Title(language.English)

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

	f, err := g.parseFile(fset)
	if err != nil {
		return nil, err
	}

	for name, proxy := range g.doc.Components.Schemas.FromOldest() {
		if decl, err := g.Type(name, proxy.Schema(), g.For(name)); err != nil {
			return nil, err
		} else {
			f.Decls = append(f.Decls, decl)
		}
	}

	return []*ast.File{f}, nil
}

func (g *Generator) Array(schema *base.Schema, config *config.Type) (*ast.Ident, error) {
	if items := schema.Items; items.IsA() {
		if s, err := items.A.BuildSchema(); err != nil {
			return nil, err
		} else {
			return ast.NewIdent(s.Type[0]), nil
		}
	}

	return nil, fmt.Errorf("items: bool not supported")
}

func (g *Generator) Field(name string, schema *base.Schema, config *config.Field) (*ast.Field, error) {
	if len(schema.Type) < 1 {
		return nil, fmt.Errorf("no types on field")
	}

	log.Debug("Generating field", "name", name, "config", config)
	if typ, err := g.Primitive(config.TypeFor(schema.Type[0])); err != nil {
		return nil, err
	} else {
		return &ast.Field{
			Names: []*ast.Ident{g.FieldName(name, schema)},
			Type:  ast.NewIdent(typ),
		}, nil
	}
}

func (g *Generator) FieldName(name string, schema *base.Schema) *ast.Ident {
	return ast.NewIdent(titleCase.String(name)) // TODO: words and stuff
}

func (g *Generator) Fields(schema *base.Schema, config *config.Type) (*ast.FieldList, error) {
	list := &ast.FieldList{}
	for name, prop := range schema.Properties.FromOldest() {
		if field, err := g.Field(name, prop.Schema(), config.For(name)); err != nil {
			return nil, err
		} else {
			list.List = append(list.List, field)
		}
	}

	return list, nil
}

func (g *Generator) Primitive(name string) (string, error) {
	switch name {
	case "boolean":
		return "bool", nil
	case "integer":
		return "int", nil
	case "string", "any":
		return name, nil
	default:
		return "", fmt.Errorf("unsupported primitive: %s", name)
	}
}

func (g *Generator) Type(name string, schema *base.Schema, config *config.Type) (*ast.GenDecl, error) {
	var err error

	log.Debug("Generating type", "name", name, "config", config)
	typ := &ast.StructType{}
	if typ.Fields, err = g.Fields(schema, config); err != nil {
		return nil, err
	}

	return &ast.GenDecl{
		Doc: &ast.CommentGroup{},
		Tok: token.TYPE,
		Specs: []ast.Spec{&ast.TypeSpec{
			Name: g.TypeName(name, schema),
			Type: typ,
		}},
	}, nil
}

func (g *Generator) TypeName(name string, schema *base.Schema) *ast.Ident {
	return ast.NewIdent(name) // overkill? yes
}

func (g *Generator) packageName() string {
	if name := g.doc.Info.Title; validPackageName(name) {
		return name
	} else {
		return g.PackageName
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

func (g *Generator) parseType(name string) (string, error) {
	switch name {
	case "boolean":
		return "bool", nil
	case "integer":
		return "int", nil
	case "string", "any":
		return name, nil
	default:
		return "", fmt.Errorf("unsupported primitive: %s", name)
	}
}

func Generate(fset *token.FileSet, doc v3.Document, config config.Config) ([]*ast.File, error) {
	return NewGenerator(doc, config).Execute(fset)
}

func validPackageName(name string) bool {
	return !strings.ContainsAny(name, " \t\n")
}
