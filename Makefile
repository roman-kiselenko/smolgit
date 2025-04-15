GOCMD=go
GOTEST=$(GOCMD) test
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
HASH := $(shell git rev-parse --short HEAD)
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

PROJECT_NAME := smolgit
LINTER_BIN ?= golangci-lint

.PHONY: all test build clean run lint /bin/$(LINTER_BIN)

all: help

## Build:
build: ## Build all the binaries and put the output in bin/
	$(GOCMD) build -ldflags "-X main.version=$(BRANCH)-$(HASH)" -o bin/$(PROJECT_NAME) .

build-docker: ## Build an image
	docker build -t $(PROJECT_NAME) .

bin/$(LINTER_BIN):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v2.1.1

## Clean:
clean: ## Remove build related file
	@-rm -fr ./bin

## Lint:
lint: ./bin/$(LINTER_BIN) ## Lint sources with golangci-lint
	./bin/$(LINTER_BIN) run

## Run:
run: clean build ## Run the smolgit `make run`
	./bin/$(PROJECT_NAME) $(ARGS)

run-docker: ## Run smolgit in the container
	docker run -it -p 3080:3080 -p 3081:3081 -v $(PWD)/:/etc/smolgit $(PROJECT_NAME)

config-docker: ## Generate smolgit config
	docker run -it $(PROJECT_NAME) config > config.yaml

config: ## Generate default config
	./bin/$(PROJECT_NAME) config > ./config.yaml

## Test:
test: ## Run the tests of the smolgit
	$(GOTEST) -v -race ./...

integration-test: build ## Run the bats tests of the smolgit
	@bats ./test/test.bats

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
