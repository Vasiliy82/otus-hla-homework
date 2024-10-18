-- Создание роли для репликации
CREATE ROLE replicator WITH REPLICATION LOGIN PASSWORD 'replica_password';
CREATE DATABASE hw;
