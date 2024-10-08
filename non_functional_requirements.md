# Нефункциональные требования
## ДЗ 1

- Любой язык программирования
- В качестве базы данных использовать PostgreSQL (при желании и необходимости любую другую SQL БД)
- Не использовать ORM
- Программа должна представлять из себя монолитное приложение.
- Не рекомендуется использовать следующие технологии:
	1. Репликация
	2. Шардирование
	3. Индексы
	4. Кэширование

Для удобства разработки и проверки задания можно воспользоваться [этой спецификацией](https://github.com/OtusTeam/highload/blob/master/homework/openapi.json "этой спецификацией") и реализовать в ней методы:
- `/login`
- `/user/register`
- `/user/get/{id}`  
      
    _Фронт опционален._  
      
    Сделать инструкцию по локальному запуску приложения, приложить Postman-коллекцию.
## ДЗ 2
- реализовать функционал автоматической загрузки в БД 1 млн анкет для задач нагрузочного тестирования
- 
# Текущий статус


| №     | Описание требования                                                                                                                                                                                                                                                      | Реализация                                                                                                                                                                      |
| ----- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| НФТ-1 | Любой язык программирования                                                                                                                                                                                                                                              | Go 1.23                                                                                                                                                                         |
| НФТ-2 | В качестве базы данных использовать PostgreSQL (при желании и необходимости любую другую SQL БД)                                                                                                                                                                         | PostgreSQL 16.4                                                                                                                                                                 |
| НФТ-3 | Не использовать ORM                                                                                                                                                                                                                                                      | выполнено                                                                                                                                                                       |
| НФТ-4 | Программа должна представлять из себя монолитное приложение.                                                                                                                                                                                                             | выполнено                                                                                                                                                                       |
| НФТ-5 | Не рекомендуется использовать следующие технологии:<br>1. Репликация<br>2. Шардирование<br>3. Индексы<br>4. Кэширование                                                                                                                                                  | Индексы использованы на минимально необходимом уровне (для обеспечения PK и уникальности, где это необходимо)                                                                   |
| НФТ-6 | Для удобства разработки и проверки задания можно воспользоваться [этой спецификацией](https://github.com/OtusTeam/highload/blob/master/homework/openapi.json "этой спецификацией") и реализовать в ней методы:<br>- `/login`<br>- `/user/register`<br>- `/user/get/{id}` | указанная спецификация противоречит ФТ (отсутствует поле Пол, отсутствует поле Фамилия, лишнее поле Отчество). Воспользовался, но затем доработал контракты.                    |
| НФТ-7 | _Фронт опционален._                                                                                                                                                                                                                                                      | фронт запроектирован, но пока не закончен                                                                                                                                       |
| НФТ-8 | Сделать инструкцию по локальному запуску приложения, приложить Postman-коллекцию.                                                                                                                                                                                        | [инструкция по локальному запуску](../README)<br>[Postman коллекция](https://github.com/Vasiliy82/otus-hla-homework/blob/main/misc/OTUS%20Homework.postman_collection.json)<br> |
| НФТ-9 | Инструментарий для загрузки тестовых данных в количестве 1 млн записей (таблица users)                                                                                                                                                                                   |                                                                                                                                                                                 |

