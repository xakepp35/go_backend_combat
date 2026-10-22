# Go Backend Combat — 05. DATA

**Version:** 1.0.0
**Topic:** Data & Transactions
**Format:** Дидактический материал ученика — PostgreSQL, SQL, migrations, pgxpool, transactions и ACID
**Project:** JobFlow
**Duration:** 120 минут
**Prerequisite:** `01_LANGUAGE.md`, `02_CONCURRENCY.md`, `04_DESIGN.md`

Плакаты:

* [05A — Data Model](../posters/05A_DATA_MODEL.svg)
* [05B — SQL Schema](../posters/05B_SQL_SCHEMA.svg)
* [05C — Transactions & Concurrency](../posters/05C_TRANSACTIONS_CONCURRENCY.svg)

---

# 1. Цель занятия

До этого JobFlow был в основном in-memory:

```text
Job
 ↓
WorkerPool
 ↓
MemoryJobRepository
```

Теперь состояние переносится в PostgreSQL.

Но цель занятия не ограничивается подключением database driver.

Основная модель:

> **State → Invariant → Constraint → SQL → Transaction → Commit**

Нужно научиться:

```text
моделировать состояние
      ↓
защищать invariant
      ↓
выражать его schema
      ↓
писать SQL
      ↓
работать с PostgreSQL из Go
      ↓
корректно менять state
      ↓
защищать concurrent transitions
```

---

# 2. Что изменится в JobFlow

Было:

```text
MemoryRepository
```

Станет:

```text
PostgreSQL schema
      ↓
SQL
      ↓
PostgresJobRepository
      ↓
pgxpool
```

А затем:

```text
concurrent workers
      ↓
transaction
      ↓
row lock / isolation
      ↓
consistent state
```

---

# 3. 05A — Data Model

![05A — Data Model](../posters/05A_DATA_MODEL.svg)

Плакат задаёт основные вопросы:

```text
ENTITY
IDENTITY
RELATIONSHIP
VALID STATE
NULL
STATE OWNER
```

---

# 4. Entity

Основная persisted entity:

```text
Job
```

В database она становится:

```text
jobs
 ├── id
 ├── type
 ├── payload
 ├── status
 ├── created_at
 └── updated_at
```

`id` определяет identity.

Остальные поля описывают текущее состояние.

---

# 5. Primary Key

Основной identity constraint:

```sql
id uuid PRIMARY KEY
```

Primary key:

* идентифицирует строку;
* не допускает duplicate values;
* не допускает `NULL`.

В JobFlow:

```text
Job.ID
    ↓
jobs.id
    ↓
PRIMARY KEY
```

---

# 6. Relationships

Job может иметь несколько попыток обработки:

```text
jobs
  │
  └── 1:N ── job_attempts
```

Например:

```text
Job #1
 ├── attempt 1
 ├── attempt 2
 └── attempt 3
```

Поэтому attempts естественно представить отдельной таблицей.

Пример:

```text
job_attempts
 ├── id
 ├── job_id
 ├── attempt
 ├── status
 ├── started_at
 ├── finished_at
 ├── error
 └── created_at
```

---

# 7. Foreign Key

Связь:

```sql
job_id uuid NOT NULL REFERENCES jobs(id)
```

означает:

> Attempt может ссылаться только на существующую Job.

Database становится владельцем этого invariant.

Нельзя создавать:

```text
attempt → несуществующая job
```

---

# 8. Constraints

Основные типы ограничений:

```text
PRIMARY KEY
FOREIGN KEY
NOT NULL
UNIQUE
CHECK
```

Они превращают business/data rules в исполняемые database constraints.

Пример:

```sql
type text NOT NULL
```

защищает:

> `jobs.type` не может быть `NULL`.

Важно:

> Application validation и database constraints не заменяют друг друга.

Application validation улучшает API и feedback.

Database constraints защищают persisted state.

---

# 9. NULL

`NULL` не равен:

```text
""
0
false
```

`NULL` означает отсутствие значения.

Неправильно:

```sql
WHERE finished_at = NULL
```

Правильно:

```sql
WHERE finished_at IS NULL
```

или:

```sql
WHERE finished_at IS NOT NULL
```

Mental model:

```text
NULL
 ↓
unknown / absent value
```

Не следует автоматически трактовать `NULL` как «пустое значение».

---

# 10. SQL Schema JobFlow

Базовая `jobs` table:

```sql
CREATE TABLE jobs (
    id uuid PRIMARY KEY,
    type text NOT NULL,
    payload jsonb NOT NULL,
    status text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
```

`job_attempts`:

```sql
CREATE TABLE job_attempts (
    id uuid PRIMARY KEY,
    job_id uuid NOT NULL REFERENCES jobs(id),
    attempt integer NOT NULL,
    status text NOT NULL,
    started_at timestamptz NOT NULL,
    finished_at timestamptz,
    error text,
    created_at timestamptz NOT NULL
);
```

---

# 11. Migration

![05B — SQL Schema](../posters/05B_SQL_SCHEMA.svg)

Migration — это versioned transition schema.

```text
current schema
     ↓
migration
     ↓
new schema
```

Например:

```text
00001 → jobs
00002 → job_attempts
00003 → constraints
00004 → indexes
```

Migration должна быть:

```text
versioned
reproducible
reviewable
ordered
```

---

# 12. Goose

Пример migration:

```sql
-- +goose Up

CREATE TABLE jobs (
    id uuid PRIMARY KEY,
    type text NOT NULL,
    payload jsonb NOT NULL,
    status text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

-- +goose Down

DROP TABLE jobs;
```

Применение:

```bash
goose -dir migrations postgres "$DATABASE_URL" up
```

Проверка:

```bash
goose -dir migrations postgres "$DATABASE_URL" status
```

Главная проверка:

> Умеет ли schema полностью восстановиться из migrations на чистой database?

---

# 13. Migration как transition

Schema:

```text
S0
 ↓
S1
 ↓
S2
```

Migration:

```text
S0 → S1
S1 → S2
```

Это позволяет одинаково воспроизводить database state на:

```text
development
test
staging
production
```

---

# 14. DDL

Основные операции:

```sql
CREATE TABLE ...
ALTER TABLE ...
CREATE INDEX ...
DROP ...
```

DDL описывает структуру persisted state.

Пример foreign key:

```sql
job_id uuid NOT NULL REFERENCES jobs(id)
```

---

# 15. INSERT

Создание Job:

```sql
INSERT INTO jobs (
    id,
    type,
    payload,
    status,
    created_at,
    updated_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    now(),
    now()
);
```

Параметры:

```text
$1
$2
$3
$4
```

нужно использовать вместо конкатенации SQL и пользовательских данных.

Причины:

```text
SQL injection protection
correct parameter encoding
clear query structure
```

---

# 16. SELECT

Получение Job:

```sql
SELECT
    id,
    type,
    payload,
    status,
    created_at,
    updated_at
FROM jobs
WHERE id = $1;
```

Главное:

> SQL должен выражать конкретный access pattern.

---

# 17. UPDATE

Переход состояния:

```sql
UPDATE jobs
SET
    status = $2,
    updated_at = now()
WHERE id = $1;
```

`WHERE` определяет множество изменяемых строк.

Без `WHERE` потенциально изменяются все строки.

---

# 18. DELETE

Удаление:

```sql
DELETE FROM jobs
WHERE id = $1;
```

Нужно понимать:

```text
which rows
why those rows
what invariant remains after deletion
```

---

# 19. JOIN

Получить Job вместе с attempts:

```sql
SELECT
    j.id,
    j.type,
    a.attempt,
    a.status
FROM jobs j
JOIN job_attempts a
    ON a.job_id = j.id
WHERE j.id = $1;
```

`JOIN` использует отношение:

```text
job_attempts.job_id
        ↓
jobs.id
```

---

# 20. INNER JOIN vs LEFT JOIN

`JOIN` / `INNER JOIN`:

> Возвращает только пары, для которых найдено соответствие.

`LEFT JOIN`:

> Сохраняет строки левой таблицы даже без matching row справа.

Например:

```sql
SELECT
    j.id,
    a.attempt
FROM jobs j
LEFT JOIN job_attempts a
    ON a.job_id = j.id;
```

Позволяет получить Job даже без attempts.

---

# 21. Index

Допустим:

```sql
SELECT *
FROM jobs
WHERE status = 'running';
```

Если такой access pattern выполняется часто, может понадобиться index:

```sql
CREATE INDEX idx_jobs_status
ON jobs(status);
```

Но:

> **Index не является автоматической оптимизацией.**

Он имеет стоимость:

```text
memory
write amplification
maintenance
storage
```

и полезен только относительно реального query workload.

---

# 22. EXPLAIN

План запроса:

```sql
EXPLAIN
SELECT *
FROM jobs
WHERE status = 'running';
```

Фактическое выполнение:

```sql
EXPLAIN ANALYZE
SELECT *
FROM jobs
WHERE status = 'running';
```

`EXPLAIN ANALYZE` даёт evidence фактического execution behavior:

```text
plan
actual rows
execution time
scan choice
```

---

# 23. pgxpool/v5

PostgreSQL из Go подключается через pool.

Пример:

```go
pool, err := pgxpool.New(ctx, databaseURL)
if err != nil {
	return err
}
defer pool.Close()
```

Pool управляет несколькими connections.

Если:

```text
1000 goroutines
```

и:

```text
MaxConns = 20
```

это не означает 1000 одновременных connections.

Максимум активно используемых pool connections ограничен конфигурацией pool.

Остальные операции ждут свободного ресурса или получают ошибку согласно выбранной policy.

Это связь:

```text
Go concurrency
      +
DB resource bound
```

---

# 24. pgxpool Configuration

Основные параметры:

```text
MinConns
MaxConns
MaxConnLifetime
MaxConnIdleTime
```

Они должны определяться workload и ресурсными ограничениями.

Не следует выбирать `MaxConns` просто как большое число.

Нужно учитывать:

```text
CPU
PostgreSQL capacity
query latency
number of application instances
other database clients
```

---

# 25. Repository Adapter

Пример concrete adapter:

```go
type PostgresJobRepository struct {
	pool *pgxpool.Pool
}
```

Получение Job:

```go
func (r *PostgresJobRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (Job, error) {
	var job Job

	err := r.pool.QueryRow(
		ctx,
		`SELECT
			id,
			type,
			payload,
			status,
			created_at,
			updated_at
		 FROM jobs
		 WHERE id = $1`,
		id,
	).Scan(
		&job.ID,
		&job.Type,
		&job.Payload,
		&job.Status,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		return Job{}, err
	}

	return job, nil
}
```

Request context должен передаваться вниз:

```text
HTTP
 ↓
service
 ↓
repository
 ↓
pgx
```

чтобы cancellation/deadline распространялись на database operation.

---

# 26. Transaction

![05C — Transactions & Concurrency](../posters/05C_TRANSACTIONS_CONCURRENCY.svg)

Transaction задаёт логическую границу изменения состояния.

Пример:

```sql
BEGIN;

UPDATE ...
INSERT ...

COMMIT;
```

Если операция не может завершиться корректно:

```text
ROLLBACK
```

Главный вопрос:

> Какие изменения должны быть зафиксированы вместе?

---

# 27. Lost Update

Пусть две concurrent transactions читают:

```text
attempt = 3
```

Transaction A:

```text
3 → 4
```

Transaction B:

```text
3 → 4
```

После двух изменений:

```text
expected = 5
actual   = 4
```

Это lost update.

Причина:

```text
read
 ↓
modify
 ↓
write
```

выполнялась конкурентно без достаточной coordination.

---

# 28. Row Lock

Для state transition можно использовать:

```sql
BEGIN;

SELECT attempt
FROM jobs
WHERE id = $1
FOR UPDATE;

UPDATE jobs
SET attempt = attempt + 1
WHERE id = $1;

COMMIT;
```

`FOR UPDATE` блокирует выбранные строки для конфликтующего изменения в текущей transaction.

Вторая конкурентная transaction должна дождаться разрешения конфликта или получить соответствующий результат согласно transaction semantics.

---

# 29. Transaction State Transition

Для JobFlow:

```text
pending
   ↓
running
```

и создание attempt могут быть одной логической операцией.

Поэтому:

```text
BEGIN
 ↓
lock Job
 ↓
check current state
 ↓
update Job
 ↓
insert Attempt
 ↓
COMMIT
```

Если любая часть неуспешна:

```text
ROLLBACK
```

Иначе можно получить частично применённое состояние.

---

# 30. MVCC

PostgreSQL использует MVCC — Multi-Version Concurrency Control.

Вместо модели:

```text
one global database mutex
```

работа с данными организована вокруг версий и visibility rules.

Это позволяет concurrent transactions получать согласованные представления данных при определённых isolation guarantees.

Основная mental model:

```text
concurrent transactions
        ↓
visibility + versions
        ↓
isolation semantics
```

---

# 31. Isolation Level

Isolation level определяет гарантии конкурентного выполнения и видимости состояния.

Важные уровни:

```text
READ COMMITTED
REPEATABLE READ
SERIALIZABLE
```

PostgreSQL по умолчанию использует:

```text
READ COMMITTED
```

Не нужно сводить Isolation к фразе:

> «Транзакции полностью изолированы.»

Isolation имеет конкретную семантику выбранного уровня.

---

# 32. Serializable

`SERIALIZABLE` стремится обеспечить эффект, эквивалентный некоторому последовательному выполнению конфликтующих транзакций.

При обнаружении невозможного concurrent schedule транзакция может завершиться serialization failure.

Такую операцию можно повторить целиком:

```text
BEGIN
 ↓
work
 ↓
serialization failure
 ↓
ROLLBACK
 ↓
retry transaction
```

Retry должен иметь:

```text
bounded attempts
backoff
idempotent operation semantics
```

если это требуется архитектурой.

---

# 33. ACID

На интервью:

## A — Atomicity

Транзакция либо фиксирует изменения как единое целое, либо не оставляет частично зафиксированный результат.

## C — Consistency

После commit база остаётся в состоянии, удовлетворяющем constraints и инвариантам.

## I — Isolation

Concurrent transactions получают определённые гарантии видимости и взаимодействия в рамках выбранного isolation level.

## D — Durability

После успешного commit изменения сохраняются в пределах гарантий СУБД и её конфигурации.

Главная поправка:

> **ACID не означает «транзакция безопасна вообще». Каждое свойство имеет конкретную семантику, особенно Isolation.**

---

# 34. Database Constraints как последняя граница

Application validation:

```text
HTTP
 ↓
service
 ↓
validation
```

не должна быть единственной защитой persisted state.

Например:

```sql
job_id uuid NOT NULL REFERENCES jobs(id)
```

database сама не позволит сохранить невозможную relationship.

Модель:

```text
application validation
        +
database constraints
        ↓
correct persisted state
```

---

# 35. Transaction в Go

Пример:

```go
tx, err := pool.Begin(ctx)
if err != nil {
	return err
}
defer tx.Rollback(ctx)
```

Operations:

```go
_, err = tx.Exec(
	ctx,
	`UPDATE jobs
	 SET attempt = attempt + 1,
	     updated_at = now()
	 WHERE id = $1`,
	id,
)
if err != nil {
	return err
}
```

Commit:

```go
return tx.Commit(ctx)
```

`defer tx.Rollback(ctx)` можно использовать даже перед успешным commit как safety net для error paths.

---

# 36. Практическая задача — Start Job

Реализовать атомарную операцию:

> `pending → running`

одновременно с созданием `job_attempt`.

Ожидаемый flow:

```text
BEGIN
  ↓
lock Job
  ↓
verify status = pending
  ↓
status → running
  ↓
create attempt
  ↓
COMMIT
```

При ошибке:

```text
ROLLBACK
```

---

# 37. Concurrent Start

Запустить две concurrent попытки:

```text
worker A
    ↓
start(job)

worker B
    ↓
start(job)
```

Обе пытаются сделать:

```text
pending → running
```

Корректное состояние должно иметь:

```text
один успешный переход
```

а не:

```text
две независимые попытки
```

Нужно проверить:

```text
final Job state
attempt count
error semantics
transaction behavior
```

---

# 38. Checkpoint

### Что такое primary key?

<details><summary>Ответ</summary>Уникальный идентификатор строки, не допускающий NULL.</details>

### Зачем foreign key?

<details><summary>Ответ</summary>Для защиты ссылочной целостности между таблицами.</details>

### Что такое `NULL`?

<details><summary>Ответ</summary>Отсутствующее/неизвестное значение, не равное пустой строке или нулю.</details>

### Почему `= NULL` неправильно?

<details><summary>Ответ</summary>Для проверки NULL используются `IS NULL` и `IS NOT NULL`.</details>

### Зачем migration?

<details><summary>Ответ</summary>Для versioned и воспроизводимого изменения schema.</details>

### Зачем pgxpool?

<details><summary>Ответ</summary>Для управления пулом database connections и ограничения database concurrency.</details>

### Что такое transaction?

<details><summary>Ответ</summary>Логическая атомарная граница группы database operations.</details>

### Что такое lost update?

<details><summary>Ответ</summary>Сценарий, когда concurrent read-modify-write приводит к потере одного из обновлений.</details>

### Что делает `FOR UPDATE`?

<details><summary>Ответ</summary>Берёт row-level lock для выбранных строк в текущей transaction.</details>

### Что такое MVCC?

<details><summary>Ответ</summary>Механизм конкурентного доступа через версии данных и правила видимости.</details>

### Что такое ACID?

<details><summary>Ответ</summary>Atomicity, Consistency, Isolation, Durability.</details>

### Что делать при serialization failure?

<details><summary>Ответ</summary>Откатить transaction и при подходящей policy повторить её целиком.</details>

---

# 39. JobFlow v3 — PostgreSQL Architecture

После урока:

```text
                    HTTP
                      │
                      ▼
                 JobService
                      │
              ┌───────┴────────┐
              ▼                ▼
       JobRepository     EventPublisher
              │
              ▼
           pgxpool
              │
              ▼
         PostgreSQL
           │      │
           ▼      ▼
         jobs   attempts
```

Schema lifecycle:

```text
Go domain model
      ↓
migration
      ↓
schema
      ↓
SQL
      ↓
repository
      ↓
transaction
      ↓
consistent state
```

---

# 40. Домашнее задание — JobFlow v3

Цель:

> **Убрать основную runtime-зависимость JobFlow от MemoryRepository и сделать PostgreSQL настоящим владельцем persisted state.**

## A. Migrations

Создать минимум:

```text
00001_jobs.sql
00002_job_attempts.sql
00003_constraints.sql
00004_indexes.sql
```

### `jobs`

```text
id
type
payload
status
created_at
updated_at
```

### `job_attempts`

```text
id
job_id
attempt
status
started_at
finished_at
error
created_at
```

Добавить подходящие:

```text
PRIMARY KEY
FOREIGN KEY
NOT NULL
CHECK
UNIQUE
INDEX
```

---

# 41. PostgreSQL Repository

Реализовать:

```text
Add
Get
List
UpdateStatus
CreateAttempt
```

Concrete implementation:

```go
type PostgresJobRepository struct {
	pool *pgxpool.Pool
}
```

---

# 42. Pool

Добавить configuration:

```text
MinConns
MaxConns
MaxConnLifetime
MaxConnIdleTime
```

Проверить:

```text
pool creation
query
scan
context cancellation
close
```

---

# 43. SQL Practice

Самостоятельно написать и выполнить:

```text
INSERT
SELECT
UPDATE
DELETE
JOIN
COUNT
ORDER BY
LIMIT
```

Минимум три запроса исследовать через:

```sql
EXPLAIN
```

и:

```sql
EXPLAIN ANALYZE
```

Для каждого уметь объяснить:

> Какой access pattern используется и почему planner выбрал этот plan?

---

# 44. Transaction Practice

Реализовать:

```text
pending → running
```

атомарно с:

```text
job_attempt INSERT
```

Требования:

```text
transaction
row locking
status verification
commit
rollback
context propagation
```

---

# 45. Concurrency Test

Запустить две concurrent операции для одной Job.

Проверить:

```text
exactly one valid transition
consistent attempt state
no invalid final state
no lost update
```

Запустить:

```bash
go test ./...
go test -race ./...
go vet ./...
```

---

# 46. Clean Database Test

Удалить schema и поднять её заново только migrations.

Проверить:

```text
clean database
 ↓
goose up
 ↓
schema exists
 ↓
constraints exist
 ↓
indexes exist
 ↓
application starts
```

Это подтверждает reproducibility schema.

---

# 47. Acceptance Criteria

## PASS

Ученик:

* проектирует `jobs` и `job_attempts`;
* использует primary/foreign keys;
* применяет `NOT NULL`, `CHECK`, `UNIQUE` там, где они нужны;
* понимает `NULL`;
* пишет DDL и DML;
* пишет `JOIN`;
* умеет объяснить index;
* умеет читать простой `EXPLAIN`;
* использует Goose;
* поднимает schema с нуля migrations;
* подключает `pgxpool/v5`;
* передаёт context;
* реализует repository;
* использует transaction;
* объясняет lost update;
* использует row locking;
* объясняет MVCC;
* объясняет ACID;
* понимает isolation levels;
* обрабатывает serialization failure;
* тестирует concurrent state transitions.

---

# 48. Финальная mental model

PostgreSQL в JobFlow — не просто место хранения `struct Job`.

Это владелец durable state:

```text
                 STATE
                   │
                   ▼
               INVARIANT
                   │
                   ▼
                CONSTRAINT
                   │
                   ▼
                 SCHEMA
                   │
                   ▼
                  SQL
                   │
                   ▼
              TRANSACTION
                   │
             ┌─────┴─────┐
             ▼           ▼
          COMMIT      CONFLICT
                         │
                         ▼
                       RETRY
```

Главная мысль:

> **Database engineering — это управление корректным состоянием при изменениях и конкуренции.**
