# --- Variables ---
POSTGRES_USER ?= app_hw
POSTGRES_PASSWORD ?= Passw0rd
POSTGRES_DATABASE ?= hw

export PATH := $(PWD)/bin:$(PATH)
export SHELL := bash

# ~~~ Docker Compose Commands ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

up: docker-up wait-for-postgres migrate-up frontend-dev ## Start all services
down: docker-stop frontend-stop  ## Stop all services
destroy: docker-teardown clean  ## Remove containers and volumes

docker-up: ## Start Docker Compose
	@ docker compose up -d --build

docker-stop: ## Stop Docker Compose
	@ docker compose down

docker-teardown: ## Remove containers and volumes
	@ docker compose down --remove-orphans -v

wait-for-postgres: ## Wait for PostgreSQL to be ready
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker compose exec postgres pg_isready -U $(POSTGRES_USER); do \
		echo "PostgreSQL is unavailable - sleeping"; \
		sleep 2; \
	done
	@echo "PostgreSQL is up - executing command"

# ~~~ Migrations ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

migrate-up: ## Apply migrations
	@ docker compose exec backend goose -dir /app/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres/$(POSTGRES_DATABASE)?sslmode=disable" up

migrate-down: ## Rollback migrations
	@ docker compose exec backend goose -dir /app/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres/$(POSTGRES_DATABASE)?sslmode=disable" down

# ~~~ Frontend Commands ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

frontend-dev: ## Start frontend in dev mode
	@ docker compose up frontend

frontend-stop: ## Stop frontend
	@ docker compose stop frontend

clean: clean-docker ## Clean up all artifacts and dangling images
clean-docker: 
	@ docker image prune -f
	