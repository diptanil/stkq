GOFLAGS    ?=
BIN        := stkq
PKG        := github.com/dcxforge/stkq
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE       := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -s -w \
              -X $(PKG)/cmd.Version=$(VERSION) \
              -X $(PKG)/cmd.Commit=$(COMMIT) \
              -X $(PKG)/cmd.Date=$(DATE)

.PHONY: all build install test vet fmt tidy clean

all: build

build:
	go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o ./bin/$(BIN) ./

install:
	go install $(GOFLAGS) -ldflags '$(LDFLAGS)' ./

test:
	go test $(GOFLAGS) -race ./...

vet:
	go vet ./...

fmt:
	gofmt -s -w .

tidy:
	go mod tidy

clean:
	rm -rf ./bin ./dist