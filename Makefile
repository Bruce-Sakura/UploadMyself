.PHONY: help build dev test lint clean docker-up docker-down

help:  ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build:  ## Build backend binary
	cd backend && go build -o ../bin/uploadmyself .

dev:  ## Start dev (backend + frontend + deps)
	docker-compose up -d redis postgres minio
	cd backend && go run . &
	cd frontend && npm run dev

test:  ## Run tests
	cd backend && go test ./... -v -coverprofile=coverage.out

lint:  ## Run linters
	cd backend && go vet ./...
	cd backend && golangci-lint run ./...

format:  ## Format Go code
	cd backend && gofmt -w .
	cd backend && goimports -w .

clean:  ## Clean build artifacts
	rm -rf bin/
	cd backend && go clean -cache
	rm -rf frontend/node_modules/ frontend/dist/

tidy:  ## Tidy Go modules
	cd backend && go mod tidy

models-download:  ## Download ML models
	bash ml/scripts/download_models.sh

docker-up:  ## Start all services with Docker
	docker-compose up -d

docker-down:  ## Stop all services
	docker-compose down

# Frontend
frontend-install:  ## Install frontend deps
	cd frontend && npm install

frontend-dev:  ## Start frontend dev server
	cd frontend && npm run dev

frontend-build:  ## Build frontend
	cd frontend && npm run build
