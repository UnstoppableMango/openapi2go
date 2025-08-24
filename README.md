# openapi2go

A simple Go code generator for OpenAPI specifications.

## Usage

```text
Generate Go code from OpenAPI specifications

Usage:
  openapi2go [flags]

Flags:
      --config string          Path to a configuration file
  -h, --help                   help for openapi2go
      --output string          Path to the generated code output
      --package-name string    Name of the output package
      --specification string   Path to an OpenAPI specification file
```

### Configuration

Configuration can be provided with `--config` to customize the generated output.
Currently only a few options are supported.

```yaml
packageName: petstore # The Go package name to use. e.g. `package petstore`
types:                # Override output by type, structure is defined by [config.go](./pkg/config/config.go)
  Pet:                # Selects the type matching the name defined in '.components.schemas.*'
    fields:           # Configuration for the fields on `Pet`
      category:       # Selects the field matching the name defined in '.components.schemas.Pet.properties.*'
        type: any     # Override the type specified by '.components.schemas.Pet.properties.category.type'
```

## Build

`make bin/openapi2go`

## Test

`make test`
