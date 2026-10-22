# Go Backend Combat — 07. BACKEND ARCHITECTURE

**Version:** 1.0.0
**Owner:** Mentor
**Duration:** 120 минут
**Format:** интерактивный mentor-led lesson
**Project:** JobFlow
**Prerequisite:** `01.LANGUAGE.md` — `06.DISTRIBUTED.md`
**Posters:** `07A_HTTP_API.svg`, `07B_APPLICATION_ARCHITECTURE.svg`, `07C_RUNTIME_CACHE.svg`
**Outcome:** ученик собирает все изученные механизмы в полноценный Go backend и самостоятельно доводит JobFlow до состояния, пригодного для демонстрации и mock interview

---

# 1. Purpose

Сегодня мы перестаём изучать отдельные механизмы.

До этого JobFlow последовательно получил:

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

Теперь собираем всё в единое приложение:

```text
                         CLIENT
                            │
                           HTTP
                            │
                            ▼
                      ┌──────────┐
                      │ HANDLERS │
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

Главная задача урока:

> **Собрать систему так, чтобы transport, policy, state, effects, cache и runtime имели явные границы и lifecycle.**

И важнее всего:

> **К концу урока JobFlow должен стать законченным проектом, а не набором учебных примеров.**

После этого восьмой урок полностью отдаём под **настоящее mock interview уровня BigTech / strong middle+ / senior**.

---

# 2. Learning Contract

После занятия ученик умеет:

### HTTP

* поднять HTTP server;
* определить routes;
* написать handlers;
* декодировать request;
* формировать response;
* маппить application errors в HTTP status;
* использовать request context.

### Architecture

* разделить transport / application / domain / infrastructure;
* определить consumer-side interfaces;
* внедрять dependencies через constructors;
* отделять state от effects;
* объяснить dependency direction.

### Redis

* использовать Redis как cache;
* реализовать cache hit / miss;
* определить TTL;
* понимать stale data;
* не делать Redis source of truth без явного требования;
* корректно деградировать при cache failure.

### Runtime

* собрать composition root;
* связать все компоненты;
* управлять lifecycle;
* сделать graceful shutdown;
* ограничить ресурсы;
* провести smoke test всей системы.

### Engineering

Ученик способен объяснить:

```text
что является state
что является effect
кто владеет resource
куда направлена dependency
что произойдёт при отказе
как это доказано
```

---

# 3. Mentor Flow

Сегодня основной режим:

```text
OBSERVE
 ↓
ASK
 ↓
DESIGN
 ↓
CODE
 ↓
RUN
 ↓
BREAK
 ↓
FIX
 ↓
REVIEW
```

И особое правило:

> **После сегодняшнего урока мы больше не строим учебный код ради отдельного понятия. Каждый новый код должен улучшать финальный JobFlow.**

---

# 4. 00:00–00:08 — Возвращаемся к проекту

Я показываю текущую систему.

Говорю:

> «До этого мы строили её по частям. Теперь я хочу посмотреть на неё глазами пользователя.»

Показываю:

```text
client
  ↓
???
```

И спрашиваю:

> «Как пользователь вообще должен создать Job?»

<details style="display: inline;"><summary>Помощь</summary>Назовите transport protocol и минимальную операцию.</details>

<details style="display: inline;"><summary>Ответ</summary>HTTP request, например `POST /jobs`.</details>

Продолжаю:

> «Хорошо. А потом пользователь хочет узнать состояние Job. Что ему нужно?»

<details><summary>Ответ</summary>`GET /jobs/{id}`.</details>

> «А отменить?»

<details><summary>Ответ</summary>`POST /jobs/{id}/cancel` или другой явно определённый command endpoint.</details>

И говорю:

> «Сегодня мы строим реальный вход в систему.»

---

# 5. 00:08–00:20 — Плакат 07A: HTTP & API

Открываем `07A`.

## Card 1 — REQUEST

Показываю:

```http
POST /jobs
Content-Type: application/json

{
  "type": "email",
  "payload": {}
}
```

Спрашиваю:

> «Что здесь является domain data, а что transport data?»

<details><summary>Ответ</summary>HTTP method, path, headers и JSON body — transport. Domain object Job — внутреннее состояние системы.</details>

> «Почему я не хочу просто передать `http.Request` в domain?»

<details><summary>Ответ</summary>Это создаёт зависимость domain/application от HTTP transport и увеличивает coupling.</details>

---

# 6. 00:20–00:30 — HANDLER

Пишем:

```text
HTTP request
   ↓
decode
   ↓
validate transport
   ↓
application command
   ↓
application result
   ↓
HTTP response
```

Спрашиваю:

> «Что должен знать handler о PostgreSQL?»

<details><summary>Ответ</summary>Ничего.</details>

> «О Kafka?»

<details><summary>Ответ</summary>Ничего.</details>

> «О бизнес-правиле перехода Job из `pending` в `running`?»

<details><summary>Ответ</summary>Handler может инициировать use case, но само бизнес-правило должно находиться в application/domain layer.</details>

И фиксирую:

> **Handler — adapter, а не место для бизнес-логики.**

---

# 7. 00:30–00:40 — ERROR MAPPING

Показываю:

```go
ErrJobNotFound
```

Спрашиваю:

> «Что вернём HTTP-клиенту?»

<details><summary>Ответ</summary>Например, HTTP 404.</details>

> «Должен ли `ErrJobNotFound` знать о коде 404?»

<details><summary>Ответ</summary>Нет. Domain/application error не должен зависеть от HTTP.</details>

Собираем:

```text
ErrJobNotFound
      ↓
HTTP adapter
      ↓
404
```

И:

```text
validation error → 400
not found         → 404
conflict          → 409
internal error    → 500
```

---

# 8. 00:40–00:50 — CONTEXT

Показываем:

```text
HTTP request
    ↓
context.Context
    ↓
JobService
    ↓
PostgreSQL
```

Вопрос:

> «Клиент закрыл соединение. Что должна сделать database operation?»

<details><summary>Ответ</summary>Она должна получить cancellation через context и прекратить ненужную работу, если конкретная операция и driver поддерживают эту семантику.</details>

> «Почему нельзя создать внутри service `context.Background()` и заменить request context?»

<details><summary>Ответ</summary>Это разрывает cancellation/deadline propagation от владельца исходной операции.</details>

---

# 9. 00:50–01:00 — Плакат 07B: APPLICATION ARCHITECTURE

Открываем `07B`.

Говорю:

> «До сих пор мы видели эти слои по отдельности. Теперь надо понять, кто за что владеет.»

---

## TRANSPORT

Вопрос:

> «Что делает transport?»

<details><summary>Ответ</summary>Переводит внешний протокол в application command/result и обратно.</details>

---

## APPLICATION

> «Где находится use case?»

<details><summary>Ответ</summary>В application layer: координация command, validation, state change, effect и нужных ports.</details>

Показываем:

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

---

## DOMAIN

> «Что такое domain в нашем JobFlow?»

<details><summary>Ответ</summary>Job state, допустимые transitions, бизнес-инварианты и поведение, не зависящее от HTTP/DB/Kafka.</details>

---

# 10. 01:00–01:10 — PORTS

Я показываю:

```go
type JobRepository interface {
	Get(...)
	Add(...)
	UpdateStatus(...)
}
```

И:

```go
type EventPublisher interface {
	Publish(...)
}
```

Спрашиваю:

> «Кто определяет этот interface?»

<details><summary>Ответ</summary>Потребитель, исходя из реально необходимого ему поведения.</details>

> «Зачем делать interface на весь Kafka client?»

<details><summary>Ответ</summary>Не нужно, если application использует только одно поведение `Publish`. Большой abstraction surface увеличивает coupling.</details>

---

# 11. 01:10–01:20 — ADAPTERS

Теперь показываем concrete implementations:

```text
PostgresJobRepository
KafkaPublisher
RedisCache
HTTPHandler
KafkaConsumer
```

Вопрос:

> «Кто знает про `pgxpool.Pool`?»

<details><summary>Ответ</summary>Postgres adapter.</details>

> «Кто знает про Kafka client?»

<details><summary>Ответ</summary>Kafka adapter.</details>

> «Кто знает про Redis?»

<details><summary>Ответ</summary>Redis adapter/cache component.</details>

> «А JobService?»

<details><summary>Ответ</summary>Он знает только свои contracts/ports.</details>

---

# 12. 01:20–01:30 — Composition Root

Открываем `07C`, Card 1.

Я говорю:

> «Хорошая архитектура ещё не запустилась. Кто создаст все эти объекты?»

<details><summary>Ответ</summary>Composition root, обычно `main` или отдельный application bootstrap package.</details>

Показываем:

```text
main
 ├── config
 ├── logger
 ├── pg pool
 ├── redis
 ├── kafka
 ├── repositories
 ├── publishers
 ├── services
 ├── handlers
 └── workers
```

И спрашиваю:

> «Кто должен создавать PostgreSQL внутри `JobService`?»

<details><summary>Ответ</summary>Никто. `JobService` получает готовый dependency через constructor.</details>

Это снова закрепляет DIP.

---

# 13. 01:30–01:40 — HTTP Server

Теперь пишем реальный server.

Минимально:

```go
srv := &http.Server{
	Addr:    cfg.HTTPAddr,
	Handler: router,
}
```

Запуск:

```text
start
 ↓
serve
```

И graceful shutdown:

```text
signal
 ↓
stop accepting
 ↓
shutdown
 ↓
wait
```

Вопрос:

> «Почему `ListenAndServe()` — недостаточно для production-like lifecycle?»

<details><summary>Ответ</summary>Нужно уметь принимать сигнал завершения, прекращать новый приём работы и корректно останавливать зависимые resources.</details>

---

# 14. 01:40–01:50 — Redis: зачем он нам вообще?

Я показываю `GET /jobs/{id}`.

Спрашиваю:

> «Почему не делать PostgreSQL каждый раз?»

<details><summary>Помощь</summary>Что если запросов на чтение значительно больше, чем изменений?</details>

<details><summary>Ответ</summary>Cache может снизить latency и нагрузку на PostgreSQL для часто читаемого состояния.</details>

Но сразу:

> «Что является source of truth?»

<details><summary>Ответ</summary>PostgreSQL.</details>

> «А Redis?»

<details><summary>Ответ</summary>Derived/temporary state — cache.</details>

---

# 15. 01:50–02:00 — Cache Hit / Miss

Показываем:

```text
GET /jobs/{id}
        │
        ▼
      Redis
      /   \
   HIT     MISS
    │        │
    │        ▼
    │     PostgreSQL
    │        │
    │        ▼
    └──── Redis
             │
             ▼
          response
```

Спрашиваю:

> «Что произойдёт при Redis miss?»

<details><summary>Ответ</summary>Читаем PostgreSQL и можем заполнить cache.</details>

> «Что произойдёт при Redis failure?»

<details><summary>Ответ</summary>Если Redis — необязательный cache, запрос должен деградировать к PostgreSQL согласно policy.</details>

> «Почему Redis нельзя считать source of truth автоматически?»

<details><summary>Ответ</summary>Потому что cache может быть удалён, протухнуть или быть недоступным; основное состояние должно оставаться в durable store.</details>

---

# 16. 02:00–02:10 — Cache Semantics

Вводим:

```text
TTL
MISS
STALE
INVALIDATION
```

Вопрос:

> «Когда обновлять cache после изменения Job?»

<details><summary>Ответ</summary>По выбранной cache policy: например, удалить/обновить cache после успешного изменения authoritative state.</details>

> «Что проще для начала: write-through, cache-aside или что-то другое?»

<details><summary>Ответ</summary>Для учебного JobFlow достаточно cache-aside: сначала cache lookup, при miss — PostgreSQL, затем запись результата в Redis.</details>

Теперь реализуем:

```text
GetJob
```

с cache-aside.

---

# 17. 02:10–02:20 — Resource Limits

Возвращаемся к 03C.

Вопрос:

> «Теперь у нас есть HTTP, PostgreSQL, Redis, Kafka и workers. Что здесь может закончиться?»

<details><summary>Ответ</summary>Goroutines, DB connections, Redis connections, Kafka consumers, in-flight requests, memory, queue capacity и другие ресурсы.</details>

Я пишу:

```text
bounded workers
bounded DB pool
bounded Redis pool
bounded in-flight work
bounded queue
```

И спрашиваю:

> «Почему важно ограничить DB pool, если HTTP сервер уже ограничивает workers?»

<details><summary>Ответ</summary>Не каждая goroutine обязательно завершает DB operation синхронно, и разные paths могут использовать DB. Кроме того, pool является отдельной ресурсной границей database concurrency.</details>

---

# 18. 02:20–02:35 — Main Integration Sprint

Теперь перестаём читать плакаты.

Говорю:

> «Сейчас у нас 15 минут инженерной сборки. Я даю только интерфейс системы. Реализацию пишете вы.»

Требования:

```text
POST /jobs
GET /jobs/{id}
POST /jobs/{id}/cancel
```

`POST /jobs`:

```text
HTTP
 ↓
JobService
 ↓
PostgreSQL transaction
 ↓
outbox
```

`GET /jobs/{id}`:

```text
HTTP
 ↓
Redis
 ├── HIT → response
 └── MISS → PostgreSQL → Redis → response
```

Worker:

```text
Kafka
 ↓
consumer
 ↓
idempotency
 ↓
Job state transition
 ↓
PostgreSQL
```

Ученик должен сам соединить компоненты.

---

# 19. 02:35–02:45 — Integration Attack

Теперь ломаем систему.

## Attack 1 — Redis down

> «Я выключил Redis. Что должно произойти?»

<details><summary>Ответ</summary>Cache path деградирует к PostgreSQL. Основной state продолжает работать.</details>

## Attack 2 — Kafka down

> «Job создаётся, Kafka недоступна. Теряем Job?»

<details><summary>Ответ</summary>Нет. Job и outbox event уже должны быть зафиксированы в PostgreSQL; relay может повторить публикацию позже.</details>

## Attack 3 — client disconnect

> «Клиент закрыл HTTP connection во время GET.»

<details><summary>Ответ</summary>Context cancellation должен распространяться на internal operations.</details>

## Attack 4 — worker crash after DB effect

> «Что произойдёт?»

<details><summary>Ответ</summary>Kafka event может быть redelivered; state transition/consumer должны быть идемпотентными.</details>

## Attack 5 — PostgreSQL unavailable

> «Что делать?»

<details><summary>Ответ</summary>Запросы, требующие authoritative state, должны корректно получить failure; система не должна выдавать фиктивный успешный результат.</details>

---

# 20. 02:45–03:00 — System Review

Теперь я прошу ученика самостоятельно провести architecture walkthrough.

Говорю:

> «Представьте, что вы завтра выходите на code review нового backend. Проведите меня по системе от HTTP request до конечного effect.»

Жду объяснение.

Подсказки только при необходимости:

> «Где transport?»

<details><summary>Ответ</summary>HTTP handler/router.</details>

> «Где policy?»

<details><summary>Ответ</summary>Application/domain.</details>

> «Где state?»

<details><summary>Ответ</summary>PostgreSQL.</details>

> «Где derived state?»

<details><summary>Ответ</summary>Redis cache.</details>

> «Где asynchronous effect?»

<details><summary>Ответ</summary>Kafka/outbox/consumer.</details>

> «Кто создаёт все зависимости?»

<details><summary>Ответ</summary>Composition root.</details>

> «Кто владеет shutdown?»

<details><summary>Ответ</summary>Application runtime / owner в composition root.</details>

---

# 21. 03:00–03:10 — Smoke Test

Теперь система должна реально работать целиком.

Запускаем:

```text
PostgreSQL
Redis
Kafka / Redpanda
JobFlow
```

Проверяем:

```bash
curl -X POST ...
curl ...
```

Создаём Job.

Потом:

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

Потом:

```text
GET /jobs/{id}
      ↓
Redis MISS
      ↓
PostgreSQL
      ↓
Redis SET
```

Следующий:

```text
GET /jobs/{id}
      ↓
Redis HIT
```

Вот здесь ученик впервые видит **весь курс как одну работающую систему**.

---

# 22. 03:10–03:20 — Финальная диагностика

Я задаю вопросы уже не по отдельным темам, а cross-system.

### Вопрос

> «Если Redis быстрее PostgreSQL, почему мы не сделаем Redis source of truth?»

<details><summary>Ответ</summary>Потому что это изменило бы durability/consistency model и потребовало бы отдельной стратегии восстановления и сохранения состояния. Текущий контракт системы предполагает PostgreSQL как durable source of truth.</details>

### Вопрос

> «Если Kafka доставляет сообщение дважды, где это безопасно?»

<details><summary>Ответ</summary>В consumer idempotency/state transition logic.</details>

### Вопрос

> «Если PostgreSQL commit прошёл, а Kafka publish нет?»

<details><summary>Ответ</summary>Outbox сохраняет message intent; relay повторит publish.</details>

### Вопрос

> «Почему handler не вызывает Kafka client напрямую?»

<details><summary>Ответ</summary>Это смешивает transport и infrastructure detail, увеличивает coupling и нарушает application boundary.</details>

### Вопрос

> «Почему application не создаёт Redis client?»

<details><summary>Ответ</summary>Dependencies создаются в composition root и внедряются через contracts.</details>

---

# 23. 03:20–03:30 — Финальный Architecture Interview

Теперь симулирую один короткий system-design вопрос.

> «У нас 1000 запросов GET `/jobs/{id}` в секунду. PostgreSQL начинает перегружаться. Что будете делать?»

Не даю подсказки сразу.

Жду.

<details><summary>Помощь</summary>Посмотрите на 07C: cache и limits.</details>

<details><summary>Ответ</summary>Добавлю cache-aside через Redis для горячих reads, определю TTL/invalidation policy, ограничу database concurrency pool и проверю фактический access pattern. После этого измерю latency, cache hit rate и DB load.</details>

Следующий:

> «Redis падает.»

<details><summary>Ответ</summary>Должен быть cache miss/failure path к PostgreSQL, если Redis не является authoritative state.</details>

Следующий:

> «Теперь Kafka consumer отстаёт.»

<details><summary>Ответ</summary>Смотрю partition count, consumer group concurrency, processing latency, backlog/lag, downstream limits и backpressure. Не просто увеличиваю goroutines бесконечно.</details>

---

# 24. 03:30–03:40 — Подготовка к финальному состоянию проекта

Я говорю:

> «Сегодня мы закончили архитектурную часть. Больше новых крупных технологий до interview не вводим.»

Теперь цель:

```text
JobFlow
    ↓
complete
    ↓
test
    ↓
verify
    ↓
explain
```

И проверяем repository tree.

Ожидаем примерно:

```text
cmd/
  jobflow/
internal/
  domain/
  application/
  transport/
  repository/
  publisher/
  consumer/
  cache/
  config/
migrations/
tests/
```

Важно:

> **Структура может отличаться. Важны ownership и dependency boundaries, а не конкретные названия каталогов.**

---

# 25. 03:40–03:50 — Final Acceptance Review

Ученик должен самостоятельно показать:

### API

```text
POST /jobs
GET /jobs/{id}
POST /jobs/{id}/cancel
```

### Persistence

```text
PostgreSQL
migrations
transactions
constraints
```

### Messaging

```text
outbox
Kafka / Redpanda
consumer group
idempotency
```

### Cache

```text
Redis
cache-aside
TTL
failure fallback
```

### Concurrency

```text
worker pool
bounded concurrency
context
graceful shutdown
race-free
```

### Architecture

```text
transport
application
domain
ports
adapters
composition root
```

---

# 26. 03:50–04:00 — Final Lesson Closing

Я заканчиваю урок так:

> «Посмотрите, что произошло с нашим проектом за курс.
>
> Мы начали с `struct Job`.
>
> Потом научили Go выполнять работу конкурентно.
>
> Затем научились проводить границы между ответственностями.
>
> Потом сделали состояние durable в PostgreSQL.
>
> Научились защищать его transaction semantics.
>
> Перешли через network boundary и научились переживать duplicate, retry и partial failure.
>
> А сегодня собрали всё это в один backend, который действительно можно запустить и использовать через HTTP.
>
> Это уже не набор учебных примеров.
>
> У нас есть приложение:
>
> HTTP transport.
>
> Application policy.
>
> PostgreSQL state.
>
> Redis derived state.
>
> Kafka effects.
>
> Workers.
>
> Outbox.
>
> Idempotency.
>
> Cancellation.
>
> Graceful shutdown.
>
> И главное — мы можем объяснить, почему каждая граница существует.»

После паузы:

> «Восьмое занятие будет другим.
>
> Никаких новых тем.
>
> Я перестану быть учителем и стану интервьюером.
>
> Мы проведём полноценное mock interview: Go, concurrency, SQL, PostgreSQL, transactions, Kafka, distributed systems, architecture, debugging и system design.
>
> Я буду спрашивать не только определения. Я буду давать failure scenarios, code review и архитектурные задачи.
>
> И вы должны будете не вспомнить правильный термин, а **спроектировать и защитить решение**.»

---

# 27. Домашнее задание — JobFlow FINAL

Это главное домашнее задание всего курса.

На следующем занятии **новый функциональный код уже не добавляем**. Проект должен быть закончен самостоятельно.

## 1. Complete HTTP API

Завершить:

```text
POST /jobs
GET /jobs/{id}
POST /jobs/{id}/cancel
GET /healthz
GET /readyz
```

---

## 2. PostgreSQL

Проверить:

```text
migrations from clean database
constraints
indexes
transactions
concurrent transitions
```

Запустить migrations на чистой БД.

---

## 3. Redis

Довести:

```text
cache-aside
TTL
cache miss
cache failure fallback
```

Проверить:

```text
Redis available
Redis unavailable
stale / expired entry
```

---

## 4. Kafka / Redpanda

Проверить:

```text
producer
outbox relay
consumer group
partition key
idempotency
retry
backoff
graceful shutdown
```

---

## 5. Failure Scenarios

Самостоятельно проверить:

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

Для каждого:

```text
expected behavior
observed behavior
evidence
```

---

# 28. Quality Gate

Перед 08 ученику нужно получить:

```bash
go test ./...
go test -race ./...
go vet ./...
```

и проверить чистую сборку.

Дополнительно:

```text
application starts
migrations apply from zero
HTTP works
Kafka flow works
Redis cache works
graceful shutdown works
```

---

# 29. Final Deliverable

К восьмому занятию JobFlow должен быть репозиторием, который можно открыть перед интервьюером.

Минимальный сценарий демонстрации:

```text
1. start infrastructure
2. run migrations
3. start application
4. POST /jobs
5. observe PostgreSQL state
6. observe outbox
7. observe Kafka event
8. observe worker processing
9. GET /jobs/{id}
10. observe Redis hit
11. inject failure
12. explain recovery
```

И ученик должен уметь сказать:

> **«Я могу показать полный путь запроса, состояние системы, асинхронный effect, cache path, failure mode и механизм восстановления.»**

Это и есть завершение практической части курса.

---

# 30. Acceptance Criteria

## PASS

Ученик:

* самостоятельно запускает JobFlow;
* объясняет request → application → state/effect;
* объясняет все три плаката;
* пишет и меняет handlers без подсказки;
* понимает application boundary;
* умеет найти dependency direction;
* использует Redis как cache;
* знает source of truth;
* объясняет outbox;
* объясняет idempotency;
* объясняет Kafka delivery semantics;
* объясняет transaction boundary;
* может провести failure analysis всей системы;
* умеет выполнить graceful shutdown;
* может показать evidence работы системы.

## Финальный уровень

Ученик должен отвечать не:

> «Мы сделали Redis, Kafka и PostgreSQL.»

а:

> **«Для этого требования выбрана такая граница, потому что она сохраняет такой invariant; Redis является derived state, PostgreSQL — source of truth, Kafka — asynchronous effect, duplicate обрабатывается через idempotency, а recovery ограничен такой policy.»**

---

# 31. Финальная mental model курса перед интервью

К этому моменту вся серия плакатов складывается в одну систему:

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
        ↓
08 MOCK INTERVIEW
```

А JobFlow:

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
                  /       \
                 ▼         ▼
              DOMAIN      EFFECT
                 │           │
                 ▼           ▼
            PostgreSQL      Kafka
                 │           │
                 │        Workers
                 │           │
                 └──────┬────┘
                        │
                     Redis
                     cache
```

Финальный принцип перед mock interview:

> **Не доказывай, что знаешь технологии. Доказывай, что умеешь управлять состоянием, зависимостями, ресурсами, отказами и изменениями.**
