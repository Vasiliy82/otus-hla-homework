-- Добавление воркеров в кластер
SELECT master_add_node('worker1', 5432);
SELECT master_add_node('worker2', 5432);
SELECT master_add_node('worker3', 5432);

-- Проверка активных нод
SELECT * FROM master_get_active_worker_nodes();