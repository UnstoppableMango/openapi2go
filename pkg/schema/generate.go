package schema

import (
	"go/types"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func Struct(schema *base.Schema) *types.Struct {
	return types.NewStruct(Fields(schema), nil)
}

func Fields(schema *base.Schema) []*types.Var {
	return []*types.Var{}
}
