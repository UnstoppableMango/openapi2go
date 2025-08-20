package config

import openapi2go "github.com/unstoppablemango/openapi2go/pkg"

const DefaultFileSuffix = ".zz_generated.go"

var Default = openapi2go.Config{
	PackageName:    "openapi2go",
	FileNameSuffix: DefaultFileSuffix,
}

func Must(config openapi2go.Config, err error) openapi2go.Config {
	if err != nil {
		panic(err)
	}

	return config
}
