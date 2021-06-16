PROJECT := gotway
VERSION=$(git describe --abbrev=0 --tags)
LD_FLAGS := -X main.version=$(VERSION) -s -w
SOURCE_FILES ?= ./internal/... ./pkg/... ./cmd/...
UNAME		:= $(uname -s)

.PHONY: all
all: help

.PHONY: help
help:	### Show targets documentation
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

.PHONY: clean
clean: ### Clean build files
	@rm -rf ./build
	@go clean

.PHONY: build
build: clean ### Build binary
	@go build -tags netgo -a -v -ldflags "${LD_FLAGS}" -o ./build/gotway ./cmd/gotway/*.go
	@chmod +x ./bin/*

.PHONY: run
run: ### Quick run
	@go run -race cmd/gotway/*.go

.PHONY: deps
deps: ### Optimize dependencies
	@go mod tidy

.PHONY: install
install: ### Install binary in your system
	@go install -v cmd/gotway/*.go

.PHONY: fmt
fmt: ### Format
	@gofmt -s -w .

.PHONY: vet
vet: ### Vet
	@go vet ./...

### Lint
.PHONY: lint
lint: fmt vet

### Clean test 
.PHONY: test-clean
test-clean: ### Clean test cache
	@go clean -testcache ./...

.PHONY: test
test: lint ### Run tests
	@go test -v  -coverprofile=cover.out -timeout 10s ./...

.PHONY: cover
cover: test ### Run tests and generate coverage
	@go tool cover -html=cover.out -o=cover.html

.PHONY: mocks
mocks: ### Generate mocks
	@mockery --all --output internal/mocks