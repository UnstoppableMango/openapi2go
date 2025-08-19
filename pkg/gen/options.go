package gen

import "github.com/spf13/afero"

type Options struct {
	Output        string
	Specification string
	Fs            afero.Fs
}
