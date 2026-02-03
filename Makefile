.PHONY: help test test-verbose test-coverage generate build run clean lint golangci-lint

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

test: ## Run all tests
	go test ./...

test-verbose: ## Run all tests with verbose output
	go test ./... -v

test-coverage: ## Run all tests with coverage report
	go test ./... -cover

test-coverage-html: ## Generate HTML coverage report
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

generate: ## Generate mocks using go generate
	go generate ./...

build: ## Build the application
	go build -o greeting-api

run: build ## Build and run the application
	./greeting-api

dev: ## Run the application without building binary
	go run main.go

clean: ## Clean build artifacts and generated files
	rm -f greeting-api coverage.out coverage.html
	rm -rf mock/

deps: ## Download dependencies
	go mod download
	go mod tidy

install-tools: ## Install development tools
	go install github.com/golang/mock/mockgen@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

golangci-lint: ## Run golangci-lint
	$$(go env GOPATH)/bin/golangci-lint run ./...

golangci-lint-fix: ## Run golangci-lint with auto-fix
	$$(go env GOPATH)/bin/golangci-lint run --fix ./...

lint: fmt vet golangci-lint ## Run all linters

all: clean generate lint test build ## Run all steps
