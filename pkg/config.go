package openapi2go

import v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

type Config struct {
	PackageName    string
	FileNameSuffix string
}

func (c Config) Generator(doc v3.Document) *Generator {
	return NewGenerator(doc, c)
}
