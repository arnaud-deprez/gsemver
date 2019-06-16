BUILDDIR   := $(CURDIR)/build
BINDIR     := $(BUILDDIR)/bin
DIST_DIRS  := find * -type d -exec
TARGETS    := darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64le windows/amd64
BINNAME    ?= gsemver

GOPATH        = $(shell go env GOPATH)
GOX           = $(GOPATH)/bin/gox
MOCKGEN		  = $(GOPATH)/bin/mockgen
GOIMPORTS     = $(GOPATH)/bin/goimports
GOLANGCI_LINT = $(GOPATH)/bin/golangci-lint
GHR           = $(GOPATH)/bin/ghr

# go option
PKG        := ./...
TAGS       :=
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    :=
SRC        := $(shell find . -type f -name '*.go' -print)

# Required for globs to work correctly
SHELL      = /bin/bash

# use gsemver to retrieve version
VERSION	   = $(shell go run internal/release/main.go)
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

.PHONY: all
all: build docs release

# ------------------------------------------------------------------------------
#  dependencies
$(GOX):
	go get -u github.com/mitchellh/gox

$(MOCKGEN):
	go get -u github.com/golang/mock/mockgen

$(GOLANGCI_LINT):
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
 
$(GOIMPORTS):
	go get -u golang.org/x/tools/cmd/goimports

$(GHR):
	go get -u github.com/tcnksm/ghr

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

.PHONY: generate
generate: $(MOCKGEN)
	go generate ./...

$(BINDIR)/$(BINNAME): generate $(SRC)
	go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BINNAME) github.com/arnaud-deprez/gsemver/cmd

.PHONY: docs
docs: $(BINDIR)/$(BINNAME)
	@mkdir -p docs/cmd
	@$(BINDIR)/$(BINNAME) docs markdown --dir docs/cmd

# ------------------------------------------------------------------------------
#  test

.PHONY: test
test: build
test: TESTFLAGS += -race -v
test: test-style
test: test-coverage

.PHONY: test-unit
test-unit:
	@echo
	@echo "==> Running unit tests <=="
	go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

.PHONY: test-coverage
test-coverage:
	@echo
	@echo "==> Running unit tests with coverage <=="
	@scripts/coverage.sh --html

.PHONY: test-style
test-style: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

# .PHONY: verify-docs
# verify-docs: build
#	@scripts/verify-docs.sh

# .PHONY: coverage
# coverage:
# 	@scripts/coverage.sh

.PHONY: format
format: $(GOIMPORTS) generate
	go list -f '{{.Dir}}' ./... | xargs $(GOIMPORTS) -w -local github.com/arnaud-deprez/gsemver

# ------------------------------------------------------------------------------
#  release

.PHONY: release
release: build-cross dist checksum

.PHONY: build-cross
build-cross: LDFLAGS += -extldflags "-static"
build-cross: generate $(GOX)
	CGO_ENABLED=0 $(GOX) -parallel=3 -output="$(BUILDDIR)/dist/{{.OS}}-{{.Arch}}/$(BINNAME)" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' github.com/arnaud-deprez/gsemver/cmd

.PHONY: dist
dist:
	( \
		cd $(BUILDDIR)/dist && \
		$(DIST_DIRS) cp ../../LICENSE {} \; && \
		$(DIST_DIRS) cp ../../README.md {} \; && \
		$(DIST_DIRS) tar -zcf $(BINNAME)-$(VERSION)-{}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r $(BINNAME)-$(VERSION)-{}.zip {} \; \
	)

.PHONY: checksum
checksum:
	for f in $(BUILDDIR)/dist/*.{gz,zip} ; do \
		shasum -a 256 "$${f}"  | awk '{print $$1}' > "$${f}.sha256" ; \
	done

.PHONY: publish
publish:
	# TODO: generate changelog and link it to release
	#$(GHR) -body="$$(ghch --latest -F markdown)" v$(VERSION) $(BUILDDIR)/dist
	$(GHR) -delete -prerelease v$$($(BINDIR)/$(BINNAME) bump) $(BUILDDIR)/dist

# ------------------------------------------------------------------------------
# clean

.PHONY: clean
clean:
	rm -rf $(BUILDDIR)
