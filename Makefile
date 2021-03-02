VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
SRC_DIR := ./
BIN_NAME := alfred-bookmarks
BIN_DIR := bin
BINARY := $(BIN_DIR)/$(BIN_NAME)
ASSETS_DIR := assets
ASSETS := $(ASSETS_DIR)/* $(BINARY) README.md
ARTIFACT_DIR := .artifact
ARTIFACT_NAME := $(BIN_NAME).alfredworkflow

## For local test
WORKFLOW_DIR := "$${HOME}/Library/Application Support/Alfred/Alfred.alfredpreferences/workflows/user.workflow.7C42A657-124F-46B8-89EE-7A1C06594E13"

GOLANGCI_LINT_VERSION := v1.30.0
export GO111MODULE=on

## Build binaries on your environment
build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(SRC_DIR)

## Format source codes
fmt:
	@(if ! type goimports >/dev/null 2>&1; then go get -u golang.org/x/tools/cmd/goimports ;fi)
	goimports -w $$(go list -f {{.Dir}} ./... | grep -v /vendor/)

## Lint
lint:
	@(if ! type golangci-lint >/dev/null 2>&1; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION) ;fi)
	golangci-lint run ./...

## Build macos binaries
darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS) -s -w" -o  $(BIN_DIR)/amd64 $(SRC_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS) -s -w" -o  $(BIN_DIR)/arm64 $(SRC_DIR)
	lipo -create bin/amd64 bin/arm64 -output $(BINARY)

## Run tests for my project
test:
	go test -v ./...

## Install Binary and Assets to Workflow Directory
install: build
	@(cp $(ASSETS)  $(WORKFLOW_DIR)/)

## create workflow artifact
package: darwin
	@(if [ ! -e $(ARTIFACT_DIR) ]; then mkdir $(ARTIFACT_DIR) ; fi)
	@(cp $(ASSETS) $(ARTIFACT_DIR))
	@(zip -j $(ARTIFACT_NAME) $(ARTIFACT_DIR)/*)

## GitHub Release and uploads artifacts
release: package
	@(if ! type ghr >/dev/null 2>&1; then go get -u github.com/tcnksm/ghr ;fi)
	@ghr -replace $(VERSION) $(ARTIFACT_NAME)

## Clean Binary
clean:
	rm -f $(BIN_DIR)/*
	rm -f $(ARTIFACT_DIR)/*

docker-test:
	go mod vendor
	docker run --rm -it -v $(PWD):/usr/src/myapp -w /usr/src/myapp golang:1.14 bash -c "./setup-test-dir.sh && make test"

## Show help
help:
	@(if ! type make2help >/dev/null 2>&1; then go get -u github.com/Songmu/make2help/cmd/make2help ;fi)
	@make2help $(MAKEFILE_LIST)

.PHONY: build test lint fmt darwin release package clean help
