## ────────────────────────────────────────────────────────────────────────────
##  MediConnect ID — Makefile
## ────────────────────────────────────────────────────────────────────────────

APP_NAME   := mediconnect-backend
MAIN_PATH  := ./cmd/server
BIN_DIR    := ./bin
COVER_OUT  := coverage.out

.PHONY: all run build test lint fmt tidy \
        docker-up docker-down docker-build \
        migrate seed clean help

## ── Development ──────────────────────────────────────────────────────────────

run: ## Run the application locally (hot env from .env)
	go run $(MAIN_PATH)/...

build: ## Compile a production binary into ./bin/
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PATH)/...
	@echo "Binary: $(BIN_DIR)/$(APP_NAME)"

## ── Quality ──────────────────────────────────────────────────────────────────

test: ## Run all tests with race detector and coverage report
	go test -v -race -count=1 -coverprofile=$(COVER_OUT) ./...
	go tool cover -func=$(COVER_OUT) | tail -1

lint: ## Run golangci-lint (install: https://golangci-lint.run/usage/install/)
	golangci-lint run ./...

fmt: ## Format all Go source files
	gofmt -s -w .

tidy: ## Tidy and verify Go module dependencies
	go mod tidy
	go mod verify

## ── Docker ───────────────────────────────────────────────────────────────────

docker-up: ## Start all services (PostgreSQL, Redis, RabbitMQ, App) via Docker Compose
	docker compose up -d --build
	@echo "Services running. Backend → http://localhost:8080"

docker-down: ## Stop and remove all containers
	docker compose down

docker-build: ## Build only the app Docker image
	docker build -t $(APP_NAME):local .

## ── Database ─────────────────────────────────────────────────────────────────

migrate: ## Apply ALL SQL migrations (migrations/*.sql) in order to mediconnect-db
	@echo ">>> Running all migrations in migrations/ ..."
	@for f in $(sort $(wildcard migrations/*.sql)); do \
		echo "  Applying $$f ..."; \
		docker exec -i mediconnect-db \
			psql -U mediconnect_user -d mediconnect_db \
			< $$f || exit 1; \
	 done
	@echo "✅ All migrations complete."

seed: ## Run Go seeder to populate the database with sample data
	@echo ">>> Running Go seeder ..."
	go run ./cmd/seed/...
	@echo "✅ Seeding complete."

## ── Utility ──────────────────────────────────────────────────────────────────

clean: ## Remove build artefacts
	rm -rf $(BIN_DIR) $(COVER_OUT)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
