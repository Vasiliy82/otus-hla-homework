# otus-hla-homework
## Домашняя работа по курсу [OTUS](https://otus.ru/) [Highload Architect](https://otus.ru/lessons/highloadarchitect/)
### 1 Как работать с проектом
#### 1.1 Требования к окружению
Разрабатывалось и тестировалось в окружении: Ubuntu 22.04 x86_64, Docker 27.2.0, GNU Make 4.3
#### 1.2 запуск
```sh
git clone https://github.com/Vasiliy82/otus-hla-homework.git
cd otus-hla-homework/environments/hw03
make up
```
[frontend](http://localhost:5173)
[backend](http://localhost:8080)
[backend metrics](http://localhost:8080/metrics)
[grafana](http://localhost:3000)
[prometheus](http://localhost:9090)
[cadvisor metrics](http://localhost:8081/metrics)
[postgres-exporter metrics](http://localhost:9187/metrics)
[postgres_slave1-exporter metrics](http://localhost:9188/metrics)
[postgres_slave2-exporter metrics](http://localhost:9189/metrics)

#### 1.3. остановка
```sh
make down
```
#### 1.4 удаление
```sh
make destroy
```
#### 1.5 коллекция postan [тут](https://github.com/Vasiliy82/otus-hla-homework/blob/main/misc/OTUS%20Homework.postman_collection.json)

#### 1.6 отладка в VS Code
Пример .vscode/launch.json
```
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Attach to Backend in Docker",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "/app",  // путь к коду внутри контейнера
            "port": 2345,          // порт, на котором Delve работает в контейнере
            "host": "localhost",   // IP хоста для подключения (localhost, так как мы перенаправляем порт)
//            "localRoot": "${workspaceFolder}/backend",  // путь к коду на локальной машине
            "cwd": "${workspaceFolder}/backend",
//            "trace": "verbose"     // включить вывод для отладки (можно убрать, если не нужно)
          }
    ]
}
```
### 2 Документация
#### 2.1 [Функциональные требования](./functional_requirements.md)
#### 2.2 [Нефункциональные требования](./non_functional_requirements.md)
#### 2.3 Отчеты о проделанной работе
##### 2.3.1 [Нагрузочное тестирование](./docs/hw2/hw-2-report.md)
##### 2.3.2 [Репликация](./docs/hw3/hw-3-report.md)

### 3 Описание компонентов
#### 3.1 [backend](./backend/README.md)
#### 3.2 [frontend](./frontend/README.md)
#### 3.3 [database](./postgresql.md)
### 4. Использованные инструменты (в разработке)
### 5. Todo
    5.1 Прописать автоматическую настройку Прометеуса и Графаны в /misc/*/Dockerfile;
    5.2 Добавить ELK и перевести туда логирование;
    5.3 ~~Реализовать конфиг файл, вынести туда настройки пулинга подключений к БД~~
    5.4 Вынести метрики пулинга БД в Прометеус и Графану
    5.5 Во всех Dockerfile / docker-compose сервисах зафиксировать конкретные версии образов
    5.6 Echo.Logger wrapper for zap
    5.7 Прошерстить весь конфиг на предмет дефолтных значений. Кроме необходимых настроек (БД, и пр.), все должны иметь значение по умолчанию.
    5.8 Подумать о том, что тип данных AppError вообще не предназначен для технических ошибок. Он должен отражать ошибки бизнесовые, с текстом. Тогда, если возникает техническая ошибка, то она будет другого типа и в UI никогда не попадет.