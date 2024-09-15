#!/bin/sh
# в случае правок в этом файле контейнер backend придется пересобрать

# Выполняем замену переменных окружения в шаблоне и создаем /etc/app.yaml. Почему не в текущей папке? Потому что не хотим, чтобы скрипты 
# внутри контейнера редактировали файл в проекте
envsubst < /app/app.yaml.tpl > /etc/app.yaml

# Выполняем миграции
goose -dir /app/migrations postgres "host=$POSTGRES_HOST port=$POSTGRES_PORT user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=disable" up

# Запускаем air для работы с приложением
air