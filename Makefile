SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

GIT_TAG := $(shell git describe --always --abbrev=0 --tags)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT := $(shell git log --pretty=format:'%h' -n 1)
VERSION="$(GIT_TAG)-$(GIT_BRANCH).$(GIT_COMMIT)"

SERVER_OUT := "syphon.bin"
SERVER_PKG_BUILD := "./cmd/syphon"
RELEASE_ZIP := "syphon.zip"

.PHONY: all docker run promiscuous test clean
all: build
docker: .image-id ## Build Docker container
build: $(SERVER_OUT) ## Build binary
release: $(RELEASE_ZIP) ## Package release artifact

.image-id:
	image_id="butlerx/syphon:$$(pwgen -1)"
	docker build --tag="$${image_id}" --build-arg $(VERSION) -f build/package/Dockerfile .
	echo "$${image_id}" > .image-id

$(SERVER_OUT): dep
	@go build -i -v -o $(SERVER_OUT) -ldflags "-X main.version=$(GIT_TAG)-$(GIT_BRANCH).$(GIT_COMMIT)" $(SERVER_PKG_BUILD)

$(RELEASE_ZIP): build
	zip --junk-paths $(RELEASE_ZIP) $(SERVER_OUT) README.md

clean: ## Remove previous builds
	@go clean
	@rm -f $(SERVER_OUT) $(RELEASE_ZIP) metrics.txt metrics_received.txt

dep: ## Get the dependencies
	@go get -v -d ./...

run: dep ## Run server
	@go run $(SERVER_PKG_BUILD) --config configs/config.toml

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

promiscuous: build ## give minary ability to listen promiscuously
	sudo setcap cap_net_raw=ep $(SERVER_OUT)

test: ## Run test on code
	@go test ./...
