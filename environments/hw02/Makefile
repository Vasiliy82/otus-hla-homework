# Имя файла для Docker Compose
DOCKER_COMPOSE_FILE = docker-compose.yml

# Задача: запуск всех контейнеров
.PHONY: up
up:
	@echo "Запуск контейнеров..."
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

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
