package gen

import "github.com/spf13/afero"

type Options struct {
	Output        string
	PackageName   string
	Specification string
	Fs            afero.Fs
}

func (o Options) readSpec() ([]byte, error) {
	return afero.ReadFile(o.Fs, o.Specification)
}
