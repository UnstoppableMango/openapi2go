package openapi

import (
	"errors"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func ParseDocument(data []byte) (v3.Document, error) {
	docModel, err := libopenapi.NewDocument(data)
	if err != nil {
		return v3.Document{}, nil
	}

	doc, errs := docModel.BuildV3Model()
	if len(errs) > 0 {
		return v3.Document{}, errors.Join(errs...)
	}

	return doc.Model, nil
}
