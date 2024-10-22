.ONESHELL:
.SHELLFLAGS = -ce

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

VERSION ?= $(shell git describe --tags --dirty --always)
REVISION ?= $(shell git rev-parse HEAD)

.PHONY: all
all: check
all: test
all: build

.PHONY: check
check:
	@exit_code=0
	scripts/check-format || exit_code=1
	scripts/check-modules || exit_code=1
	exit $$exit_code

.PHONY: test
test:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go test -race -coverprofile cover.out ./... > /dev/null
	GOOS=$(GOOS) GOARCH=$(GOARCH) go tool cover -func=cover.out > cover.report
	>&2 echo [test succeded]

.PHONY: build
build:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v \
		-buildvcs=false \
		-ldflags "-s -w \
			-X github.com/dschott/jsg/version.Version=$(VERSION) \
			-X github.com/dschott/jsg/version.Revision=$(REVISION)" .
	>&2 echo [build succeded]

.PHONY: format
format:
	@gofumpt -w .
