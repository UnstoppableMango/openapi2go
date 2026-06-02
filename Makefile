_ != mkdir -p bin

PROJECT := openapi2go

CURL		  ?= curl
DOCKER        ?= docker
DPRINT        ?= dprint
GINKGO        ?= ginkgo
GO            ?= go
GOLANGCI_LINT ?= golangci-lint
JQ            ?= jq

GO_SRC != find . -path '*.go'

build: bin/${PROJECT}
tidy: go.sum nix/gomod2nix.toml

.PHONY: test
test: bin/petstore.json
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

coverprofile.out: bin/petstore.json ${GO_SRC}
	$(GINKGO) run -r --cover

bin/${PROJECT}: go.mod ${GO_SRC}
	$(GO) build -o $@ main.go

bin/petstore.json:
	$(CURL) https://petstore3.swagger.io/api/v3/openapi.json | $(JQ) -r > $@

bin/${PROJECT}.tar: Dockerfile .dockerignore ${GO_SRC}
	$(DOCKER) build ${CURDIR} \
	--output type=tar,dest=$@ \
	--output type=image,name=${PROJECT}

go.sum: go.mod
	$(GO) mod tidy

nix/gomod2nix.toml: go.sum
	$(GOMOD2NIX) generate --dir ${CURDIR} --outdir ${@D}
