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

.PHONY: run-tests

run-tests:
	@echo "Запуск нагрузочного теста с последовательным отключением сервисов..."

	# Создаем директории для отчетов
	mkdir -p $(RESULTS_DIR)/step1
	mkdir -p $(RESULTS_DIR)/step2
	mkdir -p $(RESULTS_DIR)/step3
	mkdir -p $(RESULTS_DIR)/step4

	# Шаг 1: Тестирование на начальной конфигурации
	echo "Шаг 1: Тестирование системы без отказов";
	docker run \
	    --add-host=host.docker.internal:host-gateway \
	    --rm -v $$(pwd):/jmeter -w /jmeter $(JMETER_IMAGE) \
		-n -t $(JMETER_FILE) \
		-JHOSTNAME=host.docker.internal \
		-JPORT=8080 \
		-JTHREADS=1 \
		-JDURATION=$(TEST_DURATION) \
		-l $(RESULTS_DIR)/step1/result.jtl \
		-e -o $(RESULTS_DIR)/step1/report;

	# Шаг 2: Отключение postgres_slave1 и повторное тестирование
	echo "Отключение postgres_slave1 и повторное тестирование";
	docker compose kill postgres_slave1;
	docker run \
	    --add-host=host.docker.internal:host-gateway \
	    --rm -v $$(pwd):/jmeter -w /jmeter $(JMETER_IMAGE) \
		-n -t $(JMETER_FILE) \
		-JHOSTNAME=host.docker.internal \
		-JPORT=8080 \
		-JTHREADS=1 \
		-JDURATION=$(TEST_DURATION) \
		-l $(RESULTS_DIR)/step2/result.jtl \
		-e -o $(RESULTS_DIR)/step2/report;

	# Шаг 3: Отключение dialogs-1 и повторное тестирование
	echo "Отключение dialogs-1 и повторное тестирование";
	docker compose kill dialogs-1;
	docker run \
	    --add-host=host.docker.internal:host-gateway \
	    --rm -v $$(pwd):/jmeter -w /jmeter $(JMETER_IMAGE) \
		-n -t $(JMETER_FILE) \
		-JHOSTNAME=host.docker.internal \
		-JPORT=8080 \
		-JTHREADS=1 \
		-JDURATION=$(TEST_DURATION) \
		-l $(RESULTS_DIR)/step3/result.jtl \
		-e -o $(RESULTS_DIR)/step3/report;

	# Шаг 4: Отключение socnet-1 и повторное тестирование
	echo "Отключение socnet-1 и повторное тестирование";
	docker compose kill socnet-1;
	docker run \
	    --add-host=host.docker.internal:host-gateway \
	    --rm -v $$(pwd):/jmeter -w /jmeter $(JMETER_IMAGE) \
		-n -t $(JMETER_FILE) \
		-JHOSTNAME=host.docker.internal \
		-JPORT=8080 \
		-JTHREADS=1 \
		-JDURATION=$(TEST_DURATION) \
		-l $(RESULTS_DIR)/step4/result.jtl \
		-e -o $(RESULTS_DIR)/step4/report;

	# Восстановление всех сервисов
	echo "Восстановление всех остановленных сервисов...";
	docker compose up -d postgres_slave1 dialogs-1 socnet-1;

	@echo "Нагрузочные тесты завершены. HTML-отчеты сохранены в $(RESULTS_DIR)."
