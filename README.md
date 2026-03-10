# Тестовое задание на позицию middle Golang.

## Компания ООО "МКК ЛУНА"

## Описание проекта

Это REST API сервис, реализованный на **Go**, предназначенный для работы с задачами, командами и статистикой.
Проект построен с использованием чистой архитектуры и включает в себя:

- HTTP API
- авторизацию пользователей
- работу с задачами и командами
- сбор метрик для мониторинга
- rate-limiting для защиты API
- контейнеризацию через Docker

API реализовано с использованием роутера **Chi** и поддерживает документацию **Swagger**.

---

# Основные возможности

### HTTP API

Реализованы эндпоинты для работы с:

- авторизацией
- задачами
- командами
- статистикой

Все API ручки имеют префикс:

```
/api/v1/
```

---

### Swagger документация

Документация API доступна через Swagger UI:

```
/swagger/index.html
```

Позволяет:

- просматривать все эндпоинты
- тестировать запросы
- смотреть схемы ответов

---

### Prometheus метрики

В проекте реализован сбор метрик для мониторинга API.

Доступны следующие метрики:

#### HTTP метрики

- `http_requests_total` — общее количество HTTP запросов
- `http_request_duration_seconds` — время выполнения HTTP запросов
- `http_requests_errors_total` — количество ошибок API

#### Метки метрик

```
method
path
status
```

Метрики доступны по адресу:

```
/metrics
```

---

### Rate Limiting

Реализован middleware для ограничения количества запросов от пользователя.

Ограничение:

```
100 запросов в минуту на пользователя
```

Идентификация пользователя происходит через `UserID` из сессии.

---

### Мониторинг Go runtime

Prometheus автоматически собирает runtime метрики:

- количество goroutine
- использование памяти
- время работы GC
- CPU использование

---

# Технологии

Проект использует:

- Go
- Chi router
- Redis
- MariaDB
- Prometheus
- Docker
- Docker Compose
- Swagger

---

# Структура проекта

Основные директории проекта:

```
cmd/api            - точка входа приложения
internal/
  adapters         - внешние адаптеры (http, redis, db)
  controllers      - HTTP контроллеры
  middlewares      - middleware (rate limiter, metrics)
  services         - бизнес логика
  config           - конфигурация приложения
```

---

# Требования для запуска

Для запуска проекта необходимо установить:

- Go 1.26+
- Docker
- Docker Compose

---

# Переменные окружения

Приложение использует `.env` файл для конфигурации.

Пример `.release.env`:

```
APP_ENV=release

HTTP_SERVER_PORT=5555

MYSQL_ROOT_PASSWORD=password
MYSQL_DATABASE=app
MYSQL_USER=user123
MYSQL_PASSWORD=password123

MYSQL_HOST=mariadb
MYSQL_PORT=3306

MYSQL_DSN=${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}?parseTime=true&loc=UTC


REDIS_PASSWORD=password123
REDIS_HOST=redis
REDIS_PORT=6379

REDIS_ADDR=${REDIS_HOST}:${REDIS_PORT}

MIGRATIONS_PATH=migrations
MIGRATIONS_TABLE=migrations

ACCESS_TOKEN_SECRET=secretkeyforaccestoken
ACCESS_TOKEN_DURATION=10m

REFRESH_TOKEN_DURATION=1h

TOKEN_ISSUER=application_name

PROMETHEUS_PORT=9090
```

---

# Запуск проекта через Docker

Сборка и запуск:

```
docker compose --env-file .release.env up --build
```

После запуска будут доступны:

| Сервис             | Адрес                                    |
| ------------------ | ---------------------------------------- |
| API                | http://localhost:5555                    |
| Swagger            | http://localhost:5555/swagger/index.html |
| Prometheus metrics | http://localhost:8080/metrics            |
| Prometheus UI      | http://localhost:9090                    |

---

# Пример API запросов

Получить список задач:

```
GET /api/v1/tasks
```

Получить список команд:

```
GET /api/v1/teams
```

Получить статистику команд:

```
GET /api/v1/statistics/teams
```

Авторизация:

```
POST /api/v1/login
```

# Завершение работы

При остановке приложения корректно завершаются:

- фоновые goroutine
- отключение от бд и редис
- остановка http сервера
- очистка rate limiter

---

# Итог

В проекте реализовано:

- REST API
- Swagger документация
- Prometheus метрики
- Rate limiting
- Docker контейнеризация
- Redis кэш
- MariaDB база данных

Проект готов к локальному запуску и дальнейшему расширению.
