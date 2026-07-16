# Task Manager API

Простой REST API для управления задачами (todo-list), написанный на Go.

Проект построен по слоистой архитектуре (`handler → service → repository`), что позволяет менять хранилище данных, не затрагивая бизнес-логику и HTTP-обработчики.

## Стек

- Go 1.22+ (стандартная библиотека `net/http`, без сторонних роутеров)
- PostgreSQL 16
- `database/sql` + `lib/pq` (драйвер)
- `golang-migrate` — миграции БД
- `godotenv` — конфигурация через переменные окружения

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
migrations/           — SQL-миграции (up/down)
```

Каждый слой знает только о слое ниже себя:

- **handler** знает про **service**
- **service** знает про **repository** — но только через интерфейс `TaskRepository`, не про конкретную реализацию
- **repository** ничего не знает о вышестоящих слоях

Благодаря этому переключение хранилища с in-memory на PostgreSQL потребовало изменения всего одной строки в `main.go` — `service` и `handler` остались нетронутыми.

### Обработка ошибок между слоями

`repository` возвращает собственную ошибку `ErrNotFound`, если запись не найдена, либо любую другую ошибку в случае сбоя (например, недоступность БД). `service` транслирует `repository.ErrNotFound` в свою `service.ErrTaskNotFound`, не пропуская наверх детали реализации хранилища. `handler` проверяет конкретную ошибку через `errors.Is` и выбирает подходящий HTTP-статус (`404` для "не найдено", `500` для непредвиденных сбоев).

## Запуск

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

При необходимости отредактируйте значения в `.env` под вашу конфигурацию.

### 3. Применить миграции

Требуется установленный [golang-migrate](https://github.com/golang-migrate/migrate):

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/taskmanager?sslmode=disable" up
```

### 4. Запустить сервер

```bash
go run main.go
```

Сервер запустится на `http://localhost:8080`.

## API

| Метод  | Путь            | Описание                      |
|--------|-----------------|--------------------------------|
| GET    | `/tasks`        | Список всех задач              |
| POST   | `/tasks`        | Создать задачу                 |
| GET    | `/tasks/{id}`   | Получить задачу по ID          |
| PATCH  | `/tasks/{id}`   | Частично обновить задачу       |
| DELETE | `/tasks/{id}`   | Удалить задачу                 |

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
|--------|----------------------------------------------------------|
| 400    | Невалидный JSON, невалидный `id`, пустой `title`          |
| 404    | Задача с указанным `id` не найдена                        |
| 500    | Непредвиденная ошибка сервера (например, БД недоступна)   |

## Миграции

Файлы миграций лежат в `migrations/`, пара `up`/`down` на каждое изменение схемы:

```bash
# применить все новые миграции
migrate -path migrations -database "$DATABASE_URL" up

# откатить последнюю миграцию
migrate -path migrations -database "$DATABASE_URL" down 1
```

## Планы по развитию

- [ ] Docker + docker-compose (приложение + БД одной командой)
- [ ] Метрики (Prometheus)
- [ ] Юнит-тесты для service-слоя (с mock-реализацией `TaskRepository`)
- [ ] Логирование запросов (middleware)