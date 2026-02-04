.PHONY: help test test-verbose test-coverage generate build run clean lint golangci-lint docker-build docker-run docker-stop compose-up compose-down compose-logs k8s-deploy k8s-delete k8s-status k8s-logs k8s-kibana

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

# Docker targets
DOCKER_IMAGE := greeting-api
DOCKER_TAG := latest

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	docker run -d --name $(DOCKER_IMAGE) -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-stop: ## Stop and remove Docker container
	docker stop $(DOCKER_IMAGE) && docker rm $(DOCKER_IMAGE)

compose-up: ## Start with docker-compose
	docker compose up -d --build

compose-down: ## Stop docker-compose
	docker compose down

compose-logs: ## Show docker-compose logs
	docker compose logs -f

# Kubernetes targets
k8s-deploy: docker-build ## Deploy to Kubernetes with Filebeat sidecar
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/elasticsearch.yaml
	kubectl apply -f k8s/kibana.yaml
	kubectl apply -f k8s/filebeat-configmap.yaml
	kubectl apply -f k8s/greeting-api.yaml
	@echo "Waiting for pods to be ready..."
	kubectl wait --for=condition=ready pod -l app=elasticsearch -n greeting --timeout=120s || true
	kubectl wait --for=condition=ready pod -l app=kibana -n greeting --timeout=120s || true
	kubectl wait --for=condition=ready pod -l app=greeting-api -n greeting --timeout=120s || true

k8s-delete: ## Delete Kubernetes resources
	kubectl delete -f k8s/ --ignore-not-found

k8s-status: ## Show Kubernetes status
	kubectl get all -n greeting

k8s-logs: ## Show greeting-api logs
	kubectl logs -f -l app=greeting-api -c greeting-api -n greeting

k8s-logs-filebeat: ## Show Filebeat sidecar logs
	kubectl logs -f -l app=greeting-api -c filebeat -n greeting

k8s-kibana: ## Get Kibana URL
	@echo "Kibana URL: http://localhost:$$(kubectl get svc kibana -n greeting -o jsonpath='{.spec.ports[0].nodePort}')"
	@echo "Or use port-forward: kubectl port-forward svc/kibana 5601:5601 -n greeting"
