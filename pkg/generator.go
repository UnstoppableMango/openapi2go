package openapi2go

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/unstoppablemango/openapi2go/pkg/config"
	ux "github.com/unstoppablemango/ux/pkg"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCase = cases.Title(language.English)

type Generator struct {
	config.Config
	doc  v3.Document
	spec uuid.UUID
}

func NewGenerator(doc v3.Document, config config.Config) *Generator {
	return &Generator{config, doc, uuid.Nil}
}

func (g *Generator) Configure(b ux.Inputs) error {
	g.spec = b.Add(nil)

	return nil
}

func (g *Generator) Generate(ctx ux.Context) error {
	fset := token.NewFileSet()
	r, err := ctx.Input(g.spec)
	if err != nil {
		return err
	}

	files, err := g.Execute(fset)
	if err != nil {
		return err
	}

	return nil
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

func (g *Generator) Array(schema *base.Schema, config *config.Field) (*ast.ArrayType, error) {
	if schema.Items.IsB() {
		// https://spec.openapis.org/oas/v3.1.0#schema-object
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.1.2
		return nil, fmt.Errorf("items: bool not supported")
	}

	s, err := schema.Items.A.BuildSchema()
	if err != nil {
		return nil, err
	}

	typ, err := g.FieldType(s, config)
	if err != nil {
		return nil, err
	}

	return &ast.ArrayType{Elt: typ}, nil
}

func (g *Generator) Field(name string, schema *base.Schema, config *config.Field) (*ast.Field, error) {
	if len(schema.Type) < 1 {
		return nil, fmt.Errorf("no types on field")
	}

	log.Debug("Generating field", "name", name, "config", config)
	typ, err := g.FieldType(schema, config)
	if err != nil {
		return nil, err
	}

	return &ast.Field{
		Names: []*ast.Ident{g.FieldName(name, schema)},
		Type:  typ,
	}, nil
}

func (g *Generator) FieldName(name string, schema *base.Schema) *ast.Ident {
	return ast.NewIdent(titleCase.String(name)) // TODO: words and stuff
}

func (g *Generator) FieldType(schema *base.Schema, config *config.Field) (ast.Expr, error) {
	if len(schema.Type) < 1 {
		return nil, fmt.Errorf("no types on field")
	}

	switch typ := schema.Type[0]; typ {
	case "boolean":
		return ast.NewIdent("bool"), nil
	case "integer":
		return ast.NewIdent("int"), nil
	case "string", "any":
		return ast.NewIdent(typ), nil
	case "array":
		return g.Array(schema, config)
	case "object":
		return ast.NewIdent("any"), nil // TODO
	default:
		return ast.NewIdent(typ), nil // TODO
	}
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

func (g *Generator) Bool() *ast.Ident {
	return ast.NewIdent("bool")
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

func Generate(fset *token.FileSet, doc v3.Document, config config.Config) ([]*ast.File, error) {
	return NewGenerator(doc, config).Execute(fset)
}

func validPackageName(name string) bool {
	return !strings.ContainsAny(name, " \t\n") // TODO
}
