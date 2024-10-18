#!/bin/bash

echo pg_current_wal_lsn: postgres            $(docker compose exec postgres psql -U app_hw -d hw -t -c "SELECT pg_current_wal_lsn();" | tr -d '[:space:]')
echo pg_current_wal_lsn: postgres_slave1     $(docker compose exec postgres_slave1 psql -U app_hw -d hw -t -c "SELECT pg_current_wal_lsn();" | tr -d '[:space:]')
echo pg_current_wal_lsn: postgres_slave2     $(docker compose exec postgres_slave2 psql -U app_hw -d hw -t -c "SELECT pg_current_wal_lsn();" | tr -d '[:space:]')
echo pg_last_wal_replay_lsn: postgres        $(docker compose exec postgres psql -U app_hw -d hw -t -c "SELECT pg_last_wal_replay_lsn();" | tr -d '[:space:]')
echo pg_last_wal_replay_lsn: postgres_slave1 $(docker compose exec postgres_slave1 psql -U app_hw -d hw -t -c "SELECT pg_last_wal_replay_lsn();" | tr -d '[:space:]')
echo pg_last_wal_replay_lsn: postgres_slave2 $(docker compose exec postgres_slave2 psql -U app_hw -d hw -t -c "SELECT pg_last_wal_replay_lsn();" | tr -d '[:space:]')

# Получаем позиции WAL для postgres_slave1
wal_postgres_slave1=$(docker compose exec postgres_slave1 psql -U app_hw -d hw -t -c "SELECT pg_last_wal_replay_lsn();" | tr -d '[:space:]')

# Получаем позиции WAL для postgres_slave2
wal_postgres_slave2=$(docker compose exec postgres_slave2 psql -U app_hw -d hw -t -c "SELECT pg_last_wal_replay_lsn();" | tr -d '[:space:]')

echo "WAL position on postgres_slave1: $wal_postgres_slave1"
echo "WAL position on postgres_slave2: $wal_postgres_slave2"

# Сравниваем позиции WAL
if [[ "$wal_postgres_slave1" > "$wal_postgres_slave2" ]]; then
    echo "postgres_slave1 is the freshest."
elif [[ "$wal_postgres_slave2" > "$wal_postgres_slave1" ]]; then
    echo "postgres_slave2 is the freshest."
else
    echo "postgres_slave1 and postgres_slave2 are same"
fi
