-- Создание пользователя
CREATE ROLE app_hw_messages WITH LOGIN PASSWORD 'Passw0rd';

-- Создание базы данных
CREATE DATABASE hw_messages OWNER app_hw_messages;

-- Назначение прав доступа
GRANT ALL PRIVILEGES ON DATABASE hw_messages TO app_hw_messages;