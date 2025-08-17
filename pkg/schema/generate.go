package schema

import (
	"go/ast"
	"go/token"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func ToType(schema *base.Schema) ast.StructType {
	fset := token.NewFileSet()

	return ast.StructType{
		Fields: Fields(schema),
	}
}

func Fields(schema *base.Schema) *ast.FieldList {
	return &ast.FieldList{}
}
