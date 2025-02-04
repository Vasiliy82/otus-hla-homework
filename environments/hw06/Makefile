# Имя файла для Docker Compose
DOCKER_COMPOSE_FILE ?= compose.yml
POSTGRES_HOST ?= postgres
POSTGRES_PORT ?= 5432
POSTGRES_DATABASE ?= hw
POSTGRES_USER ?= app_hw
POSTGRES_PASSWORD ?= Passw0rd
KAFKA_TOPIC_NAME_1 ?= mysocnet-cache-feed-changed
KAFKA_TOPIC_NAME_2 ?= mysocnet-posts-post-modified
KAFKA_TOPIC_NAME_3 ?= mysocnet-posts-follower-notify

# Задача: запуск всех контейнеров
.PHONY: up
up:
	@echo "Запуск контейнеров...";
	docker compose run --rm --entrypoint "sh -c 'rm -rf /var/lib/grafana/* && tar -vxf /tmp/configuration.tar.gz -C /'" grafana;
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d;
	docker compose -f $(DOCKER_COMPOSE_FILE) exec kafka /bin/bash -c "unset KAFKA_OPTS && kafka-topics --create --topic $(KAFKA_TOPIC_NAME_1) --partitions 3  --replication-factor 1 --bootstrap-server kafka:29092"
	docker compose -f $(DOCKER_COMPOSE_FILE) exec kafka /bin/bash -c "unset KAFKA_OPTS && kafka-topics --create --topic $(KAFKA_TOPIC_NAME_2) --partitions 3  --replication-factor 1 --bootstrap-server kafka:29092"
	docker compose -f $(DOCKER_COMPOSE_FILE) exec kafka /bin/bash -c "unset KAFKA_OPTS && kafka-topics --create --topic $(KAFKA_TOPIC_NAME_3) --partitions 3  --replication-factor 1 --bootstrap-server kafka:29092"
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

.PHONY: pg-backup
pg-backup:
	@echo "Резервное копирование БД..."
	docker compose exec postgres pg_dump -U app_hw -d hw | gzip > ~/Projects/2024-11-01-hw_backup.gz
	@echo "завершено"
