DC=docker compose
PG_HOST=localhost
PG_PORT=5432
PG_SUPERUSER=postgres
PG_SUPERPASSWORD=postgres
APP_DB=hw_messages

MIGRATION_DB_URL=postgres://$(PG_SUPERUSER):$(PG_SUPERPASSWORD)@master:$(PG_PORT)/$(APP_DB)?sslmode=disable

export PG_SUPERUSER PG_SUPERPASSWORD APP_DB

# Команды
up:
	$(DC) up -d
	sleep 10 # Ждем, чтобы сервисы успели подняться
	make init-cluster
	make migrate-up

down:
	$(DC) down

cleanup:
	$(DC) down --remove-orphans

destroy:
	$(DC) down -v

migrate-up:
	$(DC) run --rm goose goose -dir /app/migrations postgres "$(MIGRATION_DB_URL)" up

migrate-down:
	$(DC) run --rm goose goose -dir /app/migrations postgres "$(MIGRATION_DB_URL)" down

init-cluster:
	$(DC) exec master psql -h $(PG_HOST) -U $(PG_SUPERUSER) -d ${APP_DB} -c "\
	SELECT master_add_node('worker1', $(PG_PORT)); \
	SELECT master_add_node('worker2', $(PG_PORT)); \
	SELECT master_add_node('worker3', $(PG_PORT)); \
	"

