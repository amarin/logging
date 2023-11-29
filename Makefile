# Top level Makefile for the entire project

PROJECT_NAME := "logging"
PKG := "git.media-tel.ru/railgo/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

# Go parameters
GOCMD=GO111MODULE=on GOFLAGS=-mod=vendor go
DEPCMD=$(GOCMD) mod
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Detect GOOS. By default build release for current OS. Redefine GOOS to cross-build.
GOOS ?= $(shell go env GOOS)

.PHONY: all dep build clean test coverage coverhtml lint build_path linux_ftpd

build: clean deps ## (Re)build tools binaries

test: ## Run unit tests
	@GOPATH=${GOPATH}:`pwd` $(GOCMD) test -v ./...

all: test

lint: ## Lint the files
	@golangci-lint run

deps: tidy ## Get the dependencies
	$(DEPCMD) vendor

.PHONY: tidy
tidy: ## Fix dependencies records in go.mod
	@go mod tidy

race: deps ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: deps ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

clean: ## Remove previous build
	@rm -rf $(BUILD_PATH)

patch_release: ## make new patch version
	@bumpversion patch

minor_release: ## make new minor version
	@bumpversion minor

## Display this help screen
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
