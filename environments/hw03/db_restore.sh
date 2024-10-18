#!/bin/bash

cat ../../misc/backup/hw_backup.gz | docker compose exec -T postgres bash -c "gunzip | pg_restore -U app_hw -d hw -v"

