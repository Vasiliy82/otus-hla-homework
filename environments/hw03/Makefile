# Имя файла для Docker Compose
DOCKER_COMPOSE_FILE ?= compose.yml
POSTGRES_HOST ?= postgres
POSTGRES_PORT ?= 5432
POSTGRES_DATABASE ?= hw
POSTGRES_USER ?= app_hw
POSTGRES_PASSWORD ?= Passw0rd

# Задача: запуск всех контейнеров
.PHONY: up
up:
	@echo "Запуск контейнеров...";
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d;
	$(MAKE) migrate;

# Задача: накатывание миграций с использованием goose
.PHONY: migrate
migrate:
	@echo "Накатывание миграций...";
	docker compose run --rm goose goose -dir /app/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@${POSTGRES_HOST}:${POSTGRES_PORT}/$(POSTGRES_DATABASE)?sslmode=disable" up

# Задача: остановка всех контейнеров
.PHONY: down
down:
	@echo "Остановка всех контейнеров..."
	docker compose -f $(DOCKER_COMPOSE_FILE) down

# Задача: удаление всех контейнеров и томов
.PHONY: destroy
destroy:
	@echo "Остановка и удаление всех контейнеров и томов..."
	docker compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans
