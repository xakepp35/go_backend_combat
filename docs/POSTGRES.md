# PostgreSQL: от SQL-песочницы к базе приложения

[Окружение 09](ENVIRONMENT.md) · [SQL12](../student/12_SQL.md) · [Go + PostgreSQL14](../student/14_POSTGRES.md)

## 1. База уже в Compose

С недели 09 у твоего проекта есть `compose.yaml`. Все следующие команды выполняй из его корня; для знакомства с эталоном можно перейти в `reference/12` или `reference/14`. Не запускай разные снимки одновременно под одним именем Compose-проекта.

```sh
docker compose up -d --wait postgres
docker compose exec postgres psql -U jobflow -d jobflow
```

`psql` работает **в контейнере**, на хост его ставить не нужно. `\conninfo` показывает подключение, `\dt` — таблицы, `\q` — выход. SQL заканчивается точкой с запятой. Приглашение `jobflow-#` вместо `jobflow=#` часто означает незавершённую команду; `\r` очищает буфер ввода.

## 2. На неделях 12–13 экспериментируем отдельно

Один раз создай учебную базу:

```sh
docker compose exec postgres psql -U jobflow -d postgres -c "CREATE DATABASE jobflow_lab;"
docker compose exec postgres psql -U jobflow -d jobflow_lab
```

Сообщение «already exists» при повторном создании означает, что база уже есть; просто подключись. Перед любым изменением проверь `\conninfo`. В `jobflow_lab` выполняй упражнения из `sql/practice.sql`: создавать и менять таблицы здесь можно без влияния на приложение. Учебные сырые SQL-файлы недель 12–13 ещё не прогоняются приложением.

Для сохранённого файла на Linux/macOS: `docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U jobflow -d jobflow_lab < sql/practice.sql`. PowerShell: `Get-Content -Raw sql/practice.sql | docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U jobflow -d jobflow_lab`. Оба варианта читают локальный файл и передают SQL в контейнер. `ON_ERROR_STOP` прекращает выполнение при ошибке; это не оборачивает весь файл в транзакцию.

## 3. С недели 14 схемой приложения владеет встроенный Goose

Приложение подключается к **jobflow**, не к jobflow_lab. SQL лежит в `migrations/*.sql`, а `migrations/migrations.go` включает его через `//go:embed *.sql`. В бинарнике конкретной версии находятся и код, и изменения её схемы.

```sh
docker compose up -d --build --wait app
docker compose logs --tail=80 app
```

Отдельной команды миграции и отдельного сервиса нет. Старт: открыть pgxpool → проверить подключение → применить pending-миграции → восстановить фоновую работу, когда она появится → запустить HTTP. Ошибка миграции завершает процесс с ненулевым кодом. `/healthz` не обслуживается до завершения миграций.

Для быстрого запуска Go на хосте сначала подними postgres, затем задай `DATABASE_URL`:

```sh
# Linux/macOS, порт из .env при необходимости замени
export DATABASE_URL='postgres://jobflow:jobflow@127.0.0.1:15432/jobflow?sslmode=disable'
go run ./cmd/jobflow
```

PowerShell: `$env:DATABASE_URL = 'postgres://jobflow:jobflow@127.0.0.1:15432/jobflow?sslmode=disable'`, затем `go run ./cmd/jobflow`. Внутри Compose адрес другой: `postgres:5432`. `.env` Compose сам по себе не экспортирует переменные в твой терминал.

## 4. Новое изменение — новый файл

```sql
-- +goose Up
ALTER TABLE tasks ADD COLUMN estimate_minutes integer NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE tasks DROP COLUMN estimate_minutes;
```

Это **пример следующей** миграции, не файл для копирования поверх `001_tasks.sql`. Согласуй правило допустимой длительности и номер следующей версии. Шаблон можно создать инструментом: `go run github.com/pressly/goose/v3/cmd/goose@v3.29.0 -dir migrations create add_task_estimate sql`. Первый запуск скачает и соберёт CLI. Он нужен для генерации файла, не является отдельным сервисом и не нужен для запуска приложения. Для команды `create` подключение к базе не требуется.

Goose хранит историю в `goose_db_version`. Проверить её можно через `psql`: `SELECT version_id, is_applied FROM goose_db_version ORDER BY id;`. Повторный старт не выполняет уже учтённые версии заново. Применённые файлы не редактируем: Goose не сравнивает их содержимое с тем, что когда-то выполнялось.

Session advisory lock сериализует миграции от нескольких экземпляров **с тем же ключом блокировки**. В учебном коде есть общий таймаут старта 60 с. Одна обычная SQL-миграция выполняется в транзакции; весь набор файлов не является одной общей транзакцией. `NO TRANSACTION` требует отдельного обоснования. Изменение данных через `Down` при каждом выключении **не выполняется**.

## 5. Почему в коде Goose виден database/sql

Все запросы приложения идут через `github.com/jackc/pgx/v5/pgxpool`. Goose принимает `*sql.DB`, поэтому `stdlib.OpenDBFromPool(pool)` даёт ему адаптер к **тому же** pgxpool. Это не отдельный драйвер для CRUD и не второй независимый сетевой пул. Адаптер закрываем после миграций, основной пул — при завершении приложения. Такое поведение описано в [pgx stdlib](https://pkg.go.dev/github.com/jackc/pgx/v5/stdlib#OpenDBFromPool).

Для настоящего production: оцени длительность DDL, совместимость со старым кодом, права на изменение схемы и стратегию выпуска. Прикладной старт с миграциями удобен для этого проекта, но не обещает безостановочный rollout любой SQL-команды. Разбор — [урок 21](../student/21_OPERATIONS.md).

## 6. Если ты уже запускал старую редакцию курса

Старая самописная история `course_migrations` не равна истории Goose. Новая сборка обнаруживает её и **останавливается с объяснением**, не помечая неизвестную схему как корректную. Данные не удаляются. Для переноса есть [отдельная памятка](POSTGRES_UPGRADE.md). Для новых проектов никаких дополнительных действий нет.

## 7. Остановка и проверки

`docker compose down` сохраняет именованный том. Не удаляй том ради исправления миграции. Для диагностики сначала логи, история версий и точная ошибка SQL. Бэкап нужен до опасного изменения, а не после удачного `DROP`.

Существующие интеграционные проверки курса используют отдельную временную схему и `TEST_DATABASE_URL`. Не направляй их на рабочую базу с важными данными. Команды и область проверок описаны в [checks/README](../checks/README.md); запуск проверок — отдельное действие.

Первоисточники: [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool), [Goose Provider](https://pkg.go.dev/github.com/pressly/goose/v3#NewProvider), [блокировка Goose](https://pkg.go.dev/github.com/pressly/goose/v3/lock#NewPostgresSessionLocker).
