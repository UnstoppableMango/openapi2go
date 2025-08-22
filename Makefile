_ != mkdir -p bin

GO_SRC != find . -path '*.go'

build: bin/openapi2go

.PHONY: test
test: bin/petstore.json
	go tool ginkgo run -r

docker: bin/openapi2go.tar

tidy:
	go mod tidy

bin/openapi2go: ${GO_SRC}
	go build -o $@ main.go

bin/petstore.json:
	curl https://petstore3.swagger.io/api/v3/openapi.json | jq -r > $@

bin/openapi2go.tar: Dockerfile .dockerignore ${GO_SRC}
	docker build ${CURDIR} \
	--output type=tar,dest=$@ \
	--output type=image,name=openapi2go
