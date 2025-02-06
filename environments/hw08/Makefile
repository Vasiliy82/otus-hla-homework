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

# Параметры для JMeter
JMETER_IMAGE ?= justb4/jmeter:latest
JMETER_FILE ?= stresstest.jmx
RESULTS_DIR ?= ./jmeter_results
HTML_REPORT_DIR ?= ./jmeter_html_reports
SOCNET_PG_CONFIG ?= ./backend/socnet-pg.yaml
SOCNET_TARANTOOL_CONFIG ?= ./backend/socnet-t.yaml

# Массив потоков
THREADS_LIST=100
TEST_DURATION=5

# Задача: запуск всех контейнеров
.PHONY: up
up:
	@echo "Запуск контейнеров...";
	mkdir -p ./tarantool/data;
	chmod -R 777 ./tarantool/data;
	docker compose run --rm --entrypoint "sh -c 'rm -rf /var/lib/grafana/* && tar -vxf /tmp/configuration.tar.gz -C /'" grafana;
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d --build;
	docker compose -f $(DOCKER_COMPOSE_FILE) exec kafka /bin/bash -c "unset KAFKA_OPTS && kafka-topics --create --topic $(KAFKA_TOPIC_NAME_1) --partitions 3  --replication-factor 1 --bootstrap-server kafka:29092";
	docker compose -f $(DOCKER_COMPOSE_FILE) exec kafka /bin/bash -c "unset KAFKA_OPTS && kafka-topics --create --topic $(KAFKA_TOPIC_NAME_2) --partitions 3  --replication-factor 1 --bootstrap-server kafka:29092";
	docker compose -f $(DOCKER_COMPOSE_FILE) exec kafka /bin/bash -c "unset KAFKA_OPTS && kafka-topics --create --topic $(KAFKA_TOPIC_NAME_3) --partitions 3  --replication-factor 1 --bootstrap-server kafka:29092";
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

# Задача: запуск нагрузочного теста с HTML-отчетами
.PHONY: run-tests
run-tests:
	@echo "Запуск нагрузочного теста для PostgreSQL и Tarantool..."

	# Создаем директории для отчетов
	mkdir -p $(RESULTS_DIR)/postgresql
	mkdir -p $(RESULTS_DIR)/tarantool
	mkdir -p $(HTML_REPORT_DIR)/postgresql
	mkdir -p $(HTML_REPORT_DIR)/tarantool


	cp $(SOCNET_PG_CONFIG) ./backend/socnet.yaml;
	docker compose -f $(DOCKER_COMPOSE_FILE) down dialogs;
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d dialogs;

	# Тестирование для PostgreSQL
	for THREADS in $(THREADS_LIST); do \
		echo "Тестирование PostgreSQL с потоками: $$THREADS"; \
		docker run \
		    --add-host=host.docker.internal:host-gateway \
			--rm -v $$(pwd):/jmeter -w /jmeter $(JMETER_IMAGE) \
			-n -t $(JMETER_FILE) \
			-JHOSTNAME=host.docker.internal \
			-JPORT=8081 \
			-JTHREADS=$$THREADS \
			-JDURATION=$(TEST_DURATION) \
			-l $(RESULTS_DIR)/postgresql/result_$$THREADS.jtl \
			-e -o $(HTML_REPORT_DIR)/postgresql/report_$$THREADS; \
	done

	cp $(SOCNET_TARANTOOL_CONFIG) ./backend/socnet.yaml;
	docker compose -f $(DOCKER_COMPOSE_FILE) down dialogs;
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d dialogs;

	# Тестирование для Tarantool
	for THREADS in $(THREADS_LIST); do \
		echo "Тестирование Tarantool с потоками: $$THREADS"; \
		docker run \
		    --add-host=host.docker.internal:host-gateway \
		    --rm -v $$(pwd):/jmeter -w /jmeter $(JMETER_IMAGE) \
			-n -t $(JMETER_FILE) \
			-JHOSTNAME=host.docker.internal \
			-JPORT=8081 \
			-JTHREADS=$$THREADS \
			-JDURATION=$(TEST_DURATION) \
			-l $(RESULTS_DIR)/tarantool/result_$$THREADS.jtl \
			-e -o $(HTML_REPORT_DIR)/tarantool/report_$$THREADS; \
	done

	@echo "Нагрузочные тесты завершены. HTML-отчеты сохранены в $(HTML_REPORT_DIR)."
