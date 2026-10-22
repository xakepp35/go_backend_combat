# Go Backend Combat — 07. BACKEND ARCHITECTURE

**Version:** 1.0.0
**Topic:** Backend Architecture
**Format:** Дидактический материал ученика — HTTP, application architecture, Redis, composition root и финальная сборка JobFlow
**Project:** JobFlow
**Duration:** 120 минут
**Prerequisite:** `01_LANGUAGE.md` — `06_DISTRIBUTED_SYSTEMS`

Плакаты:

* [07A — HTTP API](../posters/07A_HTTP_API.svg)
* [07B — Application Architecture](../posters/07B_APPLICATION_ARCHITECTURE.svg)
* [07C — Runtime & Cache](../posters/07C_RUNTIME_CACHE.svg)

---

# 1. Цель

JobFlow больше не изучается по отдельным механизмам.

К этому моменту есть:

```text
Go
 ↓
Concurrency
 ↓
Design
 ↓
PostgreSQL
 ↓
Kafka / Redpanda
```

Теперь они собираются в единый backend:

```text
                         CLIENT
                            │
                           HTTP
                            │
                            ▼
                        HANDLER
                            │
                            ▼
                      APPLICATION
                            │
             ┌──────────────┼──────────────┐
             ▼              ▼              ▼
        PostgreSQL        Redis          Kafka
             ▲                             │
             │                             ▼
             └────────────────────────── Workers
```

Главная задача:

> **Transport, policy, state, effects, cache и runtime должны иметь явные boundaries и lifecycle.**

---

# 2. Финальная модель JobFlow

К концу занятия должен существовать рабочий end-to-end flow:

```text
POST /jobs
    ↓
HTTP handler
    ↓
JobService
    ↓
PostgreSQL transaction
    ↓
Outbox
    ↓
Kafka / Redpanda
    ↓
Consumer
    ↓
Worker
    ↓
Idempotent state transition
    ↓
PostgreSQL
```

Чтение:

```text
GET /jobs/{id}
    ↓
Redis
 ┌──┴──┐
HIT   MISS
 │      │
 │      ▼
 │   PostgreSQL
 │      │
 └──────┴──────→ response
```

---

# 3. 07A — HTTP API

![07A — HTTP API](../posters/07A_HTTP_API.svg)

HTTP — transport boundary.

Его задача:

```text
HTTP request
    ↓
transport parsing
    ↓
application command
    ↓
application result
    ↓
HTTP response
```

---

# 4. Transport и Domain Data

Пример:

```http
POST /jobs
Content-Type: application/json

{
  "type": "email",
  "payload": {}
}
```

Transport data:

```text
method
path
headers
JSON
status code
```

Domain data:

```text
Job
JobStatus
JobID
business state
```

Не следует передавать `http.Request` внутрь domain/application policy без необходимости.

Иначе application начинает зависеть от HTTP transport.

---

# 5. Handler

Handler должен заниматься transport concerns:

```text
decode
validate transport input
call use case
map result
write response
```

Handler не должен знать:

```text
PostgreSQL implementation
Kafka client
Redis implementation
business state machine internals
```

Принцип:

> **Handler — adapter, а не место для бизнес-логики.**

---

# 6. API JobFlow

Минимальный внешний API:

```text
POST /jobs
GET /jobs/{id}
POST /jobs/{id}/cancel
GET /healthz
GET /readyz
```

### `POST /jobs`

Создаёт Job.

### `GET /jobs/{id}`

Возвращает текущее состояние Job.

### `POST /jobs/{id}/cancel`

Запрашивает переход Job в отменённое состояние согласно domain contract.

---

# 7. HTTP Error Mapping

Application error:

```go
ErrJobNotFound
```

не должен знать HTTP.

Mapping происходит в transport adapter:

```text
ErrJobNotFound
       ↓
HTTP handler
       ↓
404 Not Found
```

Пример mapping:

```text
validation error → 400
not found        → 404
conflict         → 409
internal failure → 500
```

Главный принцип:

> **Application semantics не зависят от HTTP status codes.**

---

# 8. Context Propagation

Типичный путь:

```text
HTTP request
     ↓
context.Context
     ↓
JobService
     ↓
PostgreSQL
```

Context позволяет передавать:

```text
cancellation
deadline
request-scoped metadata
```

Плохая практика:

```go
ctx := context.Background()
```

внутри application path, когда уже существует request context.

Это разрывает propagation cancellation/deadline.

---

# 9. Cancellation

Если клиент закрыл соединение:

```text
client disconnect
      ↓
request context cancelled
      ↓
service
      ↓
database / downstream
```

Зависимые операции должны иметь возможность остановить ненужную работу.

Context — это:

> **сигнал cancellation, а не механизм принудительного убийства goroutine.**

---

# 10. 07B — Application Architecture

![07B — Application Architecture](../posters/07B_APPLICATION_ARCHITECTURE.svg)

Основные границы:

```text
TRANSPORT
    ↓
APPLICATION
    ↓
DOMAIN / POLICY
    ↓
PORTS
    ↑
ADAPTERS
```

---

# 11. Transport

Transport переводит внешний protocol в application command и результат обратно.

Например:

```text
HTTP JSON
    ↓
CreateJobCommand
    ↓
JobService
```

и:

```text
JobResult
    ↓
HTTP response
```

Transport не владеет business policy.

---

# 12. Application Layer

Application layer координирует use case.

Упрощённый flow:

```text
command
 ↓
validate
 ↓
decide
 ↓
change state
 ↓
emit effect
```

Например:

```text
CreateJob
 ↓
validate input
 ↓
create Job
 ↓
persist Job
 ↓
record outbox event
 ↓
return result
```

---

# 13. Domain

Domain содержит то, что относится к business semantics JobFlow:

```text
Job
JobStatus
allowed transitions
business invariants
domain behavior
```

Domain не должен знать:

```text
HTTP
SQL driver
Kafka client
Redis client
```

---

# 14. Ports

Application зависит от behavior contracts.

Например:

```go
type JobRepository interface {
	Get(ctx context.Context, id string) (Job, error)
	Add(ctx context.Context, job Job) error
	UpdateStatus(ctx context.Context, id string, status JobStatus) error
}
```

И:

```go
type EventPublisher interface {
	Publish(ctx context.Context, event JobCreated) error
}
```

Interface описывает именно то, что нужно consumer.

Не нужно описывать весь внешний client API.

---

# 15. Adapters

Concrete infrastructure:

```text
PostgresJobRepository
KafkaPublisher
RedisCache
HTTPHandler
KafkaConsumer
```

Их задача — адаптировать конкретную technology к application contracts.

Например:

```text
JobRepository
     ↑
PostgresJobRepository
     │
     ▼
pgxpool
```

или:

```text
EventPublisher
     ↑
KafkaPublisher
     │
     ▼
Kafka client
```

---

# 16. Dependency Direction

Целевая зависимость:

```text
Application
     ↓
Contract
     ↑
Infrastructure
```

Не:

```text
Application
     ↓
PostgreSQL driver
```

и не:

```text
Application
     ↓
Kafka client
```

Главный вопрос:

> **Может ли infrastructure implementation измениться без переписывания policy?**

Чем лучше boundary, тем меньше blast radius.

---

# 17. Composition Root

Кто создаёт объекты?

Ответ:

> **Composition root.**

Обычно это `main` или bootstrap package.

Пример:

```text
main
 ├── config
 ├── logger
 ├── PostgreSQL pool
 ├── Redis client
 ├── Kafka client
 ├── repositories
 ├── publishers
 ├── services
 ├── handlers
 └── workers
```

`JobService` не должен сам создавать:

```text
PostgreSQL connection
Redis client
Kafka client
```

Он получает зависимости через constructor.

---

# 18. Constructor Injection

Например:

```go
type JobService struct {
	repo      JobRepository
	publisher EventPublisher
}
```

```go
func NewJobService(
	repo JobRepository,
	publisher EventPublisher,
) *JobService {
	return &JobService{
		repo:      repo,
		publisher: publisher,
	}
}
```

Таким образом dependency становится:

```text
explicit
observable
replaceable
testable
```

---

# 19. HTTP Server Lifecycle

Минимальная модель:

```go
srv := &http.Server{
	Addr:    cfg.HTTPAddr,
	Handler: router,
}
```

Lifecycle:

```text
start
  ↓
serve
  ↓
SIGTERM
  ↓
stop accepting
  ↓
shutdown
  ↓
wait
  ↓
exit
```

Production-like server должен иметь явный shutdown path.

---

# 20. 07C — Runtime & Cache

![07C — Runtime & Cache](../posters/07C_RUNTIME_CACHE.svg)

Теперь у системы появляется ещё одна resource/runtime boundary:

```text
HTTP
PostgreSQL
Redis
Kafka
Workers
```

Каждый из этих компонентов имеет:

```text
resource limits
lifecycle
failure mode
shutdown semantics
```

---

# 21. Redis: зачем он нужен?

Основное чтение:

```text
GET /jobs/{id}
```

Если PostgreSQL становится bottleneck для большого количества reads, cache может уменьшить нагрузку.

Но важно различать:

```text
PostgreSQL → source of truth
Redis      → derived / temporary state
```

Redis не становится authoritative state только потому, что он быстрее.

---

# 22. Cache-Aside

Базовая схема:

```text
GET /jobs/{id}
        │
        ▼
      Redis
      /   \
    HIT   MISS
    │       │
    │       ▼
    │    PostgreSQL
    │       │
    │       ▼
    └──── Redis
            │
            ▼
         response
```

### HIT

Вернуть cached value.

### MISS

Получить authoritative state из PostgreSQL.

После этого можно заполнить cache.

---

# 23. Cache Failure

Если Redis недоступен:

```text
GET
 ↓
Redis failure
 ↓
PostgreSQL
 ↓
response
```

если cache является optional optimization.

Главная идея:

> **Cache failure не должен автоматически означать loss of authoritative state.**

Но fallback policy должна учитывать нагрузку: массовый Redis outage может создать cache stampede и перегрузить PostgreSQL.

---

# 24. TTL

Cache entry имеет ограниченный lifetime:

```text
value
  ↓
TTL
  ↓
expired
  ↓
MISS
```

TTL управляет тем, как долго cached representation считается пригодной.

Это не заменяет consistency policy.

---

# 25. Stale Data

Redis может содержать значение, которое уже отличается от PostgreSQL.

Например:

```text
PostgreSQL
status = completed

Redis
status = running
```

Нужно определить policy:

```text
invalidate
update
short TTL
version check
```

Для базового JobFlow достаточно:

> **PostgreSQL — source of truth, Redis — cache-aside с TTL.**

---

# 26. Cache Invalidation

После изменения authoritative state:

```text
transaction
 ↓
PostgreSQL updated
```

cache должен быть:

```text
updated
```

или:

```text
invalidated
```

Простой подход:

```text
write DB
 ↓
invalidate cache
```

Но если cache operation fail, database state всё равно остаётся authoritative.

---

# 27. Resource Limits

В финальном JobFlow должны быть ограничены:

```text
worker count
DB connections
Redis connections
in-flight work
queue capacity
HTTP request resources
```

Например:

```text
workers = 16
DB MaxConns = 20
queue capacity = 1000
```

Конкретные значения зависят от workload и должны быть измеряемыми.

---

# 28. Почему DB pool — отдельная граница

Даже если:

```text
workers = 100
```

и:

```text
DB MaxConns = 20
```

database concurrency остаётся ограниченной pool capacity.

То есть:

```text
goroutines
     ↓
DB pool
     ↓
PostgreSQL
```

является отдельной resource boundary.

---

# 29. Финальная интеграция

Job creation:

```text
HTTP
 ↓
JobService
 ↓
PostgreSQL transaction
 ├── jobs
 └── outbox
 ↓
COMMIT
```

Asynchronous processing:

```text
outbox
 ↓
relay
 ↓
Kafka
 ↓
consumer
 ↓
worker pool
 ↓
idempotent transition
 ↓
PostgreSQL
```

Read:

```text
HTTP
 ↓
Redis
 ├── HIT
 └── MISS → PostgreSQL → Redis
```

---

# 30. End-to-End Flow

Полная система:

```text
                         CLIENT
                            │
                           HTTP
                            │
                            ▼
                      ┌──────────┐
                      │ HANDLER  │
                      └────┬─────┘
                           ▼
                    ┌─────────────┐
                    │ APPLICATION │
                    │    POLICY   │
                    └─────┬───────┘
                          │
             ┌────────────┼────────────┐
             ▼            ▼            ▼
        PostgreSQL      Redis        Kafka
             ▲                         │
             │                         ▼
             └──────────────────── Workers
```

---

# 31. Smoke Test

Полный минимальный сценарий:

```text
POST /jobs
    ↓
PostgreSQL ✓
    ↓
outbox ✓
    ↓
Kafka ✓
    ↓
worker ✓
    ↓
PostgreSQL state ✓
```

Затем:

```text
GET /jobs/{id}
    ↓
Redis MISS
    ↓
PostgreSQL
    ↓
Redis SET
```

И повторный запрос:

```text
GET /jobs/{id}
    ↓
Redis HIT
```

Это первый момент курса, когда вся архитектура должна работать как единая система.

---

# 32. Failure Matrix

| Failure                     | Expected behavior                        |
| --------------------------- | ---------------------------------------- |
| Redis unavailable           | fallback to PostgreSQL                   |
| Kafka unavailable           | event remains in Outbox                  |
| consumer crash after effect | redelivery + idempotency                 |
| PostgreSQL unavailable      | authoritative operation fails explicitly |
| HTTP cancellation           | context propagates downstream            |
| SIGTERM                     | graceful shutdown                        |
| cache stale                 | invalidation / TTL policy                |

---

# 33. Практическое задание

Довести JobFlow до end-to-end состояния.

Не добавлять новую технологию.

Нужно закончить существующие компоненты:

```text
HTTP
Application
PostgreSQL
Outbox
Kafka
Consumer
WorkerPool
Redis
Shutdown
```

---

# 34. API Implementation

Реализовать:

```text
POST /jobs
GET /jobs/{id}
POST /jobs/{id}/cancel
GET /healthz
GET /readyz
```

Проверить:

```text
valid request
invalid request
not found
conflict
internal failure
```

---

# 35. Health vs Readiness

Разделять:

### `/healthz`

> Процесс жив.

### `/readyz`

> Процесс готов обслуживать traffic с учётом обязательных dependencies.

Это не обязательно означает:

```text
all optional components healthy
```

Readiness policy должна соответствовать тому, какие зависимости критичны для конкретной операции.

---

# 36. Redis Practice

Реализовать cache-aside для:

```text
GET /jobs/{id}
```

Проверить:

```text
cache hit
cache miss
TTL expiration
Redis unavailable
```

Критерий:

> При обычном Redis failure authoritative state продолжает читаться из PostgreSQL, если capacity системы это позволяет.

---

# 37. Composition Root Practice

Собрать все dependencies в одном месте:

```text
config
 ↓
clients
 ↓
adapters
 ↓
services
 ↓
handlers
 ↓
workers
 ↓
server
```

Не создавать infrastructure clients внутри domain/application components.

---

# 38. Runtime Lifecycle Practice

Определить:

```text
startup order
shutdown order
ownership
cancellation
wait
resource release
```

Например:

```text
startup
 ↓
config
 ↓
DB
 ↓
Redis
 ↓
Kafka
 ↓
services
 ↓
HTTP
```

Shutdown должен учитывать dependency graph.

---

# 39. Final Architecture Review

Пройди систему сверху вниз.

### Transport

```text
HTTP → handler
```

### Application

```text
handler → use case
```

### Domain

```text
use case → Job state/invariants
```

### State

```text
PostgreSQL
```

### Cache

```text
Redis
```

### Effect

```text
Outbox → Kafka
```

### Async execution

```text
Kafka → workers
```

### Runtime

```text
composition root → lifecycle
```

---

# 40. Cross-System Questions

### Где находится authoritative state?

<details><summary>Ответ</summary>PostgreSQL.</details>

### Где находится derived state?

<details><summary>Ответ</summary>Redis cache.</details>

### Где находится asynchronous effect?

<details><summary>Ответ</summary>Kafka / Redpanda delivery pipeline.</details>

### Где application policy?

<details><summary>Ответ</summary>Application/domain layers.</details>

### Кто создаёт infrastructure?

<details><summary>Ответ</summary>Composition root.</details>

### Что происходит при Redis failure?

<details><summary>Ответ</summary>Cache path деградирует к PostgreSQL при допустимой нагрузке.</details>

### Что происходит при Kafka failure?

<details><summary>Ответ</summary>Outbox сохраняет event intent; relay повторяет delivery позже.</details>

### Что происходит при consumer crash после effect?

<details><summary>Ответ</summary>Возможна redelivery; idempotency не допускает второй business effect.</details>

---

# 41. System Design Exercise

Сценарий:

> `GET /jobs/{id}` достигает 1000 requests/sec. PostgreSQL перегружен.

Определи:

```text
1. Где bottleneck?
2. Как измерить?
3. Где добавить cache?
4. Какой TTL?
5. Что делать при Redis failure?
6. Как избежать cache stampede?
7. Как ограничить DB concurrency?
```

<details><summary>Ожидаемая линия решения</summary>

Cache-aside через Redis, TTL/invalidation policy, DB pool limit, измерение cache hit rate/latency/DB load. При Redis failure — fallback к PostgreSQL в рамках допустимой capacity. При массовом miss необходимо учитывать stampede и применять ограничение/коалесинг запросов при необходимости.

</details>

---

# 42. Checkpoint

### Что должен делать HTTP handler?

<details><summary>Ответ</summary>Transport mapping: parse request, вызвать application use case, преобразовать result/error в HTTP response.</details>

### Должен ли domain знать HTTP?

<details><summary>Ответ</summary>Нет.</details>

### Должен ли application знать PostgreSQL client?

<details><summary>Ответ</summary>Нет, application должен зависеть от необходимого contract.</details>

### Кто создаёт dependencies?

<details><summary>Ответ</summary>Composition root.</details>

### Что такое source of truth?

<details><summary>Ответ</summary>Авторитетный владелец состояния, от которого восстанавливается корректное значение.</details>

### Что такое cache-aside?

<details><summary>Ответ</summary>Сначала проверяется cache, при miss читается source of truth и результат помещается в cache.</details>

### Что делать при Redis failure?

<details><summary>Ответ</summary>Использовать fallback policy к authoritative store, если это допустимо по resource limits.</details>

### Почему context нельзя заменять `context.Background()`?

<details><summary>Ответ</summary>Это разрывает cancellation и deadline propagation исходной операции.</details>

### Что такое composition root?

<details><summary>Ответ</summary>Место, где создаются concrete dependencies и собирается runtime graph приложения.</details>

### Как проверить архитектуру?

<details><summary>Ответ</summary>Изменениями и failure scenarios: storage replacement, transport replacement, Redis failure, Kafka failure, cancellation, shutdown и end-to-end smoke test.</details>

---

# 43. Final Quality Gate

Перед mock interview должны проходить:

```bash
go test ./...
go test -race ./...
go vet ./...
```

Также:

```text
clean migrations
application startup
HTTP smoke test
PostgreSQL persistence
Outbox publish
Kafka consume
worker processing
Redis cache hit
Redis fallback
graceful shutdown
```

---

# 44. Финальная демонстрация

Минимальный сценарий:

```text
1. Start infrastructure
2. Run migrations
3. Start JobFlow
4. POST /jobs
5. Check PostgreSQL
6. Check outbox
7. Observe Kafka event
8. Observe worker processing
9. GET /jobs/{id}
10. Observe Redis MISS
11. GET /jobs/{id}
12. Observe Redis HIT
13. Disable Redis
14. Verify fallback
15. Disable Kafka
16. Verify outbox recovery
17. Shutdown gracefully
```

Ученик должен уметь объяснить каждый переход.

---

# 45. Домашнее задание — JobFlow FINAL

Это последнее coding homework перед mock interview.

## Завершить HTTP API

```text
POST /jobs
GET /jobs/{id}
POST /jobs/{id}/cancel
GET /healthz
GET /readyz
```

## Проверить PostgreSQL

```text
migrations from zero
constraints
indexes
transactions
concurrent state transitions
```

## Проверить Redis

```text
cache-aside
TTL
MISS
HIT
failure fallback
```

## Проверить Kafka

```text
outbox
relay
consumer group
partition key
idempotency
retry
backoff
graceful shutdown
```

## Проверить failures

```text
PostgreSQL down
Redis down
Kafka down
consumer crash
relay crash
duplicate event
HTTP cancellation
worker shutdown
```

Для каждого записать:

```text
expected behavior
observed behavior
evidence
```

---

# 46. Acceptance Criteria

## PASS

JobFlow:

```text
starts from clean state
migrations are reproducible
HTTP works
PostgreSQL is authoritative
Redis is derived cache
Kafka flow works
Outbox is atomic with local state
consumer is idempotent
worker concurrency is bounded
context propagates
shutdown is graceful
tests pass
race detector passes
```

Ученик может пройти путь:

```text
HTTP
 ↓
application
 ↓
state
 ↓
effect
 ↓
worker
 ↓
state update
```

и объяснить failure mode каждого boundary.

---

# 47. Финальная mental model

Все семь практических этапов курса:

```text
01 ENGINEERING MINDSET
        ↓
02 GO LANGUAGE
        ↓
03 CONCURRENCY
        ↓
04 DESIGN
        ↓
05 DATA
        ↓
06 DISTRIBUTED
        ↓
07 ARCHITECTURE
```

собираются в одну систему:

```text
                         CLIENT
                            │
                           HTTP
                            │
                            ▼
                       TRANSPORT
                            │
                            ▼
                       APPLICATION
                       /         \
                      ▼           ▼
                  DOMAIN       EFFECT
                     │            │
                     ▼            ▼
                PostgreSQL      Kafka
                     │            │
                     │         Workers
                     │            │
                     └─────┬──────┘
                           │
                        Redis
                         cache
```

Главный принцип перед интервью:

> **Не доказывай, что знаешь технологии. Доказывай, что умеешь управлять состоянием, зависимостями, ресурсами, отказами и изменениями.**
