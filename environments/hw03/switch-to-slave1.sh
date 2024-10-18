#!/bin/bash

# убиваем мастер, если он вдруг еще жив
docker compose kill postgres

# тут же промоутим postgres_slave1 до мастера
docker compose exec postgres_slave1 psql -U app_hw -d hw -c "SELECT pg_promote();"
docker compose exec postgres_slave1 psql -U app_hw -d hw -c "ALTER SYSTEM SET synchronous_standby_names = 'ANY 1 (postgres_slave2)';"
docker compose exec postgres_slave1 psql -U app_hw -d hw -c "SELECT pg_reload_conf();"
docker compose exec postgres_slave1 psql -U app_hw -d hw -c "SELECT * FROM pg_create_physical_replication_slot('postgres_slave2');"
docker compose exec postgres_slave1 psql -U app_hw -d hw -c "SELECT pg_reload_conf();"

# а postgres_slave1 переключаем на первый
docker compose exec postgres_slave2 psql -U app_hw -d hw -c "ALTER SYSTEM SET primary_slot_name = 'postgres_slave2';"
docker compose exec postgres_slave2 psql -U app_hw -d hw -c "ALTER SYSTEM SET primary_conninfo = 'user=replicator password=replica_password host=postgres_slave1 port=5432 application_name=postgres_slave2';"
docker compose exec postgres_slave2 psql -U app_hw -d hw -c "SELECT pg_reload_conf();"

# на первой ноде попробуем добавить запись
docker compose exec postgres_slave1 psql -U app_hw -d hw -c "INSERT INTO test DEFAULT VALUES;"

# а на второй тут же ее почитать
docker compose exec postgres_slave2 psql -U app_hw -d hw -c "SELECT MAX(id) FROM test;"

# Попробуем теперь вылечить бывший master и сделать его дополнительной репликой
# Добавляем файл standby.signal для старого мастера и настраиваем его на работу в роли слейва
docker compose run postgres bash -c "touch /var/lib/postgresql/data/standby.signal"

# попробуем поднять убитый мастер в качестве слейва
docker compose up -d postgres


# Ожидаем статуса healthy
echo "Waiting for postgres to become healthy..."
while [[ $(docker inspect --format='{{.State.Health.Status}}' postgres) != "healthy" ]]; do
  echo "Postgres is not healthy yet..."
  sleep 1
done

echo "Postgres is healthy!"

docker compose exec postgres psql -U app_hw -d hw -c "ALTER SYSTEM SET primary_slot_name = 'postgres';"
docker compose exec postgres psql -U app_hw -d hw -c "ALTER SYSTEM SET primary_conninfo = 'user=replicator password=replica_password host=postgres_slave1 port=5432 application_name=postgres';"
docker compose exec postgres psql -U app_hw -d hw -c "SELECT pg_reload_conf();"