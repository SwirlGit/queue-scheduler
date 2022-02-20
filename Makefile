CONFIG_PATH ?= cmd/$(APP_NAME)/config.yaml

.PHONY: build
build:
	go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

.PHONY: build-qs-api
build-qs-api:
	make APP_NAME=qs-api build

.PHONY: build-qs-worker
build-qs-worker:
	make APP_NAME=qs-worker build

.PHONY: build-qs-checker
build-qs-checker:
	make APP_NAME=qs-checker build

.PHONY: build-all
build-all:
	make build-qs-api
	make build-qs-worker
	make build-qs-checker

.PHONY: run
run:
	CONFIG_PATH=$(CONFIG_PATH) go run cmd/$(APP_NAME)/main.go

.PHONY: run-qs-api
run-qs-api:
	make APP_NAME=qs-api run

.PHONY: run-qs-worker
run-qs-worker:
	make APP_NAME=qs-worker run

.PHONY: run-qs-checker
run-qs-checker:
	make APP_NAME=qs-checker run

.PHONY: lint
lint:
	golangci-lint run -v ./...