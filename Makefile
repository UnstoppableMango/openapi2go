_ != mkdir -p bin

PROJECT := openapi2go

CURL		  ?= curl
DOCKER        ?= docker
DPRINT        ?= dprint
GINKGO        ?= ginkgo
GO            ?= go
GOLANGCI_LINT ?= golangci-lint

GO_SRC != find . -path '*.go'

build: bin/${PROJECT}
tidy: go.sum nix/gomod2nix.toml

.PHONY: test
test: bin/petstore.yaml
	$(GINKGO) run -r

cover: coverprofile.out
docker: bin/${PROJECT}.tar

lint:
	$(GOLANGCI_LINT) run

format fmt:
	$(GO) fmt ./...
	$(DPRINT) fmt

update:
	nix flake update

check:
	nix flake check

coverprofile.out: bin/petstore.yaml ${GO_SRC}
	$(GINKGO) run -r --cover

bin/${PROJECT}: go.mod ${GO_SRC}
	$(GO) build -o $@ main.go

bin/petstore.yaml:
	$(CURL) -fsSL https://raw.githubusercontent.com/swagger-api/swagger-petstore/5cb909fc7780bc708151bce0014939acba8b212a/src/main/resources/openapi.yaml > $@

bin/${PROJECT}.tar: Dockerfile .dockerignore ${GO_SRC}
	$(DOCKER) build ${CURDIR} \
	--output type=tar,dest=$@ \
	--output type=image,name=${PROJECT}

go.sum: go.mod
	$(GO) mod tidy

nix/gomod2nix.toml: go.sum
	$(GOMOD2NIX) generate --dir ${CURDIR} --outdir ${@D}
