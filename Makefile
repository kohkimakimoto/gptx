.DEFAULT_GOAL := help

SHELL := bash
PATH := $(CURDIR)/.dev/gopath/bin:$(PATH)
VERSION := 0.0.0
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/kohkimakimoto/gptx/internal.CommitHash=$(COMMIT_HASH) -X github.com/kohkimakimoto/gptx/internal.Version=$(VERSION)"

# Load .env file if it exists.
ifneq (,$(wildcard ./.env))
  include .env
  export
endif


.PHONY: help
help: ## Show help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[/0-9a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'


.PHONY: bump
bump: ## Bump up version
	@PART=${PART}; \
	if [ -z "$$PART" ]; then \
		PART=minor; \
	fi; \
	if [ $$PART = "major" ]; then \
		perl -i.bak -pe 's/(VERSION := )(\d+)(\.(\d+)\.(\d+))/$$1 . ($$2 + 1) .  ".0.0"/e' Makefile; \
	elif [ $$PART = "minor" ]; then \
		perl -i.bak -pe 's/(VERSION := (\d+)\.)(\d+)(\.(\d+))/$$1 . ($$3 + 1) . ".0"/e' Makefile; \
	elif [ $$PART = "patch" ]; then \
	  	perl -i.bak -pe 's/(VERSION := (\d+\.\d+\.))(\d+)/$$1 . ($$3 + 1)/e' Makefile; \
	else \
		echo "Invalid part: $$PART"; exit 1; \
	fi && rm Makefile.bak
	@new_version=$$(perl -ne 'print $$1 if /VERSION := (\d+\.\d+\.\d+)/' Makefile) && \
	git commit -am "Bump up version to $$new_version" && \
	git tag "v$$new_version"


.PHONY: dev/setup
dev/setup: ## Setup development environment
	@mkdir -p .dev/gopath
	@export GOPATH=$(CURDIR)/.dev/gopath && \
		go install honnef.co/go/tools/cmd/staticcheck@latest && \
		go install github.com/Songmu/goxz/cmd/goxz@latest && \
		go install github.com/axw/gocov/gocov@latest && \
		go install github.com/matm/gocov-html/cmd/gocov-html@latest && \
		go install go.etcd.io/bbolt/cmd/bbolt@latest


.PHONY: dev/clean
dev/clean: ## Clean up development environment
	@export GOPATH=$(CURDIR)/.dev/gopath && go clean -modcache
	@rm -rf .dev


.PHONY: format
format: ## Format source code
	@go fmt ./...


.PHONY: test
test: ## Run tests
	@go test -race -timeout 30m -cover ./...


.PHONY: test/verbos
test/verbose: ## Run tests with verbose outputting.
	@go test -race -timeout 30m -cover -v ./...


.PHONY: test/coverage
test/coverage: ## Run tests with coverage report
	@mkdir -p .dev
	@go test -race -timeout 30m -cover ./... -coverprofile=.dev/coverage.out
	@gocov convert .dev/coverage.out | gocov-html > .dev/coverage.html


.PHONY: lint
lint: ## Static code analysis
	@go vet ./...
	@staticcheck ./...


.PHONY: open-coverage-html
open-coverage-html: ## Open coverage report
	@open .dev/coverage.html


.PHONY: build
build: ## build dev binary
	@mkdir -p .dev/build/dev
	@go build -ldflags=$(BUILD_LDFLAGS) -o .dev/build/dev/gptx ./cmd/gptx


.PHONY: build/release
build/release: ## build release binary
	@mkdir -p .dev/build/release
	@goxz -n gptx -pv=v$(VERSION) -os=linux,darwin -static -build-ldflags=$(BUILD_LDFLAGS) -d=.dev/build/release ./cmd/gptx


.PHONY: clean
clean: ## Clean generated files
	@rm -rf .dev/build
	@rm -rf .dev/coverage.html
	@rm -rf .dev/coverage.out


# check variable definition
guard-%:
	@if [[ -z '${${*}}' ]]; then echo 'ERROR: variable $* not set' && exit 1; fi
