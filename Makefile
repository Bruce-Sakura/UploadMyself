.PHONY: help build dev test lint clean

help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build binary
	cd backend && go build -o ../bin/uploadmyself .

dev: ## Run dev server
	cd backend && go run .

test: ## Run tests
	cd backend && go test ./... -v -count=1

lint: ## Lint Go code
	cd backend && go vet ./...

format: ## Format code
	cd backend && gofmt -w .

clean:
	rm -rf bin/

docker-up: ## Start deps (Redis + PG)
	docker-compose up -d redis postgres

docker-down: ## Stop deps
	docker-compose down
