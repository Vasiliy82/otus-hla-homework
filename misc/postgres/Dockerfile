# Используем конкретную версию PostgreSQL 16.4
FROM postgres:16.4

# Копируем ваш конфигурационный файл в контейнер
COPY postgresql.conf /etc/postgresql/postgresql.conf

# Указываем PostgreSQL использовать наш конфиг при запуске
CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]