SERVER_OUT := "syphon.bin"
SERVER_PKG_BUILD := "./cmd/syphon"
GIT_TAG := $(shell git describe --always --abbrev=0 --tags)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT := $(shell git log --pretty=format:'%h' -n 1)
METRIC_FILE := "metrics.txt"
METRIC_REC_FILE := "metrics_received.txt"
CONFIG := "./assets/config.toml"

.PHONY: all build run promiscuous

all: build

docker-build:
	@docker build -t butlerx/syphon:latest --build-arg VERSION="$(GIT_TAG)-$(GIT_BRANCH).$(GIT_COMMIT)" build

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build binary
	@go build -i -v -o $(SERVER_OUT) -ldflags "-X main.version=$(GIT_TAG)-$(GIT_BRANCH).$(GIT_COMMIT)" $(SERVER_PKG_BUILD)

clean: ## Remove previous builds
	@rm $(SERVER_OUT) $(METRIC_FILES) $(METRIC_REC_FILES)

run: dep ## Run server
	@go run $(SERVER_PKG_BUILD) --config $(CONFIG)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

promiscuous: build ## give minary ability to listen promiscuously
	sudo setcap cap_net_raw=ep $(SERVER_OUT)

test: ## Run test on code
	@go test ./...
