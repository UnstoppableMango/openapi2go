package openapi2go

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/unstoppablemango/openapi2go/pkg/openapi"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)

type Config struct {
	PackageName    string
	FileNameSuffix string
	Types          map[string]TypeConfig
}

func (c Config) Generator(doc v3.Document) *Generator {
	return NewGenerator(doc, c)
}

func (c Config) Type(name string, schema *base.Schema) (*ast.GenDecl, error) {
	if conf, ok := c.Types[name]; ok {
		return conf.Type(name, schema)
	} else {
		return openapi.Type(name, schema)
	}
}

type TypeConfig struct {
	Fields map[string]FieldConfig
}

func (c TypeConfig) Type(name string, schema *base.Schema) (*ast.GenDecl, error) {
	var err error

	typ := &ast.StructType{}
	if typ.Fields, err = c.FieldList(schema); err != nil {
		return nil, err
	}

	return &ast.GenDecl{
		Doc: &ast.CommentGroup{},
		Tok: token.TYPE,
		Specs: []ast.Spec{&ast.TypeSpec{
			Name: ast.NewIdent(c.Name(name, schema)),
			Type: typ,
		}},
	}, nil
}

func (TypeConfig) Name(name string, schema *base.Schema) string {
	return name // overkill? yes
}

func (c TypeConfig) FieldList(schema *base.Schema) (*ast.FieldList, error) {
	list := &ast.FieldList{}
	for name, prop := range schema.Properties.FromOldest() {
		if field, err := c.Field(name, prop.Schema()); err != nil {
			return nil, err
		} else {
			list.List = append(list.List, field)
		}
	}

	return list, nil
}

func (c TypeConfig) Field(name string, schema *base.Schema) (*ast.Field, error) {
	if conf, ok := c.Fields[name]; ok {
		return conf.Field(name, schema)
	} else {
		return openapi.Field(name, schema)
	}
}

type FieldConfig struct {
	Type string
}

func (c FieldConfig) Field(name string, schema *base.Schema) (*ast.Field, error) {
	if len(schema.Type) < 1 {
		return nil, fmt.Errorf("at least one type is required")
	}

	typ, err := openapi.Primitive(schema.Type[0])
	if err != nil {
		// return nil, err
		typ = schema.Type[0] // TODO
	}

	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(c.Name(name, schema))},
		Type:  ast.NewIdent(typ),
	}, nil
}

func (c FieldConfig) Name(name string, schema *base.Schema) string {
	return titleCaser.String(name) // TODO: words and stuff
}
