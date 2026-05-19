# AGENTS.md

This file provides guidance to ai agents when working with code in this repository.

## Commands

```bash
make test          # Run tests (uses ginkgo)
make lint          # Run golangci-lint
make fmt           # Format with go fmt + dprint
make tidy          # go mod tidy
make bin/openapi2go  # Build binary
```

Run a single test file or suite via ginkgo directly:
```bash
ginkgo ./test/e2e/...
```

CI runs inside `nix develop`: `nix develop -c make test`

## Architecture

`openapi2go` converts OpenAPI specs into Go struct type declarations.

**Data flow:**
```
OpenAPI JSON/YAML → libopenapi (v3.Document) → Generator → Go AST ([]*ast.File) → go/format → Go source
```

**Key packages:**

- `pkg/generator.go` — Core engine. `Generate(fset, doc, config)` is the public API. Walks `Components.Schemas`, builds Go AST nodes (`ast.GenDecl`, `ast.StructType`, `ast.Field`). Type resolution order: config override → schema type → fallback `any`.

- `pkg/config/` — YAML config for per-type and per-field Go type overrides. `Config.For(typeName)` returns overrides. Default package: `openapi2go`, output suffix: `.zz_generated.go`.

- `pkg/gen/options.go` — CLI flag parsing. `Options.ReadSpec()` parses OpenAPI via libopenapi; `OutputWriter()` returns file or stdout writer.

- `pkg/openapi/openapi.go` — Thin wrapper around libopenapi's `ParseDocument`.

- `cmd/root.go` — Cobra root command: reads spec → calls `generator.Generate()` → writes formatted Go source.

- `pkg/ux/cmd.go` — Interactive plugin mode via `github.com/unstoppablemango/ux`.

**Tests** live in `test/e2e/` (Ginkgo), compare generated output against `test/e2e/testdata/petstore.go` using the Petstore v3 OpenAPI example.

## Config format

```yaml
types:
  Pet:
    type: MyPet
    fields:
      id:
        type: int64
```

## Nix

Dev shell defined in `flake.nix`. All tools (go, ginkgo, golangci-lint, dprint, make) available via `nix develop`.
