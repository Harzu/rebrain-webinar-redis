SHELL := bash
.ONESHELL:
MAKEFLAGS += --no-builtin-rules
NOCACHE := $(if $(NOCACHE),"--no-cache")

export APP_NAME := rebrain-webinars
export DOCKER_REPOSITORY := redis
export VERSION := $(if $(VERSION),$(VERSION),$(if $(COMMIT_SHA),$(COMMIT_SHA),$(shell git rev-parse --verify HEAD)))

.PHONY: help
help: ## List all available targets with help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init: ## Prepare project for development
	@git config core.hooksPath .githooks
	@go mod tidy && make generate

.PHONY: build
build:
	@docker build ${NOCACHE} --pull -f ./build/Dockerfile -t ${DOCKER_REPOSITORY}/${APP_NAME}:${VERSION} .

.PHONY: run
run: ## Run develop docker-compose
	@docker-compose up app

run-replicas: ## Run develop app with some replicas
	REPLICAS_COUNT=3 docker-compose up app

.PHONY: stop
stop: ## Stop all develop containers
	@docker-compose down -v
