.PHONY: help install dev test lint clean

help:  ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install:  ## Install all dependencies
	pip install -e ".[dev]"
	cd frontend && npm install

dev:  ## Start dev servers
	docker-compose up -d redis postgres minio
	uvicorn backend.main:app --reload --port 8000 &
	cd frontend && npm run dev

test:  ## Run tests
	pytest tests/ -v --cov=backend

lint:  ## Run linters
	ruff check backend/ tests/
	ruff format --check backend/ tests/
	mypy backend/

format:  ## Format code
	ruff format backend/ tests/
	ruff check --fix backend/ tests/

clean:  ## Clean build artifacts
	find . -type d -name __pycache__ -exec rm -rf {} +
	find . -type d -name "*.egg-info" -exec rm -rf {} +
	rm -rf .pytest_cache .mypy_cache .ruff_cache dist build

models-download:  ## Download ML models
	bash ml/scripts/download_models.sh

docker-up:  ## Start all services with Docker
	docker-compose up -d

docker-down:  ## Stop all services
	docker-compose down
