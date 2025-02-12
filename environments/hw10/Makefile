# Имя файла для Docker Compose
DOCKER_COMPOSE_FILE ?= compose.yml
POSTGRES_HOST ?= postgres
POSTGRES_PORT ?= 5432
POSTGRES_DATABASE ?= hw
POSTGRES_USER ?= app_hw
POSTGRES_PASSWORD ?= Passw0rd

POSTGRES_COUNTERS_HOST ?= postgres-counters
POSTGRES_COUNTERS_PORT ?= 5432
POSTGRES_COUNTERS_DATABASE ?= hw
POSTGRES_COUNTERS_USER ?= app_hw
POSTGRES_COUNTERS_PASSWORD ?= Passw0rd

KAFKA_TOPIC_NAME_1 ?= mysocnet-cache-feed-changed
KAFKA_TOPIC_NAME_2 ?= mysocnet-posts-post-modified
KAFKA_TOPIC_NAME_3 ?= mysocnet-posts-follower-notify
KAFKA_TOPIC_NAME_4 ?= mysocnet-saga-bus

# Массив потоков
THREADS_LIST=100
TEST_DURATION=5

# Задача: запуск всех контейнеров
.PHONY: up
up:
	@echo "Запуск контейнеров...";
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d --build;
	docker compose -f $(DOCKER_COMPOSE_FILE) exec kafka /bin/bash -c "unset KAFKA_OPTS && kafka-topics --create --topic $(KAFKA_TOPIC_NAME_1) --partitions 3  --replication-factor 1 --bootstrap-server kafka:29092";
	docker compose -f $(DOCKER_COMPOSE_FILE) exec kafka /bin/bash -c "unset KAFKA_OPTS && kafka-topics --create --topic $(KAFKA_TOPIC_NAME_2) --partitions 3  --replication-factor 1 --bootstrap-server kafka:29092";
	docker compose -f $(DOCKER_COMPOSE_FILE) exec kafka /bin/bash -c "unset KAFKA_OPTS && kafka-topics --create --topic $(KAFKA_TOPIC_NAME_3) --partitions 3  --replication-factor 1 --bootstrap-server kafka:29092";
	docker compose -f $(DOCKER_COMPOSE_FILE) exec kafka /bin/bash -c "unset KAFKA_OPTS && kafka-topics --create --topic $(KAFKA_TOPIC_NAME_4) --partitions 3  --replication-factor 1 --bootstrap-server kafka:29092";
	$(MAKE) migrate;

# Задача: накатывание миграций с использованием goose
.PHONY: migrate
migrate:
	@echo "Накатывание миграций...";
	docker compose run -v $(CURDIR)/../../backend/migrations:/app/migrations --rm goose goose -dir /app/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@${POSTGRES_HOST}:${POSTGRES_PORT}/$(POSTGRES_DATABASE)?sslmode=disable" up
	@echo "Накатывание миграций counters...";
	docker compose run -v $(CURDIR)/../../counters/migrations:/app/migrations --rm goose goose -dir /app/migrations postgres "postgres://$(POSTGRES_COUNTERS_USER):$(POSTGRES_COUNTERS_PASSWORD)@${POSTGRES_COUNTERS_HOST}:${POSTGRES_COUNTERS_PORT}/$(POSTGRES_COUNTERS_DATABASE)?sslmode=disable" up

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

.PHONY: pg-backup
pg-backup:
	@echo "Резервное копирование БД..."
	docker compose exec postgres pg_dump -U app_hw -d hw | gzip > ~/Projects/2024-11-01-hw_backup.gz
	@echo "завершено"

