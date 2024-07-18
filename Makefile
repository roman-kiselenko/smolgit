GOCMD=go
GOTEST=$(GOCMD) test
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
HASH := $(shell git rev-parse --short HEAD)
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build clean run

PROJECT_NAME := smolgit

all: help

## Build:
build: ## Build all the binaries and put the output in bin/
	$(GOCMD) build -ldflags "-X main.version=$(BRANCH)-$(HASH)" -o bin/smolgit .

## Config:
config: ## Generate default config
	./bin/smolgit config > ./config.yaml

## Clean:
clean: ## Remove build related file
	@-rm -fr ./bin

## Run:
run: clean build config ## Run the smolgit `make run`
	./bin/smolgit $(ARGS)

## Test:
test: ## Run the tests of the smolgit
	$(GOTEST) -v -race ./...

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