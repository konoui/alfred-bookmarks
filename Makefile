VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main. version=$(VERSION)' -X 'main. revision=$(REVISION)'
SRC_DIR := ./
BIN_DIR := bin
BINARY := bin/alfred-firefox-bookmarks
WORKFLOW_DIR := "$${HOME}/Library/Application Support/Alfred 3/Alfred.alfredpreferences/workflows/user.workflow.7C42A657-124F-46B8-89EE-7A1C06594E13"
ASSETS_DIR := assets
ARTIFACT_DIR := .artifact
ARTIFACT := ${ARTIFACT_DIR}/alfred-firefox-bookmarks.alfredworkflow

## Build binaries on your environment
build: deps
	CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o ${BINARY} ./${SRC_DIR}

## Setup
setup:
	#installing golint
	@(if ! type golint >/dev/null 2>&1; then go get -u golang.org/x/lint/golint ;fi)
	#installing golangci-lint
	@(if ! type golangci-lint >/dev/null 2>&1; then go get -u github.com/golangci/golangci-lint/cmd/golangci-lint ;fi)
	#installing dep
	@(if ! type dep >/dev/null 2>&1; then go get -u github.com/golang/dep/cmd/dep ;fi)
	#installing goimports
	@(if ! type goimports >/dev/null 2>&1; then go get -u golang.org/x/tools/cmd/goimports ;fi)
	#installing ghr
	@(if ! type ghr >/dev/null 2>&1; then go get -u github.com/tcnksm/ghr ;fi)
	#installing make2help
	@(if ! type make2help >/dev/null 2>&1; then go get -u github.com/Songmu/make2help/cmd/make2help ;fi)

## Install dependencies to vendor
deps: setup
	dep ensure #-vendor-only

## Update dependencies
update: setup
	dep ensure -update

## Format source codes
fmt: deps
	goimports -w $$(go list -f {{.Dir}} ./... | grep -v /vendor/)

## Lint
lint: deps
	golangci-lint run ./...

## Build linux binaries
linux: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BINARY} ./${SRC_DIR}

## Run tests for my project
test: deps
	go test -v ./...

## Initialize directory
init:
	@(if [ ! -e ${SRC_DIR} ]; then mkdir ${SRC_DIR}; fi)
	@(if [ ! -e ${BIN_DIR} ]; then mkdir ${BIN_DIR}; fi)
	@(if [ ! -e Gopkg.lock ]; then dep init; fi)

## Install Binary and Assets to Workflow Directory
install: build
	@(cp ${BINARY} ${WORKFLOW_DIR}/)
	@(cp ${ASSETS_DIR}/*  ${WORKFLOW_DIR}/)

release: build
	@(if [ ! -e ${ARTIFACT_DIR} ]; then mkdir ${ARTIFACT_DIR} ; fi)
	@(cp ${BINARY} ${ARTIFACT_DIR})
	@(cp ${ASSETS_DIR}/* ${ARTIFACT_DIR})
	@(zip -j ${ARTIFACT} ${ARTIFACT_DIR}/*)
	@(export GITHUB_TOKEN=$(shell aws secretsmanager get-secret-value --secret-id github_token --query 'SecretString' --output text) ;\
	ghr -replace ${VERSION} ${ARTIFACT})

## Clean Binary
clean:
	rm -f ${BIN_DIR}/*

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: build setup deps update test lint fmt linux init clean help
