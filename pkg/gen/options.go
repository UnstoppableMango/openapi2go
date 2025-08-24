package gen

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/spf13/afero"
	"github.com/unstoppablemango/openapi2go/pkg/config"
)

type Options struct {
	Config        string
	Output        string
	PackageName   string
	Specification string
	Fs            afero.Fs
}

func (o Options) Apply(config *config.Config) {
	if len(o.PackageName) > 0 {
		config.PackageName = o.PackageName
	}
}

func (o Options) OutputWriter(fsys afero.Fs) (io.Writer, error) {
	if len(o.Output) == 0 {
		return os.Stdout, nil
	}

	out := filepath.Join(o.Output, "petstore.go")
	log.Info("Writing", "out", out)
	return fsys.Create(out)
}

func (o Options) ReadConfig(fs afero.Fs) (*config.Config, error) {
	config, err := config.ReadFile(fs, o.Config)
	if err != nil {
		return nil, err
	}
	if len(o.PackageName) > 0 {
		config.PackageName = o.PackageName
	}

	return config, nil
}

func (o Options) ReadSpec(fsys afero.Fs) (v3.Document, error) {
	spec, err := afero.ReadFile(fsys, o.Specification)
	if err != nil {
		return v3.Document{}, err
	}

	docModel, err := libopenapi.NewDocument(spec)
	if err != nil {
		return v3.Document{}, nil
	}

	doc, errs := docModel.BuildV3Model()
	if len(errs) > 0 {
		return v3.Document{}, errors.Join(errs...)
	}

	return doc.Model, nil
}
