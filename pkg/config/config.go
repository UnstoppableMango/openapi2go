package config

import openapi2go "github.com/unstoppablemango/openapi2go/pkg"

var Default = openapi2go.Config{
	PackageName: "openapi2go",
}

func Must(config openapi2go.Config, err error) openapi2go.Config {
	if err != nil {
		panic(err)
	}

	return config
}
