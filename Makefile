_ != mkdir -p bin

PROJECT := openapi2go
GO_SRC != find . -path '*.go'

build: bin/${PROJECT}
tidy: go.sum

.PHONY: test
test: bin/petstore.json
	go tool ginkgo run -r

docker: bin/${PROJECT}.tar

lint:
	go tool golangci-lint run

format fmt:
	go fmt
	dprint fmt

bin/${PROJECT}: go.mod ${GO_SRC}
	go build -o $@ main.go

bin/petstore.json:
	curl https://petstore3.swagger.io/api/v3/openapi.json | jq -r > $@

bin/${PROJECT}.tar: Dockerfile .dockerignore ${GO_SRC}
	docker build ${CURDIR} \
	--output type=tar,dest=$@ \
	--output type=image,name=${PROJECT}

go.sum: go.mod
	go mod tidy
