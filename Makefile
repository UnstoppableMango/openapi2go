GO_SRC != find . -path '*.go'

build: bin/openapi2go

.PHONY: test
test: bin/petstore.json
	go tool ginkgo run -r

tidy:
	go mod tidy

bin/openapi2go: ${GO_SRC}
	go build -o $@ main.go

bin/petstore.json:
	curl https://petstore3.swagger.io/api/v3/openapi.json | jq -r > $@
