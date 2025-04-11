BUILDDIR   := $(CURDIR)/build
BINDIR     := $(BUILDDIR)/bin
DIST_DIRS  := find * -type d -exec
BINNAME    ?= gsemver

GO_NOMOD      := GO111MODULE=off go
GOPATH        := $(shell go env GOPATH)
MOCKGEN		  := $(GOPATH)/bin/mockgen
GOIMPORTS     := $(GOPATH)/bin/goimports
GOLANGCI_LINT := $(GOPATH)/bin/golangci-lint
GIT_CHGLOG    := $(GOPATH)/bin/git-chglog

# go option
PKG        := ./...
TAGS       := 
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    :=
SRC        := $(shell find . -type f -name '*.go' -print)

# Required for globs to work correctly
SHELL      := /bin/bash

# use gsemver to retrieve version
GIT_BRANCH ?= $(shell git symbolic-ref --short HEAD)
GIT_COMMIT ?= $(shell git rev-parse HEAD)
VERSION	   =  $(shell go run internal/release/main.go)
GIT_DIRTY  =  $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
BUILD_DATE =  $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LAST_TAG   =  $(shell git describe --tags --abbrev=0 --first-parent --match v[0-9]*.[0-9]*.[0-9]* $(GIT_COMMIT)~ || echo "") 

LDFLAGS += -X github.com/arnaud-deprez/gsemver/internal/version.version=v$(VERSION)
LDFLAGS += -X github.com/arnaud-deprez/gsemver/internal/version.gitCommit=$(GIT_COMMIT)
LDFLAGS += -X github.com/arnaud-deprez/gsemver/internal/version.gitTreeState=$(GIT_DIRTY)
LDFLAGS += -X github.com/arnaud-deprez/gsemver/internal/version.buildDate=$(BUILD_DATE)


.PHONY: all
all: build docs release

# ------------------------------------------------------------------------------
#  dependencies
$(MOCKGEN):
	go install go.uber.org/mock/mockgen@v0.5.1

$(GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(GOPATH)/bin v2.0.2
 
$(GOIMPORTS):
	go install golang.org/x/tools/cmd/goimports@latest

$(GIT_CHGLOG):
	go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

# ------------------------------------------------------------------------------
#  build

build: download-dependencies generate $(BINDIR)/$(BINNAME) docs

download-dependencies: 
	go mod download

.PHONY: generate
generate: download-dependencies $(MOCKGEN)
	go generate ./...

.PHONY: build
$(BINDIR)/$(BINNAME): generate $(SRC)
	go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BINNAME) github.com/arnaud-deprez/gsemver

.PHONY: download-dependencies docs
docs: $(BINDIR)/$(BINNAME)
	mkdir -p docs/cmd
	$(BINDIR)/$(BINNAME) docs markdown --dir docs/cmd

# ------------------------------------------------------------------------------
#  test

.PHONY: test
test: build
test: TESTFLAGS += -race -v
test: test-style
test: test-coverage

.PHONY: test-coverage
test-coverage:
	@echo
	@echo "==> Running unit tests with coverage <=="
	scripts/coverage.sh

.PHONY: test-style
test-style: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

.PHONY: test-unit
test-unit: test
	@echo
	@echo "==> Running unit tests <=="
	go test $(GOFLAGS) -run $(TESTS) $(PKG) -short $(TESTFLAGS)

.PHONY: test-integration
test-integration: test
	@echo
	@echo "==> Running integration tests <=="
	go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS) -failfast

# .PHONY: verify-docs
# verify-docs: build
#	@scripts/verify-docs.sh

.PHONY: format
format: $(GOIMPORTS) generate
	go list -f '{{.Dir}}' ./... | xargs $(GOIMPORTS) -w -local github.com/arnaud-deprez/gsemver

# ------------------------------------------------------------------------------
#  release

.PHONY: test-release
test-release: $(GIT_CHGLOG)
	@echo "Test release $(VERSION) on $(GIT_BRANCH), last version was $(LAST_TAG)"
	# Because of https://github.com/git-chglog/git-chglog/issues/45, it will generate changelog for both LAST_TAG and VERSION
	export GIT_DIRTY=$(GIT_DIRTY) && curl -sL https://git.io/goreleaser | bash -s -- release --config=./.goreleaser.yml --snapshot --skip=publish --verbose --clean --release-notes <($(GIT_CHGLOG) --next-tag v$(VERSION) $(strip $(LAST_TAG))..)

.PHONY: release
release: $(GIT_CHGLOG)
	@echo "Release $(VERSION) on $(GIT_BRANCH), last version was $(LAST_TAG)"
	git tag -am "Release v$(VERSION) by ci script" v$(VERSION)
	# This is a bit weird: https://github.com/git-chglog/git-chglog/issues/45
	export GIT_DIRTY=$(GIT_DIRTY) && curl -sL https://git.io/goreleaser | bash -s -- release --config=./.goreleaser.yml --clean --release-notes <($(GIT_CHGLOG) v$(VERSION))

# ------------------------------------------------------------------------------
# clean

.PHONY: clean
clean:
	@rm -rf $(BUILDDIR)
