APP_NAME := open-sandbox
BUILD_DIR := build

.PHONY: build test

build:
	go build -o $(BUILD_DIR)/$(APP_NAME).exe ./cmd/server

test:
	go test ./...
