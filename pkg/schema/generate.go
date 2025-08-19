package schema

import (
	"go/ast"
	"go/token"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func Struct(schema *base.Schema) ast.Decl {
	return &ast.GenDecl{
		Doc: &ast.CommentGroup{},
		Tok: token.TYPE,
		Specs: []ast.Spec{&ast.TypeSpec{
			Name: ast.NewIdent(schema.Title),
			Type: &ast.StructType{
				Fields: Fields(schema),
			},
		}},
	}
}

func Fields(schema *base.Schema) *ast.FieldList {
	list := &ast.FieldList{}
	for name := range schema.Properties.FromOldest() {
		list.List = append(list.List, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(name)},
			Type:  ast.NewIdent("string"),
		})
	}

	return list
}
