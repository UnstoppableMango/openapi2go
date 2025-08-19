package openapi

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)

func Field(name string, schema *base.Schema) (*ast.Field, error) {
	if len(schema.Type) < 1 {
		return nil, fmt.Errorf("at least one type is required")
	}

	typ, err := Primitive(schema.Type[0])
	if err != nil {
		// return nil, err
		typ = schema.Type[0]
	}

	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(FieldName(name, schema))},
		Type:  ast.NewIdent(typ),
	}, nil
}

func FieldName(name string, schema *base.Schema) string {
	return titleCaser.String(name) // TODO: words and stuff
}

func Fields(schema *base.Schema) (*ast.FieldList, error) {
	list := &ast.FieldList{}
	for name, prop := range schema.Properties.FromOldest() {
		if field, err := Field(name, prop.Schema()); err != nil {
			return nil, err
		} else {
			list.List = append(list.List, field)
		}
	}

	return list, nil
}

func PackageName(doc v3.Document) string {
	return doc.Info.Title // TODO: Super smart algorithm
}

func Primitive(name string) (string, error) {
	switch name {
	case "integer":
		return "int", nil
	case "string":
		return name, nil
	default:
		return "", fmt.Errorf("unsupported primitive: %s", name)
	}
}

func Type(name string, schema *base.Schema) (*ast.GenDecl, error) {
	var err error

	typ := &ast.StructType{}
	if typ.Fields, err = Fields(schema); err != nil {
		return nil, err
	}

	return &ast.GenDecl{
		Doc: &ast.CommentGroup{},
		Tok: token.TYPE,
		Specs: []ast.Spec{&ast.TypeSpec{
			Name: ast.NewIdent(TypeName(name, schema)),
			Type: typ,
		}},
	}, nil
}

func TypeName(name string, schema *base.Schema) string {
	return name // overkill? yes
}
