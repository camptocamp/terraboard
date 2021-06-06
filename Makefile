NAME          := terraboard
FILES         := $(wildcard */*.go)
VERSION       := $(shell git describe --always)
.DEFAULT_GOAL := help

export GO111MODULE=on

.PHONY: setup
setup: ## Install required libraries/tools for build tasks
	@command -v cover 2>&1 >/dev/null       || GO111MODULE=off go get -u -v golang.org/x/tools/cmd/cover
	@command -v goveralls 2>&1 >/dev/null   || GO111MODULE=off go get -u -v github.com/mattn/goveralls
	@command -v goimports 2>&1 >/dev/null   || GO111MODULE=off go get -u -v golang.org/x/tools/cmd/goimports
	@command -v golangci-lint 2>&1 >/dev/null || GO111MODULE=off go get -v github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: fmt
fmt: setup ## Format source code
	goimports -w $(FILES)

.PHONY: lint
lint: setup ## Run all lint related tests against the codebase (will use the .golangci.yml config)
	golangci-lint run

.PHONY: test
test: ## Run the tests against the codebase
	go test -v -race ./...

.PHONY: build
build: main.go $(FILES) ## Build the binary
	CGO_ENABLED=1 GOOS=linux \
		go build \
		-trimpath \
		-ldflags "-linkmode external -extldflags -static -X main.version=$(VERSION)" \
		-o $(NAME) $<
	strip $(NAME)

.PHONY: install
install: ## Install the binary using local environment
	go install .

.PHONY: vendor
vendor: # Vendor go modules
	go mod vendor

.PHONY: coverage
coverage: ## Generates coverage report
	rm -f coverage.out
	go test -v ./... -coverpkg=./... -coverprofile=coverage.out

.PHONY: publish-coveralls
publish-coveralls: setup ## Publish coverage results on coveralls
	goveralls -service=travis-ci -coverprofile=coverage.out

.PHONY: clean
clean: ## Remove binary if it exists
	rm -f $(NAME)

.PHONY: all
all: lint test build coverage

.PHONY: help
help: ## Displays this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
