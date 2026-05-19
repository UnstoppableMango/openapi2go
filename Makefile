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
tidy: go.sum

.PHONY: test
test: bin/petstore.json
	$(GINKGO) run -r

docker: bin/${PROJECT}.tar

lint:
	$(GOLANGCI_LINT) run

format fmt:
	$(GO) fmt ./...
	$(DPRINT) fmt

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
