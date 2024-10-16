# otus-hla-homework
## Домашняя работа по курсу [OTUS](https://otus.ru/) [Highload Architect](https://otus.ru/lessons/highloadarchitect/)
### 1 Как работать с проектом
#### 1.1 Требования к окружению
Разрабатывалось и тестировалось в окружении: Ubuntu 22.04 x86_64, Docker 27.2.0, GNU Make 4.3
#### 1.2 запуск
```sh
git clone https://github.com/Vasiliy82/otus-hla-homework.git
make up
```
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
##### 2.3.1 [Нагрузочное тестирование](docs/hw-2-report.md)
##### 2.3.2 [Репликация](hw-3-report.md)

### 3 Описание компонентов
#### 3.1 [backend](./backend/README.md)
#### 3.2 [frontend](./frontend/README.md)
#### 3.3 [database](./postgresql.md)
### 4. Использованные инструменты (в разработке)