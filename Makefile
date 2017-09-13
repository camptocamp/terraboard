DEPS = $(wildcard */*.go)
VERSION = $(shell git describe --always)

all: test terraboard

terraboard: main.go $(DEPS)
	CGO_ENABLED=1 GOOS=linux go build \
	  -ldflags "-linkmode external -extldflags -static -X main.version=$(VERSION)" \
	-o $@ $<
	strip $@

lint:
	@ go get -v github.com/golang/lint/golint
	@for file in $$(git ls-files '*.go' | grep -v '_workspace/'); do \
		export output="$$(golint $${file} | grep -v 'type name will be used as docker.DockerInfo')"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
	done; \
	exit $${status:-0}

vet: main.go
	go vet $<

imports: main.go
	goimports -d $<

test: lint vet imports
	go test -v ./...

coverage:
	rm -rf *.out
	go test -coverprofile=coverage.out
	for i in config util s3 db api compare auth; do \
	 	go test -coverprofile=$$i.coverage.out github.com/camptocamp/terraboard/$$i; \
		tail -n +2 $$i.coverage.out >> coverage.out; \
		done

clean:
	rm -f terraboard

.PHONY: all lint vet imports test coverage clean
