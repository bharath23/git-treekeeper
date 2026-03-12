VERSION ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo dev)
BUILD_DIR ?= build
BIN ?= git-tk
LDFLAGS := -X github.com/bharath23/git-treekeeper/cmd.version=$(VERSION)

.PHONY: build install test

build:
	mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BIN) .

install:
	go install -ldflags "$(LDFLAGS)" ./...

test:
	go test ./tests/... -v
