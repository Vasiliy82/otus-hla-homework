# --- Variables ---
POSTGRES_USER ?= app_hw
POSTGRES_PASSWORD ?= Passw0rd
POSTGRES_DATABASE ?= hw

export PATH := $(PWD)/bin:$(PATH)
export SHELL := bash

# ~~~ Docker Compose Commands ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

up: docker-up ## Start all services
down: docker-down  ## Stop all services
destroy: docker-teardown clean  ## Remove containers and volumes

docker-up: ## Start Docker Compose
	@ docker compose up -d --build

docker-down: ## Stop Docker Compose
	@ docker compose down

docker-teardown: ## Remove containers and volumes
	@ docker compose down --remove-orphans -v

# ~~~ Migrations ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

migrate-up: ## Apply migrations
	@ docker compose exec backend goose -dir /app/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@${POSTGRES_HOST}:${POSTGRES_PORT}/$(POSTGRES_DATABASE)?sslmode=disable" up

migrate-down: ## Rollback migrations
	@ docker compose exec backend goose -dir /app/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@${POSTGRES_HOST}:${POSTGRES_PORT}/$(POSTGRES_DATABASE)?sslmode=disable" down

# ~~~ Clean Commands ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

clean: clean-docker ## Clean up all artifacts and dangling images
clean-docker:
	@ docker system prune -f --volumes
#	@docker stop $$(docker ps -q) || true
#	@docker rm $$(docker ps -a -q) || true
#	@docker rmi $$(docker images -q) || true
#	@docker volume rm $$(docker volume ls -q) || true
#	@docker network rm $$(docker network ls -q) || true
#	@docker system prune -a --volumes -f || true