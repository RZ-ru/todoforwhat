# TODO Service

REST API сервис для управления задачами, написанный на Go с использованием:

* Gin
* PostgreSQL
* RabbitMQ
* Outbox Pattern
* Docker / Docker Compose

Проект демонстрирует построение backend-приложения с событийной архитектурой и надежной доставкой событий через Outbox Pattern.

---

# Возможности

* CRUD для задач
* REST API
* PostgreSQL как основная БД
* RabbitMQ для обработки событий
* Outbox Pattern
* Consumer для обработки сообщений
* Docker Compose окружение
* Layered architecture

---

# Стек технологий

| Технология     | Назначение               |
| -------------- | ------------------------ |
| Go             | Backend                  |
| Gin            | HTTP framework           |
| PostgreSQL     | Хранение данных          |
| RabbitMQ       | Очереди сообщений        |
| Docker         | Контейнеризация          |
| Docker Compose | Локальная инфраструктура |

---

# Архитектура проекта

```text
TODO/
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── gin/
│   ├── main.go
│   ├── migrations/
│   └── internal/
│       ├── consumer/
│       ├── handlers/
│       ├── models/
│       ├── outbox/
│       ├── repository/
│       └── services/
```

## Слои приложения

### handlers

HTTP слой.

Отвечает за:

* обработку запросов;
* валидацию данных;
* возврат HTTP ответов.

### services

Бизнес-логика приложения.

### repository

Работа с PostgreSQL.

### outbox

Реализация Outbox Pattern:

* сохранение событий;
* публикация в RabbitMQ;
* retry-механизм.

### consumer

Подписчик RabbitMQ.

---

# Outbox Pattern

Проект использует Outbox Pattern для надежной доставки событий.

## Как это работает

1. Данные сохраняются в PostgreSQL.
2. В той же транзакции создается запись в таблице `outbox`.
3. Worker читает необработанные события.
4. Событие публикуется в RabbitMQ.
5. После успешной отправки событие помечается как обработанное.

Это позволяет избежать потери событий при сбоях.

---

# API Endpoints

## Healthcheck

### GET /ping

```bash
curl http://localhost:8080/ping
```

Response:

```json
{
  "message": "pong"
}
```

---

## Создать задачу

### POST /tasks

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Learn Go"
  }'
```

---

## Получить список задач

### GET /tasks

```bash
curl http://localhost:8080/tasks
```

---

## Получить задачу по ID

### GET /tasks/:id

```bash
curl http://localhost:8080/tasks/1
```

---

## Обновить задачу

### PATCH /tasks/:id

```bash
curl -X PATCH http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Learn Go deeply"
  }'
```

---

## Удалить задачу

### DELETE /tasks/:id

```bash
curl -X DELETE http://localhost:8080/tasks/1
```

---

# Запуск проекта

## Требования

* Docker
* Docker Compose

---

## Клонирование

```bash
git clone https://github.com/RZ-ru/todoforwhat.git
cd todoforwhat
```

---

## Запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен:

```text
http://localhost:8080
```

---

# Переменные окружения

| Переменная  | Значение |
| ----------- | -------- |
| DB_HOST     | postgres |
| DB_PORT     | 5432     |
| DB_USER     | postgres |
| DB_PASSWORD | postgres |
| DB_NAME     | appdb    |

---

# RabbitMQ

Используется:

* Exchange: `tasks_exchange`
* Queue: `tasks`
* Exchange type: `fanout`

События:

* `task_created`
* `task_updated`
* `task_deleted`

---

# Пример события

```json
{
  "task_id": "1",
  "description": "Learn Go"
}
```

---

# Docker

## Dockerfile

Сервис собирается внутри контейнера:

```dockerfile
FROM golang:1.25
```

---

# Возможные улучшения

* JWT авторизация
* Swagger документация
* Unit tests
* Integration tests
* Kubernetes deployment
* CI/CD pipeline
* Redis cache
* gRPC communication
* Structured logging
* Metrics + Prometheus
* Graceful shutdown
* Config management

---

# Статус проекта

Pet-project / learning project.

Проект демонстрирует:

* чистую архитектуру;
* работу с PostgreSQL;
* RabbitMQ integration;
* Outbox Pattern;
* event-driven подход.

---

# Лицензия

MIT
