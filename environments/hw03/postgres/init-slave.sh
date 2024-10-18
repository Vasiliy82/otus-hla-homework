#!/bin/bash

# Путь к директории данных PostgreSQL
PGDATA="/var/lib/postgresql/data"

echo "===================================================="
echo "===================================================="
echo "                  init-slave.sh                     "
echo $(hostname)
echo "===================================================="
echo "===================================================="

echo "Data directory is empty, initializing replication..."

# Очистка директории данных (на всякий случай)
rm -rf $PGDATA/*

# Что важнее, удобство или безопасность? :)
export PGPASSWORD='replica_password'

# Проверяем, существует ли слот репликации
if psql -h postgres -U replicator -d hw -c "SELECT 1 FROM pg_replication_slots WHERE slot_name='$(hostname)';" | grep -q 1; then
  echo "Slot '$(hostname)' already exists"
  psql -h postgres -U replicator -d hw -c "SELECT pg_drop_replication_slot('$(hostname)');"
fi

# Выполнение pg_basebackup для инициализации реплики
pg_basebackup -h postgres -D $PGDATA -U replicator -Fp -Xs -P -R -C -S "$(hostname)" -v
  
# Конечно, безопасность!
unset PGPASSWORD

