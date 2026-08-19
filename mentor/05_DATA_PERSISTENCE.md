# Go Backend Combat — 05. DATA

**Version:** 1.0.0
**Owner:** Mentor
**Duration:** 120 минут
**Format:** интерактивный mentor-led lesson
**Project:** JobFlow
**Prerequisite:** `01.LANGUAGE.md`, `02.CONCURRENCY.md`, `04.DESIGN.md`
**Posters:** `05A_DATA_MODEL.svg`, `05B_SQL_SCHEMA.svg`, `05C_TRANSACTIONS_CONCURRENCY.svg`
**Outcome:** ученик проектирует PostgreSQL schema, пишет SQL и migrations, подключает `pgxpool/v5`, использует transactions и понимает ACID/MVCC/isolation на практическом примере

---

# 1. Purpose

До этого момента JobFlow был в основном in-memory системой.

У нас есть:

```text
Job
 ↓
WorkerPool
 ↓
MemoryJobRepository
```

Теперь мы переносим состояние в PostgreSQL.

Но задача урока не:

> «Научиться подключаться к PostgreSQL».

Задача:

> **Научиться проектировать, изменять и конкурентно изменять устойчивое состояние backend-системы.**

Поэтому урок идёт через три модели:

```text
05A DATA MODEL
    ↓
05B SQL & SCHEMA
    ↓
05C TRANSACTIONS & CONCURRENCY
```

И одновременно развивается JobFlow:

```text
MemoryRepository
      ↓
PostgreSQL schema
      ↓
Go repository
      ↓
pgxpool
      ↓
SQL queries
      ↓
transactions
      ↓
concurrent state transitions
```

Главная формула:

> **State → Invariant → Constraint → SQL → Transaction → Commit**

---

# 2. Learning Contract

После занятия ученик умеет:

### PostgreSQL

* объяснить table / row / column;
* определить primary key;
* использовать `NOT NULL`, `UNIQUE`, `CHECK`, `FOREIGN KEY`;
* объяснить `NULL`;
* писать `CREATE TABLE`;
* писать `INSERT`, `SELECT`, `UPDATE`, `DELETE`;
* писать `JOIN`;
* понимать базовое назначение index;
* читать простой `EXPLAIN`;
* объяснить transaction;
* объяснить MVCC;
* различать isolation levels;
* объяснить lost update;
* использовать row locking;
* понимать serialization failure;
* объяснить ACID на уровне интервью.

### Go

* подключать PostgreSQL через `pgxpool/v5`;
* создавать pool;
* задавать pool configuration;
* выполнять query через `context.Context`;
* сканировать строки;
* выполнять transaction;
* корректно `Commit` / `Rollback`;
* отделять SQL adapter от application policy.

### Tooling

* создавать Goose migration;
* применять migration;
* проверять schema;
* откатывать migration там, где это предусмотрено;
* держать schema changes versioned и reproducible.

### Engineering

Ученик умеет пройти:

```text
требование
→ модель
→ invariant
→ schema
→ migration
→ SQL
→ Go repository
→ transaction
→ test
```

---

# 3. Mentor Flow

Работаем по схеме:

```text
LOOK
 ↓
ASK
 ↓
PREDICT
 ↓
WRITE SQL
 ↓
RUN
 ↓
OBSERVE
 ↓
IMPLEMENT IN GO
 ↓
BREAK WITH CONCURRENCY
 ↓
FIX
 ↓
VERIFY
```

Особое правило:

> **Сначала ученик должен предсказать результат SQL или транзакции, потом мы запускаем.**

Особенно это важно для:

* `NULL`;
* `JOIN`;
* constraints;
* `UPDATE`;
* transactions;
* concurrent sessions;
* isolation.

---

# 4. 00:00–00:08 — Вход в PostgreSQL

Открываю JobFlow.

Говорю:

> «До сих пор если процесс умер, наши Job исчезли. Давайте теперь зададим вопрос, который backend задаёт постоянно: где находится состояние системы после того, как процесс завершился?»

Ответ:

<details style="display: inline;"><summary>Подсказка</summary>Нам нужен внешний durable state.</details>

<details style="display: inline;"><summary>Ответ</summary>В постоянном внешнем хранилище, в нашем случае PostgreSQL.</details>

Продолжаю:

> «Но просто положить объект в таблицу недостаточно. Теперь database становится владельцем состояния, а значит мы должны научиться описывать допустимые состояния и изменения этого состояния.»

Открываю `05A`.

---

# 5. 00:08–00:20 — Плакат 05A: DATA MODEL

## Card 01 — ENTITY

> «Что является сущностью в нашем JobFlow?»

<details><summary>Ответ</summary>`Job` и его persisted state.</details>

Показываем:

```text
jobs
id
type
payload
status
created_at
```

Спрашиваю:

> «Что здесь является состоянием, а что identity?»

<details><summary>Ответ</summary>`id` задаёт identity, остальные поля описывают текущее состояние Job.</details>

---

## Card 02 — IDENTITY

> «Как database должна гарантировать, что две Job не имеют одну identity?»

<details><summary>Ответ</summary>Через `PRIMARY KEY`.</details>

Пишем:

```sql
id uuid PRIMARY KEY
```

Что ещё даёт primary key кроме уникальности?

<details><summary>Ответ</summary>Он задаёт identity строки и запрещает `NULL` в ключе.</details>

---

## Card 03 — RELATIONSHIP

Теперь я говорю:

> «Одной таблицы нам скоро станет мало.»

Создадим попытки выполнения Job:

```text
jobs
  │
  └── 1:N ── job_attempts
```

Почему это отдельная сущность, а не ещё три поля в `jobs`?

<details><summary>Подсказка</summary>Нам нужно хранить много attempts для одной Job.</details>

<details><summary>Ответ</summary>Attempt имеет собственную множественность и lifecycle, поэтому отдельная таблица лучше выражает отношение 1:N.</details>

---

## Card 04 — VALID STATE

Теперь:

> «Где лучше защищать invariant: в Go или PostgreSQL?»

Например:

> `jobs.type` никогда не должен быть `NULL`.

<details><summary>Ответ</summary>Критический invariant должен защищаться у владельца состояния. Если PostgreSQL владеет persisted state, `NOT NULL` должен быть в schema, даже если Go также валидирует вход.</details>

И:

```sql
type text NOT NULL
```

Я подчёркиваю:

> «Go validation отвечает за хороший API. Database constraints защищают само состояние.»

---

## Card 05 — NULL

Проверяем понимание.

> «`NULL` — это пустая строка?»

<details><summary>Ответ</summary>Нет.</details>

> «Ноль?»

<details><summary>Ответ</summary>Нет.</details>

> «Отсутствие значения?»

<details><summary>Ответ</summary>Да.</details>

Показываю:

```sql
WHERE finished_at = NULL
```

Что здесь неправильно?

<details><summary>Ответ</summary>Сравнение с `NULL` через `=` не даёт ожидаемого результата. Используется `IS NULL` или `IS NOT NULL`.</details>

---

# 6. 00:20–00:30 — Проектируем JobFlow schema

Теперь ученик проектирует.

Говорю:

> «Не даю готовую схему. Сначала вы назовёте, какие поля нужны `jobs`.»

Ожидаем:

```text
id
type
payload
status
created_at
updated_at
```

И `job_attempts`:

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

Вопрос:

> «Какие constraints нам нужны?»

<details><summary>Ответ</summary>`PRIMARY KEY`, `NOT NULL`, `FOREIGN KEY`, а для необходимых бизнес-инвариантов — `CHECK`/`UNIQUE`.</details>

---

# 7. 00:30–00:40 — Плакат 05B: MIGRATION

Открываем `05B`.

Показываю:

```text
001 → 002 → 003
```

Спрашиваю:

> «Зачем нам migrations, если можно просто открыть psql и сделать `CREATE TABLE`?»

<details><summary>Ответ</summary>Чтобы изменение schema было versioned, reproducible, reviewable и одинаково применялось на разных окружениях.</details>

Продолжаю:

> «Schema — состояние базы. Migration — переход между состояниями.»

Это важная формула.

---

# 8. 00:40–00:50 — Goose

Теперь показываю Goose.

Создаём каталог:

```text id="8d7l8s"
migrations/
    00001_init.sql
    00002_job_attempts.sql
    00003_job_indexes.sql
```

И спрашиваю:

> «Что должно содержаться в `00001_init.sql`?»

<details><summary>Ответ</summary>Начальная schema для воспроизводимого создания JobFlow database.</details>

Пример:

```sql id="gky7ni"
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

Почему `Down` должен удалять только то, что создала эта migration?

<details><summary>Ответ</summary>Чтобы каждая migration имела собственную обратную операцию и не разрушала состояние, принадлежащее другим версиям schema.</details>

Запускаем Goose:

```bash id="ii3y6r"
goose -dir migrations postgres "$DATABASE_URL" up
```

Проверяем статус:

```bash id="x11a9c"
goose -dir migrations postgres "$DATABASE_URL" status
```

И обязательно:

> «Теперь удалите database/schema и поднимите её с нуля migrations. Почему это важная проверка?»

<details><summary>Ответ</summary>Проверяется, что schema действительно воспроизводима с чистого состояния, а не случайно зависит от ручных изменений.</details>

---

# 9. 00:50–01:00 — CREATE TABLE / DDL

Теперь ученик пишет SQL.

Задача:

```sql id="6m1k23"
CREATE TABLE job_attempts (
    id uuid PRIMARY KEY,
    job_id uuid NOT NULL REFERENCES jobs(id),
    attempt integer NOT NULL,
    status text NOT NULL,
    started_at timestamptz NOT NULL,
    finished_at timestamptz,
    error text
);
```

Вопрос:

> «Зачем `REFERENCES jobs(id)`?»

<details><summary>Ответ</summary>Чтобы база не позволяла создать attempt для несуществующей Job.</details>

Это отличный момент ещё раз связать:

```text
05A
relationship
 ↓
05B
foreign key
```

---

# 10. 01:00–01:10 — DML: INSERT / SELECT / UPDATE

Теперь пишем настоящий SQL.

### INSERT

```sql
INSERT INTO jobs (
    id, type, payload, status, created_at, updated_at
)
VALUES (
    $1, $2, $3, $4, now(), now()
);
```

Спрашиваю:

> «Почему `$1`, `$2`, `$3`, а не `fmt.Sprintf` с JSON внутри строки?»

<details><summary>Ответ</summary>Параметризованные запросы отделяют SQL от значений и позволяют driver безопасно передавать параметры.</details>

### SELECT

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

### UPDATE

```sql
UPDATE jobs
SET
    status = $2,
    updated_at = now()
WHERE id = $1;
```

Вопрос:

> «Почему в production UPDATE почти всегда должен иметь явный `WHERE`?»

<details><summary>Ответ</summary>Без `WHERE` изменятся все строки. Условие должно выражать именно ту группу состояния, которую мы собираемся изменить.</details>

---

# 11. 01:10–01:18 — JOIN

Теперь:

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

Спрашиваю:

> «Почему JOIN использует `a.job_id = j.id`?»

<details><summary>Ответ</summary>Это выражение отношения foreign key: attempt принадлежит конкретной Job.</details>

> «Что случится, если у Job нет attempt и используется обычный `JOIN`?»

<details><summary>Ответ</summary>Строка Job не попадёт в результат, потому что `JOIN` здесь является inner join.</details>

А если нам нужна Job даже без attempt?

<details><summary>Ответ</summary>`LEFT JOIN`.</details>

---

# 12. 01:18–01:28 — INDEX + EXPLAIN

Теперь спрашиваю:

> «Если Job стало десять миллионов, а мы постоянно делаем `WHERE status = 'running'`, что будем смотреть?»

<details><summary>Ответ</summary>Access pattern и индекс, подходящий для этого запроса.</details>

Пишем:

```sql
CREATE INDEX idx_jobs_status
ON jobs(status);
```

Но я спрашиваю:

> «Это автоматически означает, что PostgreSQL будет использовать индекс?»

<details><summary>Ответ</summary>Нет. Planner выбирает план на основании статистики, стоимости операций и структуры запроса.</details>

Теперь:

```sql
EXPLAIN
SELECT *
FROM jobs
WHERE status = 'running';
```

И:

```sql
EXPLAIN ANALYZE
SELECT *
FROM jobs
WHERE status = 'running';
```

Вопрос:

> «Что мы теперь получили?»

<details><summary>Ответ</summary>Наблюдаемое evidence того, как PostgreSQL планирует и фактически выполняет запрос.</details>

---

# 13. 01:28–01:35 — pgxpool v5

Теперь переносим SQL в Go.

Говорю:

> «SQL у нас уже есть. Теперь нам нужен надёжный lifecycle соединений.»

Показываю:

```go id="b9qj3d"
pool, err := pgxpool.New(ctx, databaseURL)
if err != nil {
	return err
}
defer pool.Close()
```

Вопрос:

> «Зачем pool, а не одно соединение?»

<details><summary>Ответ</summary>Backend имеет множество concurrent requests и операций. Pool управляет набором соединений и позволяет ограничивать database concurrency.</details>

Это очень хорошо связывается с 03C:

```text
Go concurrency bound
        +
DB connection pool bound
```

Спрашиваю:

> «Если у нас 1000 goroutines и pool max connections = 20, сколько операций одновременно может иметь connection к PostgreSQL?»

<details><summary>Ответ</summary>Не более лимита пула, в данном примере 20, хотя goroutines могут ждать свободное соединение.</details>

Это важнейшая production mental model.

---

# 14. 01:35–01:45 — PostgreSQL Repository

Теперь ученик переносит MemoryRepository в PostgreSQL.

Например:

```go id="42d4xm"
type PostgresJobRepository struct {
	pool *pgxpool.Pool
}
```

И:

```go id="v7k2mg"
func (r *PostgresJobRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (Job, error) {
	var job Job

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, type, payload, status, created_at, updated_at
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

Спрашиваю:

> «Что здесь важно передать дальше из HTTP request?»

<details><summary>Ответ</summary>`context.Context`, чтобы cancellation/deadline запроса распространялись на DB operation.</details>

---

# 15. 01:45–01:55 — Плакат 05C: TRANSACTION

Открываем третий плакат.

> «Теперь начинается самая важная часть.»

У нас есть два worker.

Оба хотят изменить одну Job.

Что может пойти не так?

<details><summary>Подсказка</summary>Вспомните lost update из плаката.</details>

<details><summary>Ответ</summary>Оба могут прочитать одно и то же состояние, независимо его изменить и затем записать результат. Один update затрёт эффект другого.</details>

---

# 16. 01:55–02:05 — Lost Update руками

Делаем две SQL sessions:

```text
psql A
psql B
```

В обеих:

```sql
SELECT attempt
FROM jobs
WHERE id = $1;
```

Получаем:

```text
3
3
```

Теперь:

```text
A → 4
B → 4
```

И спрашиваем:

> «Сколько попыток должно было быть после двух increment?»

<details><summary>Ответ</summary>5.</details>

> «Почему получилось 4?»

<details><summary>Ответ</summary>Оба worker исходили из одного старого состояния 3, поэтому второй write потерял эффект первого.</details>

---

# 17. 02:05–02:15 — Понимаем transaction

Теперь:

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

Вопрос:

> «Зачем `FOR UPDATE`?»

<details><summary>Ответ</summary>Он блокирует выбранную строку для конфликтующего изменения, чтобы другая транзакция не могла выполнить несовместимое изменение одновременно.</details>

Вторая transaction теперь ждёт.

И это важная связь:

```text
03A
worker concurrency

      ↓

05C
database concurrency
```

---

# 18. 02:15–02:25 — MVCC

Теперь плакат.

Я говорю:

> «PostgreSQL не решает concurrency просто большим глобальным mutex.»

Что используется?

<details><summary>Ответ</summary>MVCC — Multi-Version Concurrency Control.</details>

Вопрос:

> «Для чего нужны версии/snapshots?»

<details><summary>Ответ</summary>Они позволяют транзакциям видеть согласованное представление данных при параллельной работе, уменьшая необходимость блокировать обычные чтения записью.</details>

И важный принцип:

> **Isolation определяет, какие наблюдения и конфликты допускаются между транзакциями.**

---

# 19. 02:25–02:35 — Isolation Levels

Три уровня, которые нам нужны практически:

```text
READ COMMITTED
REPEATABLE READ
SERIALIZABLE
```

Спрашиваю:

> «Что означает isolation level?»

<details><summary>Ответ</summary>Какие эффекты конкурентного выполнения транзакции она может наблюдать и какие гарантии получает относительно других транзакций.</details>

Не заставляем зубрить академические таблицы.

Но спрашиваем:

> «Какой уровень в PostgreSQL является обычным default?»

<details><summary>Ответ</summary>Read Committed.</details>

И:

> «Что нужно понять про Serializable?»

<details><summary>Ответ</summary>Он стремится обеспечить эффект эквивалентный некоторому последовательному выполнению, поэтому при конфликте транзакция может завершиться serialization failure и потребовать retry.</details>

---

# 20. 02:35–02:45 — ACID Interview Drill

Теперь отдельный interview block.

Говорю:

> «На собеседовании могут спросить: что такое ACID? Наша задача — не читать учебник, а дать точный ответ.»

---

### A — Atomicity

> «Что такое atomicity?»

<details><summary>Ответ</summary>Транзакция выполняется как единая атомарная операция: commit фиксирует её изменения целиком, а rollback не оставляет частично применённого результата транзакции.</details>

### C — Consistency

> «Consistency?»

<details><summary>Ответ</summary>После успешного commit система остаётся в состоянии, удовлетворяющем своим инвариантам и constraints.</details>

### I — Isolation

> «Isolation?»

<details><summary>Ответ</summary>Конкурентные транзакции получают определённые гарантии видимости и взаимодействия, зависящие от выбранного isolation level.</details>

### D — Durability

> «Durability?»

<details><summary>Ответ</summary>После успешного commit изменения сохраняются в пределах durability guarantees конкретной СУБД и её конфигурации.</details>

И важная поправка:

> «ACID не означает "transaction safe вообще". У каждого свойства есть конкретная семантика, особенно у Isolation.»

---

# 21. 02:45–02:55 — Transaction в Go

Теперь делаем transaction через `pgxpool/v5`.

Например:

```go id="v7zvka"
tx, err := pool.Begin(ctx)
if err != nil {
	return err
}
defer tx.Rollback(ctx)
```

Дальше:

```go id="m2tywy"
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

return tx.Commit(ctx)
```

Вопрос:

> «Почему `defer tx.Rollback(ctx)` можно оставить даже перед успешным Commit?»

<details><summary>Ответ</summary>После успешного commit rollback уже ничего не изменит; такой idiom защищает error paths и гарантирует попытку rollback при раннем выходе.</details>

---

# 22. 02:55–03:05 — Main Practice

Теперь ученик получает задачу.

> «У нас есть `Job` со статусом `pending`. Сделайте атомарный переход `pending → running` вместе с созданием `job_attempt`.»

Требования:

```text
transaction
row lock
job state transition
attempt insert
commit
rollback on error
```

Вопрос 1:

> «Что должно быть внутри одной transaction?»

<details><summary>Ответ</summary>Изменение Job и создание Attempt, потому что это одна логическая state transition.</details>

Вопрос 2:

> «Как убедиться, что Job ещё `pending`?»

<details><summary>Ответ</summary>Прочитать/заблокировать строку и проверить текущий status внутри transaction.</details>

Вопрос 3:

> «Что делать, если status уже `running`?»

<details><summary>Ответ</summary>Не выполнять повторный переход; вернуть соответствующий business error/результат согласно контракту.</details>

---

# 23. 03:05–03:15 — Transaction Conflict

Теперь специально создаём две concurrent transactions.

Обе пытаются выполнить:

```text
pending → running
```

Что должна сделать система?

<details><summary>Ответ</summary>Одна транзакция успешно получает изменение согласно выбранной locking/isolation policy, другая должна ждать, увидеть изменившееся состояние или получить конфликт согласно выбранной стратегии.</details>

Теперь:

> «Что нужно доказать тестом?»

<details><summary>Ответ</summary>Что две конкурентные попытки не создают некорректное состояние Job и Attempt.</details>

---

# 24. 03:15–03:25 — Database Constraints как Evidence

Теперь я специально ломаю application code.

Допустим, один путь пытается создать:

```text
job_attempt.job_id = несуществующий job
```

Что должно произойти?

<details><summary>Ответ</summary>PostgreSQL отклонит insert из-за foreign key constraint.</details>

И это важный момент:

> «Вот зачем мы не полагаемся только на Go validation.»

Если два разных сервиса работают с одной БД, каждый может ошибиться.

Database constraint остаётся последней границей состояния.

---

# 25. 03:25–03:35 — Final JobFlow Architecture

Теперь показываем, что выросло за урок:

```text id="7x1o6f"
                  HTTP
                   │
                   ▼
              JobService
                   │
          ┌────────┴────────┐
          ▼                 ▼
   JobRepository      EventPublisher
          │
          ▼
       pgxpool
          │
          ▼
     PostgreSQL
       │     │
       ▼     ▼
     jobs  attempts
```

Migration path:

```text id="j8j8xw"
Go model
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
```

---

# 26. 03:35–03:45 — Финальный Interview Round

### Вопрос

> «Чем transaction отличается от connection?»

<details><summary>Ответ</summary>Connection — ресурс соединения с БД. Transaction — логическая атомарная граница изменения состояния, использующая connection во время выполнения.</details>

### Вопрос

> «Что такое MVCC?»

<details><summary>Ответ</summary>Механизм управления конкурентными версиями данных, позволяющий транзакциям работать с согласованными представлениями данных.</details>

### Вопрос

> «Что делает `FOR UPDATE`?»

<details><summary>Ответ</summary>Берёт row-level lock на выбранные строки для дальнейшего изменения в текущей transaction.</details>

### Вопрос

> «Что такое lost update?»

<details><summary>Ответ</summary>Сценарий, когда два concurrent операции читают одно старое состояние и одна запись затирает результат другой.</details>

### Вопрос

> «Что такое ACID?»

<details><summary>Ответ</summary>Atomicity, Consistency, Isolation, Durability — свойства транзакций, определяющие атомарность, сохранение инвариантов, конкурентную видимость/изоляцию и долговечность commit.</details>

### Вопрос

> «Что делать при serialization failure?»

<details><summary>Ответ</summary>Откатить транзакцию и повторить её по безопасной retry policy.</details>

### Вопрос

> «Почему index не всегда хорошо?»

<details><summary>Ответ</summary>Он имеет стоимость памяти и поддержания при изменениях данных; полезен только если улучшает реальные access patterns.</details>

---

# 27. 03:45–03:55 — Code Review

Я прошу ученика открыть свой repository и показать:

```text
migration
schema
repository
query
transaction
test
```

По каждому спрашиваю:

> «Какой invariant здесь защищается?»

> «Кто владеет состоянием?»

> «Какая failure path?»

> «Что будет при concurrent execution?»

> «Что доказывает test?»

Если ученик отвечает только:

> «Это SQL».

Я возвращаю:

> «Нет. Какое состояние эта SQL-команда изменяет и какой контракт она реализует?»

---

# 28. 03:55–04:00 — Closing

Я заканчиваю:

> «Сегодня JobFlow перестал быть игрушкой в памяти.
>
> Мы описали состояние как relational model.
>
> Мы сделали schema reproducible через migrations.
>
> Мы написали реальный SQL.
>
> Мы увидели, что index — это не магическая оптимизация, а часть конкретного access pattern.
>
> Мы подключили PostgreSQL из Go через connection pool.
>
> А затем самое важное — заставили две concurrent операции встретиться на одной строке и увидели, почему transactions и isolation вообще существуют.
>
> Запомните главный переход.
>
> В Go concurrency мы защищали shared state в памяти.
>
> В PostgreSQL мы теперь защищаем shared state, который живёт дольше процесса и одновременно доступен множеству workers.
>
> Поэтому database concurrency — это продолжение Go concurrency, только граница состояния стала внешней.»

И:

> «На следующем этапе мы будем двигать JobFlow за пределы одного процесса. Появятся сообщения, повторная доставка, idempotency и распределённые failure modes.»

---

# 29. Домашнее задание — JobFlow v3: PostgreSQL

Домашнее задание должно не просто «попробовать SQL», а сделать JobFlow реально зависимым от PostgreSQL.

## 1. Migrations

Создать минимум:

```text id="2r8n1c"
00001_jobs.sql
00002_job_attempts.sql
00003_constraints.sql
00004_indexes.sql
```

### Требования

`jobs`:

```text
id
type
payload
status
created_at
updated_at
```

`job_attempts`:

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

Добавить:

* primary keys;
* foreign key;
* `NOT NULL`;
* необходимые `CHECK`;
* необходимые `UNIQUE`;
* индексы под реальные queries.

---

# 30. Repository

Удалить runtime dependency JobFlow от `MemoryJobRepository` как основной реализации.

Создать:

```go
type PostgresJobRepository struct {
    pool *pgxpool.Pool
}
```

Реализовать:

```text
Add
Get
List
UpdateStatus
CreateAttempt
```

---

# 31. pgxpool

Добавить configuration:

```text
MinConns
MaxConns
MaxConnLifetime
MaxConnIdleTime
```

Числа должны быть осмысленными для local development и тестов.

Проверить:

```text
pool creation
query
scan
close
context cancellation
```

---

# 32. Transactions

Реализовать:

> **Start Job**

Атомарно:

```text
lock Job
 ↓
verify status = pending
 ↓
status → running
 ↓
create attempt
 ↓
commit
```

При любой ошибке:

```text
rollback
```

---

# 33. Concurrency Test

Запустить две concurrent попытки начать одну Job.

Ожидаем:

```text
exactly one successful transition
one valid failure / no duplicate transition
one consistent attempt state
```

Запустить:

```bash
go test ./...
go test -race ./...
```

---

# 34. SQL Practice

Ученик должен самостоятельно написать минимум:

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

и минимум три `EXPLAIN`.

Например:

```sql
EXPLAIN ANALYZE
SELECT ...
```

Не просто получить результат, а объяснить:

> «Почему PostgreSQL выбрал этот plan?»

---

# 35. Acceptance Criteria

## PASS

Ученик:

* поднимает PostgreSQL schema с нуля только migrations;
* понимает структуру `jobs` и `job_attempts`;
* использует constraints;
* пишет SELECT / JOIN / UPDATE самостоятельно;
* создаёт и использует index;
* понимает базовый `EXPLAIN`;
* подключает `pgxpool/v5`;
* использует `context`;
* реализует repository;
* использует transaction;
* понимает lost update;
* умеет применить row lock;
* объясняет ACID;
* понимает MVCC;
* может объяснить isolation level;
* корректно обрабатывает rollback;
* тестирует concurrent state transition.

## Engineering quality

```text
schema reproducible
state owner explicit
invariants enforced
SQL parameterized
transactions bounded
contexts propagated
concurrency tested
failure paths observable
```

---

# 36. Финальная mental model

К концу занятия ученик должен видеть PostgreSQL не как «базу, куда мы сохраняем structs», а как систему управления состоянием:

```text id="bdr1p5"
                    JOBFLOW STATE

                       MODEL
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
                  ┌──────┴──────┐
                  ▼             ▼
               COMMIT        CONFLICT
                                │
                                ▼
                              RETRY
```

Финальная фраза урока:

> **Database engineering — это не хранение данных. Это управление корректным состоянием при изменениях и конкуренции.**
