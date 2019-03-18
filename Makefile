DOCKER_REGISTRY   ?= gcr.io
DEV_IMAGE         ?= golang:1.12
TARGETS           ?= darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64le linux/s390x windows/amd64
TARGET_OBJS       ?= darwin-amd64.tar.gz darwin-amd64.tar.gz.sha256 linux-amd64.tar.gz linux-amd64.tar.gz.sha256 linux-386.tar.gz linux-386.tar.gz.sha256 linux-arm.tar.gz linux-arm.tar.gz.sha256 linux-arm64.tar.gz linux-arm64.tar.gz.sha256 linux-ppc64le.tar.gz linux-ppc64le.tar.gz.sha256 linux-s390x.tar.gz linux-s390x.tar.gz.sha256 windows-amd64.zip windows-amd64.zip.sha256
DIST_DIRS         = find * -type d -exec

# go option
GO        ?= go
TAGS      :=
TESTS     := .
TESTFLAGS :=
LDFLAGS   := -w -s
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin
BINARIES  := helm tiller