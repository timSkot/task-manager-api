# Task Manager API

Простой REST API для управления задачами (todo-list), написанный на Go.

Проект построен по слоистой архитектуре (`handler → service → repository`), что позволяет менять хранилище данных, не затрагивая бизнес-логику и HTTP-обработчики.

## Стек

- Go 1.26+ (стандартная библиотека `net/http`, без сторонних роутеров)
- PostgreSQL 16
- `database/sql` + `lib/pq` (драйвер)
- `golang-migrate` — миграции БД
- `godotenv` — конфигурация через переменные окружения
- Docker + Docker Compose
- Prometheus — метрики и мониторинг

## Архитектура

```
main.go              — точка входа, сборка (wiring) слоёв, подключение к БД
internal/
  model/              — структуры данных (Task, UpdateTaskInput)
  repository/         — слой хранения данных
    task_repository.go          — интерфейс TaskRepository + in-memory реализация
    postgres_task_repository.go — реализация на PostgreSQL
  service/            — бизнес-логика, валидация, обработка ошибок
  handler/            — HTTP-обработчики (парсинг запросов, статус-коды)
    middleware.go      — логирование запросов и запись метрик
    metrics.go          — определения метрик Prometheus
migrations/           — SQL-миграции (up/down)
Dockerfile             — multi-stage сборка приложения
docker-compose.yml     — приложение + БД + Prometheus + автоприменение миграций
prometheus.yml         — конфигурация сбора метрик
```

Каждый слой знает только о слое ниже себя:

- **handler** знает про **service**
- **service** знает про **repository** — но только через интерфейс `TaskRepository`, не про конкретную реализацию
- **repository** ничего не знает о вышестоящих слоях

Благодаря этому переключение хранилища с in-memory на PostgreSQL потребовало изменения всего одной строки в `main.go` — `service` и `handler` остались нетронутыми.

### Обработка ошибок между слоями

`repository` возвращает собственную ошибку `ErrNotFound`, если запись не найдена, либо любую другую ошибку в случае сбоя (например, недоступность БД). `service` транслирует `repository.ErrNotFound` в свою `service.ErrTaskNotFound`, не пропуская наверх детали реализации хранилища. `handler` проверяет конкретную ошибку через `errors.Is` и выбирает подходящий HTTP-статус (`404` для "не найдено", `500` для непредвиденных сбоев).

### Логирование и метрики

Каждый запрос проходит через `LoggingMiddleware`, который оборачивает `http.ResponseWriter`, чтобы перехватить итоговый статус-код (паттерн decorator), и:

- логирует метод, путь, статус и длительность запроса;
- обновляет метрики Prometheus — счётчик запросов (`http_requests_total`, с разбивкой по методу/пути/статусу) и гистограмму длительности (`http_request_duration_seconds`).

Метрики доступны на `GET /metrics` в формате, который умеет читать Prometheus.

## Запуск (Docker, рекомендуемый способ)

Единственная команда для полного поднятия проекта — приложение, база данных, Prometheus и применение миграций происходят автоматически:

```bash
docker-compose up --build
```

Что при этом происходит:

1. Поднимается контейнер с PostgreSQL, Docker Compose дожидается его реальной готовности (healthcheck)
2. Отдельный контейнер применяет все накопленные миграции из папки `migrations/`
3. Собирается и запускается контейнер с приложением — только после успешного завершения миграций
4. Поднимается контейнер Prometheus, который начинает опрашивать `/metrics` приложения каждые 5 секунд

Доступные адреса:

- API: `http://localhost:8080`
- Метрики (сырые, текстовый формат): `http://localhost:8080/metrics`
- Prometheus (веб-интерфейс, графики, PromQL): `http://localhost:9090`

Остановить и убрать контейнеры:

```bash
docker-compose down
```

Остановить и полностью очистить данные БД (например, чтобы проверить запуск "с нуля"):

```bash
docker-compose down -v
```

## Запуск локально (без Docker)

Если нужно запускать приложение напрямую на хосте (например, для отладки), а базу — в контейнере:

### 1. Поднять PostgreSQL

```bash
docker run --name task-manager-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=taskmanager \
  -p 5432:5432 \
  -d postgres:16
```

### 2. Настроить переменные окружения

```bash
cp .env.example .env
```

### 3. Применить миграции

Требуется установленный [golang-migrate](https://github.com/golang-migrate/migrate):

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/taskmanager?sslmode=disable" up
```

### 4. Запустить сервер

```bash
go run main.go
```

## API

| Метод  | Путь            | Описание                      |
|--------|-----------------|--------------------------------|
| GET    | `/tasks`        | Список всех задач              |
| POST   | `/tasks`        | Создать задачу                 |
| GET    | `/tasks/{id}`   | Получить задачу по ID          |
| PATCH  | `/tasks/{id}`   | Частично обновить задачу       |
| DELETE | `/tasks/{id}`   | Удалить задачу                 |
| GET    | `/metrics`      | Метрики Prometheus             |

### Пример запроса — создание задачи

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go"}'
```

Ответ:

```json
{"id":1,"title":"Learn Go","done":false}
```

### Пример запроса — частичное обновление

```bash
curl -X PATCH http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"done":true}'
```

### Пример запроса — удаление

```bash
curl -X DELETE http://localhost:8080/tasks/1
```

Ответ: `204 No Content`

## Обработка ошибок

| Статус | Когда возвращается                                     |
|--------|------------------------------------------------------------|
| 400    | Невалидный JSON, невалидный `id`, пустой `title`             |
| 404    | Задача с указанным `id` не найдена                           |
| 500    | Непредвиденная ошибка сервера (например, БД недоступна)      |

## Docker-образ

Сборка многоэтапная (multi-stage): на первом этапе компилируется бинарник в полном окружении Go, на втором — копируется только сам бинарник в минимальный `alpine`-образ. Итоговый образ приложения — около 18 МБ, без установленного Go, компилятора и исходников внутри.

## Миграции

Файлы миграций лежат в `migrations/`, пара `up`/`down` на каждое изменение схемы. При запуске через `docker-compose up` применяются автоматически. Для ручного управления (при локальном запуске без Docker):

```bash
# применить все новые миграции
migrate -path migrations -database "$DATABASE_URL" up

# откатить последнюю миграцию
migrate -path migrations -database "$DATABASE_URL" down 1
```

## Метрики и мониторинг

Пример PromQL-запросов, которые можно выполнить в веб-интерфейсе Prometheus (`http://localhost:9090`):

```promql
# текущее общее число запросов по каждому методу/пути/статусу
http_requests_total

# частота запросов в секунду за последнюю минуту
rate(http_requests_total[1m])
```

## Планы по развитию

- [ ] Юнит-тесты для service-слоя (с mock-реализацией `TaskRepository`)
- [ ] Grafana для визуализации метрик поверх Prometheus