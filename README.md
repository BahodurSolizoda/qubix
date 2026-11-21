## Qubix — Quiz Platform

## Описание проекта
## Qubix — backend-приложение для создания и прохождения тестов.
## Пользователи могут регистрироваться, авторизоваться, просматривать тесты и проходить их. Авторы могут создавать и редактировать тесты, вопросы и варианты ответов.
## После завершения попытки пользователь получает результат прохождения теста.

## Используемые технологии
## Go — язык программирования.
## net/http — обработка HTTP-запросов.
## PostgreSQL — реляционная база данных.
## database/sql — работа с базой данных.
## lib/pq — драйвер PostgreSQL.
## bcrypt — хеширование паролей.
## JWT — аутентификация пользователей.

## Архитектура проекта:
Проект использует трёхслойную архитектуру:
## Handler — принимает HTTP-запросы и возвращает HTTP-ответы.
## Service — содержит бизнес-логику и проверяет данные.
## Repository — взаимодействует с базой данных PostgreSQL.

## Структура проекта:

qubix/
├── cmd/
│   └── api/
│       └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── database/
│   ├── handler/
│   ├── model/
│   ├── repository/
│   └── service/
├── middleware/
├── migrations/
├── seeds/
├── go.mod
├── go.sum
└── README.md


## Основные возможности:
## 1. Аутентификация и авторизация:
- Регистрация пользователей.
- Вход в систему.
- Хеширование паролей с помощью bcrypt.
- Авторизация через JWT.
- Защита приватных маршрутов.
- Разделение ролей user и author.

## 2. Управление тестами:
- Создание тестов.
- Получение списка тестов.
- Получение теста по ID.
- Редактирование и удаление тестов.
- Проверка прав автора на управление собственными тестами.

## 3. Управление вопросами и вариантами ответов:
- Создание вопросов.
- Получение вопросов.
- Редактирование и удаление вопросов.
- Создание вариантов ответов.
- Указание правильных вариантов.
- Защита правильных ответов от выдачи через публичный API.

## 4. Прохождение тестов:
- Создание попытки прохождения теста.
- Отправка ответов.
- Завершение попытки.
- Автоматическая проверка ответов.
- Подсчёт баллов.
- Получение результата.

##  5. Категории:
- Создание категорий.
- Получение списка категорий.
- Получение категории по ID.
- Редактирование и удаление категорий.

## 6. База данных:
- Проект использует PostgreSQL.
Основные таблицы:
- users
- categories
- quizzes
- questions
- options
- attempts
- answers
- answer_options
- results

##  7. Структура базы данных создаётся с помощью SQL-миграций, расположенных в папке migrations.


##  Настройка и запуск

##  1. Клонирование проекта:

git clone <repository-url>
cd qubix

##  2. Настройка переменных окружения
Перед запуском необходимо настроить:
PORT=8080
DATABASE_URL=postgres://postgres:password7070@localhost:5432/qubix?sslmode=disable
JWT_SECRET=your-secret-key

## Для получения secret key JWT: 
$env:JWT_SECRET

Укажи собственные данные PostgreSQL и секретный ключ JWT. Не публикуй реальные секреты.

## 3. Установка зависимостей:
go mod tidy

## 4. Запуск сервера:
go run ./cmd/api

По умолчанию сервер запускается на: http://localhost:8080


## Проверка проекта:
go test ./...




## Основные API-маршруты

| Метод | Маршрут                         | Назначение |

| POST | /sign-up                         | Регистрация |
| POST | /login                           | Авторизация |
| GET  | /private                         | Приватный маршрут |
| GET  | /categories                      | Список категорий |
| POST | /categories                      | Создание категории |
| GET  | /quizzes                         | Список тестов |
| POST | /quizzes                         | Создание теста |
| GET  | /quizzes/{id}                    | Получение теста |
| POST | /quizzes/{quiz_id}/questions     | Создание вопроса |
| GET  | /quizzes/{quiz_id}/questions     | Вопросы теста |
| POST | /questions/{question_id}/options | Создание варианта ответа |
| POST | /quizzes/{quiz_id}/attempts      | Начало попытки |
| POST | /attempts/{id}/answers           | Отправка ответа |
| POST | /attempts/{id}/finish            | Завершение попытки |
| GET  | /attempts/{id}/result            | Получение результата |


Основной функционал backend-приложения реализован.
Проверены основные сценарии регистрации, авторизации, управления тестами, прохождения тестов и получения результатов.




Swagger Editor

Qubix API can be tested and documented using Swagger Editor.
1. Start the Qubix server
--> Open PowerShell in the project directory and run:

go run ./cmd/api


--> The server starts at:
http://localhost:8080


--> Keep this terminal window open while testing the API.

2. Open Swagger Editor
--> Open the following website:

https://editor.swagger.io/

Import the project's `openapi.yaml` file into Swagger Editor.
You can also paste the contents of `openapi.yaml` into the editor.

3. Authorize using JWT
Some API endpoints require Bearer Token authentication.

First, log in through the `POST /login` endpoint or use PowerShell:

$body = @{
    login = "YOUR_LOGIN"
    password = "YOUR_PASSWORD"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -ContentType "application/json" -Body $body


--> Copy the JWT token to the clipboard:

Set-Clipboard -Value $response.token


In Swagger Editor:

1. Click `Authorize`.
2. Paste the JWT token into the `BearerAuth` field.
3. Do not add the word `Bearer` manually.
4. Click `Authorize`, then `Close`.

4. Test API endpoints

1. Select an endpoint, for example `GET /private`.
2. Click `Try it out`.
3. Click `Execute`.
4. Check the server response.

Expected successful response:


Code: 200
Access granted!


### 5. CORS configuration

Qubix includes CORS middleware that allows requests from:


https://editor.swagger.io


It allows the following HTTP methods:

- GET
- POST
- PUT
- DELETE
- OPTIONS

It also allows the `Authorization` and `Content-Type` headers required for API requests.

`Important notes`

- The Qubix server must be running before testing API endpoints.
- The API address is `http://localhost:8080`.
- Protected endpoints require a valid JWT token.
- Never share your JWT token, database password, or `JWT_SECRET`.
