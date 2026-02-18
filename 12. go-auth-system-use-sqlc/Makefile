# ------------------------------------------------------------------------------
# 1. Environment Configuration
# ------------------------------------------------------------------------------
include .env
export

# DB_URL for host-based tools (migrate, sqlc)
# Ensure POSTGRES_HOST=127.0.0.1 in your .env for IPv4 compatibility
DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL_MODE)

# ------------------------------------------------------------------------------
# 2. Command Groups (Targets)
# ------------------------------------------------------------------------------
.PHONY: tidy swagger sqlc migrate-up migrate-down new-migration docker-up docker-down dev clean help

# 3.1 Detect OS for cross-platform support
ifeq ($(OS),Windows_NT)
    # Windows timeout command: /t 5 seconds, /nobreak to prevent keypress interruption
    SLEEP := timeout /t 5 /nobreak > NUL
else
    # Linux/Unix sleep command
    SLEEP := sleep 5
endif

# 3.2 Optimized Development Workflow
# Standard Flow: Infrastructure -> Migrations -> CodeGen -> Documentation -> Run
dev: docker-up
	@echo "[INFO] Waiting for Database to be ready (5s)..."
	@$(SLEEP)
	@echo "[DB] Running migrations..."
	$(MAKE) migrate-up
	@echo "[SQLC] Generating type-safe Go code..."
	$(MAKE) sqlc
	@echo "[SWAGGER] Generating API documentation..."
	$(MAKE) swagger
	@echo "[GO] Tidying modules..."
	$(MAKE) tidy
	@echo "[SERVER] Starting API server on port $(API_SERVER_PORT)..."
	go run cmd/api/main.go

# ------------------------------------------------------------------------------
# 4. Database & Tooling Details
# ------------------------------------------------------------------------------

# Create a new migration file (Usage: make new-migration name=init_schema)
new-migration:
	@echo "[DB] Creating new migration files for: $(name)..."
	migrate create -ext sql -dir db/migration -seq $(name)

migrate-up:
	@echo "[DB] Applying all up migrations..."
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrate-down:
	@echo "[DB] Rolling back migrations..."
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	@echo "[SQLC] Generating Go code from SQL queries..."
	sqlc generate
	@echo "[SQLC] Generation completed successfully."

swagger:
	@echo "[SWAGGER] Initializing Swagger documentation..."
	@echo "[SWAGGER] Documentation generated at ./api/docs"

docker-up:
	@echo "[DOCKER] Starting infrastructure containers (Postgres & Redis)..."
	docker compose up -d db redis redisinsight
	@echo "[DOCKER] Containers are running in detached mode."

docker-down:
	@echo "[DOCKER] Stopping and removing containers..."
	docker compose down

clean:
	@echo "[CLEAN] Tearing down environment..."
	docker compose down -v
	rm -rf api/docs
	@echo "[CLEAN] Cleaned up all containers, volumes, and generated docs."

help:
	@echo "Available commands:"
	@echo "  make dev            : Run the full development stack"
	@echo "  make migrate-up     : Apply database migrations"
	@echo "  make sqlc           : Generate Go code from SQL"
	@echo "  make swagger        : Generate API docs"
	@echo "  make docker-up      : Start DB and Redis containers"
	@echo "  make clean          : Remove all generated files and volumes"
