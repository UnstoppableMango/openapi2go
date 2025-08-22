package config

const DefaultFileSuffix = ".zz_generated.go"

var Default = Config{
	PackageName:    "openapi2go",
	FileNameSuffix: DefaultFileSuffix,
}

type Config struct {
	PackageName    string
	FileNameSuffix string
	Types          map[string]Type
}

func (c Config) ForType(name string) *Type {
	if t, ok := c.Types[name]; ok {
		return &t
	} else {
		return nil
	}
}

type Type struct {
	Fields map[string]Field
}

func (t *Type) ForField(name string) *Field {
	if t == nil {
		return nil
	}

	if f, ok := t.Fields[name]; ok {
		return &f
	} else {
		return nil
	}
}

type Field struct {
	Type string
}

func (c *Field) TypeFor(given string) string {
	if c != nil && c.Type != "" {
		return c.Type
	} else {
		return given
	}
}

func Must(config Config, err error) Config {
	if err != nil {
		panic(err)
	}

	return config
}
