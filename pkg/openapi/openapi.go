package openapi

import (
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func ParseDocument(data []byte) (v3.Document, error) {
	docModel, err := libopenapi.NewDocument(data)
	if err != nil {
		return v3.Document{}, nil
	}

	doc, err := docModel.BuildV3Model()
	if err != nil {
		return v3.Document{}, err
	}
	return doc.Model, nil
}
