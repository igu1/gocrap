# Simple Makefile for Go project

APP_NAME := server
PKG := github.com/igu1/gocrap
BIN_DIR := bin
CMD_DIR := ./cmd/server

.PHONY: help tidy fmt lint test build run clean docker-build docker-run

help:
	@echo "Targets:"
	@echo "  tidy         - go mod tidy"
	@echo "  fmt          - go fmt"
	@echo "  lint         - golangci-lint (if installed)"
	@echo "  test         - run unit tests"
	@echo "  build        - build $(APP_NAME) to $(BIN_DIR)/$(APP_NAME)"
	@echo "  run          - run the server (localhost:8080)"
	@echo "  clean        - remove $(BIN_DIR)"
	@echo "  docker-build - build docker image $(APP_NAME):latest"
	@echo "  docker-run   - run docker image mapping 8080:8080"

 tidy:
	go mod tidy

 fmt:
	go fmt ./...

 lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed, skipping"

 test:
	go test ./...

 build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

 run:
	go run $(CMD_DIR)

 clean:
	rm -rf $(BIN_DIR)

 docker-build:
	docker build -t $(APP_NAME):latest .

 docker-run:
	docker run --rm -p 8080:8080 $(APP_NAME):latest