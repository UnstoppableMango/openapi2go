# openapi2go

[![CI](https://github.com/UnstoppableMango/openapi2go/actions/workflows/ci.yml/badge.svg)](https://github.com/UnstoppableMango/openapi2go/actions/workflows/ci.yml)
[![Last commit](https://img.shields.io/github/last-commit/UnstoppableMango/openapi2go)](https://github.com/UnstoppableMango/openapi2go/commits/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/unstoppablemango/openapi2go)](https://goreportcard.com/report/github.com/unstoppablemango/openapi2go)
[![Go Reference](https://pkg.go.dev/badge/github.com/unstoppablemango/openapi2go.svg)](https://pkg.go.dev/github.com/unstoppablemango/openapi2go)
[![codecov](https://codecov.io/gh/UnstoppableMango/openapi2go/graph/badge.svg)](https://codecov.io/gh/UnstoppableMango/openapi2go)
[![GitHub release](https://img.shields.io/github/v/release/UnstoppableMango/openapi2go)](https://github.com/UnstoppableMango/openapi2go/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

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

> [!WARN]
> Actually none of this works yet.

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
