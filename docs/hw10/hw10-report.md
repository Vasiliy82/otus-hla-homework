## Домашнее задание

Сервис счетчиков
Цель:
В результате выполнения ДЗ вы создадите сервис счетчиков. Сервис будет хранить такие счетчики, как число непрочитанных сообщений.
В данном задании тренируются навыки:
- разработка отказоустойчивых сервисов;
- использование кешей.

Описание/Пошаговая инструкция выполнения домашнего задания:
**Реализовать функционал:**  
**Вариант 1**
1. Разработайте сервис счетчиков.
2. Учтите то, что на этот сервис будет большая нагрузка, особенно на чтение.
3. Продумайте, как обеспечить консистентность между счетчиком и реальным числом непрочитанных сообщений. Например, используйте паттерн SAGA.
4. Внедрите сервис для отображения счетчиков.
**Вариант 2**  
Разработайте и внедрите функционал, в котором будет необходимость применения паттерна SAGA, этот функционал должны быть потенциально нагруженным.

---
**Требования**
- обеспечение консистентности (SAGA в приоритете)
---
 
**Форма сдачи ДЗ**
- Предоставить ссылку на исходный код (github, gitlab, etc)
- Предоставить докеризированное приложение, которое можно запустить при помощи docker-compose (может лежать рядом с исходным кодом) ИЛИ развернутое приложение, доступное извне ИЛИ инструкция по запуску
- Предоставить отчет об архитектуре (схема, словесное описание)

Критерии оценки:
Оценка происходит по принципу зачет/незачет.
Требования:
1. Верно описан выбранный паттерн обеспечения консистентности.
2. Выбранная архитектура сервиса подходит для решения задачи.

  Компетенции:
- Разработка и проектирование микросервисов
    - - знать принципы организации микросервисной архитектуры
    - - знать тактики достижения Consistency в MSA

## Отчет о проделанной работе

### 1. Описание паттерна SAGA и его использование

SAGA представляет собой распределенную транзакцию, состоящую из последовательности локальных транзакций.  
Каждая локальная транзакция имеет операцию отката (компенсацию), если последующие шаги завершаются неудачно.

В проекте SAGA реализована следующим образом:

1. `DialogService.SendMessage` сохраняет сообщение и отправляет событие `SagaMessageSent` в Kafka.
2. Kafka Consumer в сервисе счетчиков получает `SagaMessageSent` и увеличивает счетчик.
3. Если увеличение счетчика успешно, отправляется `SagaCounterIncremented`, иначе `SagaCounterIncrementFailed`.
4. Если получено `SagaCounterIncrementFailed`, `SagaOrchestrator` помечает сообщение как `failed`.

Удаление данных не является обязательной частью отката транзакции. В данной реализации сообщения не удаляются, а помечаются как `failed`. Это позволяет сохранять историю событий, но при этом исключать такие сообщения из активной бизнес-логики.

### 2. Основные доработки

#### 2.1. Добавление объектов SAGA

`SagaOrchestrator` управляет процессами SAGA, отправляет и обрабатывает события.
[`backend/internal/usecases/saga_orchestrator.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/backend/internal/usecases/saga_orchestrator.go#L1-L57)

`ProcessSagaEventUseCase` обрабатывает `SagaMessageSent` и увеличивает счетчик.
[`counters/internal/usecases/process_saga_event.go`]()

`SagaEvent` - структура для событий SAGA.
[`counters/internal/domain/models.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/counters/internal/domain/models.go#L11-L17)

`KafkaConsumer` слушает топик SAGA и передает события в `ProcessSagaEventUseCase`.
[`counters/internal/infrastructure/kafka/consumer.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/counters/internal/infrastructure/kafka/consumer.go#L1-L58)

`KafkaProducer` отправляет события `SagaMessageSent`, `SagaCounterIncremented`, `SagaCounterIncrementFailed`.
[`counters/internal/infrastructure/kafka/producer.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/counters/internal/infrastructure/kafka/producer.go#L1-L42)

#### 2.2. Доработка `DialogService`

`SendMessage` теперь возвращает `TransactionID`, сохраняет сообщение как `pending` и отправляет `SagaMessageSent`.  
[`backend/internal/services/dialogs.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/backend/internal/services/dialogs.go)

`SaveMessageWithSaga` сохраняет сообщения с `transaction_id` и `saga_status`.  
[`backend/internal/repository/dialog.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/backend/internal/repository/dialog.go)

`UpdateSagaStatus` обновляет `saga_status` в БД (`pending → committed/failed`).  
[`backend/internal/repository/dialog.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/backend/internal/repository/dialog.go)

#### 2.3. Поддержка SAGA в `counters` (сервис счетчиков)

`CounterService` - интерфейс управления счетчиками.  
[`counters/internal/domain/services.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/counters/internal/domain/services.go)

`CombinedCounterRepository` работает с PostgreSQL и Redis.  
[`counters/internal/infrastructure/repository/combined_repo.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/counters/internal/infrastructure/repository/combined_repo.go)

`PGCounterRepository` - работа с БД.  
[`counters/internal/infrastructure/repository/pg_repo.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/counters/internal/infrastructure/repository/pg_repo.go)

`RedisCounterRepository` - кеширование в Redis.  
[`counters/internal/infrastructure/repository/redis_repo.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/counters/internal/infrastructure/repository/redis_repo.go)

`CounterService` - инкремент и сброс счетчика с учетом SAGA.  
[`counters/internal/usecases/counter_service.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/counters/internal/usecases/counter_service.go)

#### 2.4. Интеграция с Kafka

`SendSagaEvent` - отправка сообщений SAGA в Kafka.  
[`backend/internal/infrastructure/broker/producer.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/backend/internal/infrastructure/broker/producer.go)

`SagaBusProcessor` - читает сообщения SAGA из Kafka и вызывает `HandleSagaEvent`.  
[`backend/internal/services/saga_bus_processor.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/backend/internal/services/saga_bus_processor.go)

`KafkaConsumer` читает SAGA-события в сервисе счетчиков.  
[`counters/internal/infrastructure/kafka/consumer.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/counters/internal/infrastructure/kafka/consumer.go)

`KafkaProducer` публикует SAGA-события (`SagaCounterIncremented`, `SagaCounterIncrementFailed`).  
[`counters/internal/infrastructure/kafka/producer.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/counters/internal/infrastructure/kafka/producer.go)

#### 2.5. Обновление контекста и `X-Request-ID`

`ExtractRequestIDFromContext`, `AddRequestIDToContext`, `ExtractRequestIDFromKafka`, `AddRequestIDToKafka`.  
[`common/utils/request_id.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/common/utils/request_id.go)

`RequestIDMiddleware` теперь использует `utils/request_id.go`.  
[`common/infrastructure/http/middleware/request_id.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/common/infrastructure/http/middleware/request_id.go)

Добавлено извлечение `X-Request-ID` из Kafka и передача в контекст.  
[`backend/internal/infrastructure/broker/workerpool.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/backend/internal/infrastructure/broker/workerpool.go)

#### 2.6. Подключение компонентов в `app.go`

Добавлен `SagaOrchestrator`, `SagaBusProcessor`, `KafkaProducer`.  
[`backend/cmd/dialogs/main.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/backend/cmd/dialogs/main.go)

Добавлен запуск Kafka Consumer для обработки `SagaMessageSent`.  
[`counters/cmd/main.go`](https://github.com/Vasiliy82/otus-hla-homework/blob/4af71dbe365cbfe337fe08eba327a12fbc004c92/counters/cmd/main.go)

### 3. Порядок запуска проекта и результаты тестирования

```bash
git clone https://github.com/Vasiliy82/otus-hla-homework.git
cd otus-hla-homework
git checkout tags/hw10
cd environments/hw10
make up

# Отправим несколько сообщений
curl --location 'http://localhost:8081/api/dialog/b3dd0beb-26a1-4437-a7bd-3a57e40bb4e3/send' \
--header 'X-User-Id: d86c3ba3-0628-4716-88b8-a51ec02bc603' \
--header 'Content-Type: application/json' \
--header 'X-Request-ID: a770e39b-e947-42a3-ada1-a441bb537b51' \
--data '{"text": "Test message from User1 to User2"}'

curl --location 'http://localhost:8081/api/dialog/d86c3ba3-0628-4716-88b8-a51ec02bc603/send' \
--header 'X-User-Id: b3dd0beb-26a1-4437-a7bd-3a57e40bb4e3' \
--header 'Content-Type: application/json' \
--header 'X-Request-ID: 436de90f-e42f-42c3-9710-d5a5ed1ed267' \
--data '{"text": "Test response from User2 to User1"}'

curl --location 'http://localhost:8081/api/dialog/b3dd0beb-26a1-4437-a7bd-3a57e40bb4e3/send' \
--header 'X-User-Id: d86c3ba3-0628-4716-88b8-a51ec02bc603' \
--header 'Content-Type: application/json' \
--header 'X-Request-ID: 8de1c260-cd85-42bd-95b5-13e0ba3d6580' \
--data '{"text": "Second test message from User1 to User2"}'

curl --location 'http://localhost:8081/api/dialog/d86c3ba3-0628-4716-88b8-a51ec02bc603/send' \
--header 'X-User-Id: b3dd0beb-26a1-4437-a7bd-3a57e40bb4e3' \
--header 'Content-Type: application/json' \
--header 'X-Request-ID: 57dfd8e5-9ae2-413b-bac9-5d0ace7875b4' \
--data '{"text": "Second test response from User2 to User1"}'

# Чтобы вызвать ошибку сохранения счетчика, добавим CHECK CONSTRAINT
docker compose exec postgres-counters psql -U app_hw -d hw -t -c "ALTER TABLE dialog_counters ADD CONSTRAINT unread_count_check CHECK (unread_count <= 2);"

# Отправим еще одно сообщение, которому не суждено быть доставленным
curl --location 'http://localhost:8081/api/dialog/b3dd0beb-26a1-4437-a7bd-3a57e40bb4e3/send' \
--header 'X-User-Id: d86c3ba3-0628-4716-88b8-a51ec02bc603' \
--header 'Content-Type: application/json' \
--header 'X-Request-ID: 15af67c2-1e63-4ae0-8b9e-928560559117' \
--data '{"text": "This message will fail"}'

# Убедимся, что сообщения в Kafka присутствуют
kafka-console-consumer --bootstrap-server localhost:9092 --topic mysocnet-saga-bus --from-beginning --property print.headers=true

# Убедимся, что сообщения в БД лежат, значения счетчиков актуальны

docker compose exec postgres psql -U app_hw -d hw -t --pset pager=off -c "SELECT * FROM messages;"
docker compose exec postgres-counters psql -U app_hw -d hw -t --pset pager=off -c "SELECT * FROM dialog_counters;"

# Вернем все обратно
docker compose exec postgres-counters psql -U app_hw -d hw -t -c "ALTER TABLE dialog_counters DROP CONSTRAINT unread_count_check;"

# Т.к. с прошлого задания  у нас сохранился RequestID, можно детально посмотреть обработку каждой операции в логах
# Успешно доставленные сообщения
docker compose logs dialogs | grep a770e39b-e947-42a3-ada1-a441bb537b51
docker compose logs counters | grep a770e39b-e947-42a3-ada1-a441bb537b51
docker compose logs dialogs | grep 436de90f-e42f-42c3-9710-d5a5ed1ed267
docker compose logs counters | grep 436de90f-e42f-42c3-9710-d5a5ed1ed267
docker compose logs dialogs | grep 8de1c260-cd85-42bd-95b5-13e0ba3d6580
docker compose logs counters | grep 8de1c260-cd85-42bd-95b5-13e0ba3d6580
docker compose logs dialogs | grep 57dfd8e5-9ae2-413b-bac9-5d0ace7875b4
docker compose logs counters | grep 57dfd8e5-9ae2-413b-bac9-5d0ace7875b4

# Сообщение, во время обработки которого произошла ошибка
docker compose logs dialogs | grep 15af67c2-1e63-4ae0-8b9e-928560559117
docker compose logs counters | grep 15af67c2-1e63-4ae0-8b9e-928560559117
```

Результаты выполнения команд в [report.txt](./report.txt)

## 4. Итоги

4. Реализована SAGA для синхронизации сообщений и счетчиков.
5. Добавлены механизмы отката транзакций (Rollback) в случае ошибок.
6. Используется Kafka для передачи событий между сервисами.
7. Добавлена поддержка `X-Request-ID` для трассировки запросов.
8. Обновлен `DialogService`, теперь он хранит сообщения с `transaction_id` и `saga_status`.
9. Добавлен `SagaBusProcessor`, который обрабатывает события в Kafka.
10. Обновлены `WorkerPool` и `RequestIDMiddleware` для поддержки `X-Request-ID` в Kafka.

Сервис сообщений и сервис счетчиков работают в распределенной транзакции (SAGA) с полной поддержкой обработки ошибок.

