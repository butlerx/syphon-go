SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
.DEFAULT_GOAL := help
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

GO111MODULES=on
GIT_TAG := $(shell git describe --always --abbrev=0 --tags)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT := $(shell git log --pretty=format:'%h' -n 1)
VERSION="$(GIT_TAG)-$(GIT_BRANCH).$(GIT_COMMIT)"
IMAGE_ID="butlerx/syphon:$(VERSION)"

SYPHON := "syphon.bin"
SYPHON_PKG_BUILD := "./cmd/syphon"
RELEASE_ZIP := "syphon.zip"

.PHONY: all docker build release
all: build docker release
docker: .image-id ## Build Docker container
build: $(SYPHON) ## Build Binary
release: $(RELEASE_ZIP) ## Package release artifact

.image-id: build/package/Dockerfile
	@echo "🍳 Building container $(IMAGE_ID)"
	docker build --tag="$(IMAGE_ID)" --build-arg $(VERSION) -f build/package/Dockerfile .
	echo "$(IMAGE_ID)" > .image-id

$(SYPHON): dep
	@echo "🍳 Building $(SYPHON)"
	@go build -i -v -o $(SYPHON) -ldflags "-X main.version=$(GIT_TAG)-$(GIT_BRANCH).$(GIT_COMMIT)" $(SYPHON_PKG_BUILD)

$(RELEASE_ZIP): $(SYPHON)
	@echo "🍳 Building $(RELEASE_ZIP)"
	zip --junk-paths $(RELEASE_ZIP) $(SYPHON) README.md

.PHONY:clean
clean: ## Remove previous builds
	@echo "🧹 Cleaning old build"
	@go clean
	@rm -f $(SYPHON) $(RELEASE_ZIP) metrics.txt metrics_received.txt

.PHONY: dep
dep: ## go get all dependencies
	@echo "🛎 Updatig Dependencies"
	@go get -v -d ./...

.PHONY: run
run: dep ## Compiles and runs server
	@go run -race $(SYPHON_PKG_BUILD) --config configs/config.toml

.PHONY: help
help:  ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: promiscuous
promiscuous: $(SYPHON)  ## give binary ability to listen promiscuously
	sudo setcap cap_net_raw=ep $(SYPHON)

.PHONY: test
test: ## Runs go test with default values
	@echo "🍜 Testing $(SYPHON)"
	@go test -v -count=1 -race ./...
