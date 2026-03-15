SHELL := /bin/bash

VERSION ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo dev)
BUILD_DIR ?= build
BIN ?= git-tk
LDFLAGS := -X github.com/bharath23/git-treekeeper/cmd.version=$(VERSION)
GOTESTFLAGS ?=
VERBOSE ?=
V ?=
TESTJSON ?= test.json

.PHONY: build install check test test-acceptance test-v test-verbose test-ci

ifneq ($(strip $(VERBOSE)$(V)),)
GOTESTFLAGS += -v
endif

build:
	mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BIN) .

install:
	go install -ldflags "$(LDFLAGS)" ./...

check:
	@unformatted=$$(gofmt -l $$(git ls-files '*.go')); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt needed on:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	go vet ./...

test:
	go test ./... $(GOTESTFLAGS)
	$(MAKE) test-acceptance

test-acceptance: build
	GIT_TK_BIN="$(BUILD_DIR)/$(BIN)" scripts/acceptance-test.sh

test-v:
	$(MAKE) test VERBOSE=1

test-verbose:
	$(MAKE) test VERBOSE=1

test-ci:
	set -o pipefail; go test ./... -json | tee $(TESTJSON)
	$(MAKE) test-acceptance
