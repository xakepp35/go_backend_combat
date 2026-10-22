# Go Backend Combat — 06. DISTRIBUTED SYSTEMS

**Version:** 1.0.0
**Owner:** Mentor
**Duration:** 120 минут
**Format:** Урок 6 — распределённые системы через Kafka/Redpanda, delivery semantics, idempotency и Outbox на JobFlow
**Project:** JobFlow
**Prerequisite:** `05.DATA.md`
**Posters:** `06A_DISTRIBUTED_MODEL.svg`, `06B_DELIVERY_MESSAGING.svg`, `06C_FAILURE_EFFECTS.svg`
**Outcome:** ученик понимает семантику распределённого взаимодействия и умеет проектировать Kafka/Redpanda consumer с at-least-once delivery, idempotency, retry и bounded recovery

---

# 1. Purpose

До этого момента JobFlow имел такую границу:

```text
HTTP
 ↓
JobService
 ↓
PostgreSQL
```

Теперь мы делаем принципиально другой шаг.

После изменения состояния Job нам нужно сообщить об этом другой части системы:

```text
PostgreSQL
    │
    │ event
    ▼
Kafka / Redpanda
    │
    ▼
Workers
    │
    ▼
PostgreSQL
```

И здесь появляется класс проблем, которых нет внутри одной функции:

```text
delay
loss
duplicate
reorder
partial failure
unknown result
```

Главная идея занятия:

> **Распределённая система не обещает, что операция произойдёт ровно один раз. Она должна иметь явную семантику доставки, повторов и внешнего эффекта.**

---

# 2. Learning Contract

После занятия ученик умеет:

### Distributed model

* объяснить partial failure;
* объяснить, почему timeout не означает failure;
* объяснить duplicate delivery;
* понимать ordering scope;
* различать local и distributed failure.

### Messaging

* объяснить producer / broker / consumer;
* понимать topic / partition;
* понимать consumer group;
* понимать offset;
* объяснить at-most-once;
* объяснить at-least-once;
* аккуратно объяснить exactly-once;
* понимать связь partition и parallelism.

### Reliability

* объяснить idempotency;
* определить idempotency key;
* различать retryable / non-retryable error;
* использовать exponential backoff;
* понимать jitter;
* объяснить outbox pattern;
* понимать почему DB commit и Kafka publish нельзя считать одной локальной transaction.

### JobFlow

Ученик вместе с ментором:

```text
PostgreSQL
   ↓
Outbox
   ↓
Kafka / Redpanda
   ↓
Consumer
   ↓
Worker Pool
   ↓
Idempotent state transition
```

и понимает, **где возможен duplicate и как он обезвреживается**.

---

# 3. Mentor Flow

На этом уроке основной цикл:

```text
OBSERVE
 ↓
ASK
 ↓
PREDICT FAILURE
 ↓
DESIGN SEMANTICS
 ↓
IMPLEMENT
 ↓
INJECT FAILURE
 ↓
OBSERVE
 ↓
FIX
 ↓
VERIFY
```

Особое правило:

> **Я не даю ученику готовую distributed architecture в начале. Мы сначала создаём failure, а потом строим механизм, который делает этот failure безопасным.**

---

# 4. 00:00–00:10 — Начало: «Kafka нам вообще зачем?»

Я открываю JobFlow после пятого урока.

Говорю:

> «До сих пор у нас была хорошая система. Мы создали Job, сохранили её в PostgreSQL, умеем делать transaction, умеем защищать state.»
>
> «Теперь я хочу добавить одну очень обычную backend-возможность: после создания Job сообщить worker, что её нужно обработать.»

Показываю:

```text
POST /jobs
   ↓
JobService
   ↓
PostgreSQL
```

И добавляю:

```text
PostgreSQL
   ↓
"job.created"
   ↓
worker
```

Спрашиваю:

> «Как проще всего это сделать?»

<details style="display: inline;"><summary>Помощь</summary>Можно ли просто вызвать функцию worker из `JobService`?</details>

<details style="display: inline;"><summary>Ответ</summary>Можно, но тогда producer и worker жёстко связаны одним процессом и одним execution path. Нам нужна асинхронная граница между ними.</details>

Продолжаю:

> «Вот здесь появляется брокер сообщений. В нашем случае Kafka/Redpanda.»

---

# 5. 00:10–00:25 — Плакат 06A: DISTRIBUTED MODEL

Открываем `06A`.

Я говорю:

> «Посмотрите на первый плакат. До этого урока многие операции выглядели как function call: вызвали и получили result. Теперь между сторонами появляется сеть.»

---

## 06A / Card 1 — BOUNDARY

> «Что принципиально меняется между локальным вызовом и distributed call?»

<details><summary>Помощь</summary>Что нового появляется между `send` и `result`?</details>

<details><summary>Ответ</summary>Сеть, отдельный процесс, независимые lifecycle, latency и partial failure.</details>

Пишу:

```text
local:

call → result


distributed:

send → network → ???
```

И спрашиваю:

> «Почему я ставлю `???`?»

<details><summary>Ответ</summary>Потому что результат может быть задержан, потерян, продублирован или стать неизвестным для инициатора.</details>

---

## 06A / Card 2 — LATENCY

> «Что означает timeout?»

<details><summary>Ответ</summary>Истёк допустимый срок ожидания результата. Это не доказывает, что сама операция не произошла.</details>

Очень важный пример:

```text
producer
   │
   │ request
   ▼
consumer
   │
   │ commit
   ▼
state

producer ← timeout
```

Спрашиваю:

> «Операция произошла или нет?»

<details><summary>Ответ</summary>Producer может этого не знать.</details>

И я фиксирую:

> **UNKNOWN RESULT — это нормальная distributed-systems проблема.**

---

## 06A / Card 3 — LOSS / UNKNOWN

Теперь:

```text
producer ─── X ─── consumer
```

> «Если producer потерял connection, что именно он знает?»

<details><summary>Ответ</summary>Он знает только, что не получил ожидаемый результат. Он не обязательно знает, дошло ли сообщение и был ли выполнен effect.</details>

Я говорю:

> «Вот поэтому distributed retry сложнее обычного retry function call.»

---

## 06A / Card 4 — DUPLICATE

Теперь:

```text
message
 ↓
process
 ↓
ACK ✗
 ↓
message again
```

> «Это ошибка брокера?»

<details><summary>Ответ</summary>Не обязательно. Duplicate delivery может быть нормальной семантикой системы, например при at-least-once processing и crash между effect и фиксацией delivery.</details>

И спрашиваю:

> «Что будет проблемой, если effect не идемпотентен?»

<details><summary>Ответ</summary>Один логический message может создать несколько нежелательных внешних эффектов.</details>

---

## 06A / Card 5 — ORDER

Показываю:

```text
send:
1 → 2 → 3

receive:
1 → 3 → 2
```

> «Можно ли вообще говорить: "в Kafka сообщения всегда приходят по порядку"?»

<details><summary>Ответ</summary>Нет. Порядок гарантируется только в конкретной области упорядочивания, например внутри partition.</details>

---

## 06A / Card 6 — PARTIAL FAILURE

Теперь главный пример:

```text
DB COMMIT ✓
NETWORK     X
CLIENT      ?
```

> «Это успех или failure?»

<details><summary>Ответ</summary>С точки зрения database state commit успешен, но взаимодействие целиком не завершилось для другой стороны. Это partial failure и потенциально unknown result.</details>

Я говорю:

> «Именно поэтому распределённая корректность нельзя свести к одному boolean `success`.»

---

# 6. 00:25–00:35 — Первый combat: наивный Kafka publish

Теперь перестаём смотреть на плакат.

Я предлагаю ученику:

```go
func (s *JobService) Create(ctx context.Context, job Job) error {
    if err := s.repo.Add(ctx, job); err != nil {
        return err
    }

    return s.publisher.Publish(ctx, JobCreated{
        JobID: job.ID,
    })
}
```

Спрашиваю:

> «Выглядит нормально?»

<details><summary>Ответ</summary>Снаружи выглядит логично, но состояние БД и публикация Kafka находятся в разных системах и не имеют одной локальной transaction boundary.</details>

Теперь ломаем.

### Failure A

```text
DB commit ✓
Kafka publish ✗
```

> «Что потеряли?»

<details><summary>Ответ</summary>В БД Job существует, но сообщение worker может никогда не получить.</details>

### Failure B

```text
Kafka publish ✓
DB commit ✗
```

> «Что теперь?»

<details><summary>Ответ</summary>Worker получил событие о состоянии, которое не было успешно зафиксировано в БД.</details>

Я говорю:

> «Вот мы пришли к необходимости гарантии между state change и message intent.»

И открываю карточку `OUTBOX` на 06C.

---

# 7. 00:35–00:50 — Плакат 06C: OUTBOX

## Card 5 — OUTBOX

Я рисую:

```text
PostgreSQL transaction

jobs
  +
outbox_events
```

и затем:

```text
outbox_events
     ↓
relay
     ↓
Kafka
```

Спрашиваю:

> «Что изменилось?»

<details><summary>Ответ</summary>State change и намерение отправить событие теперь фиксируются одной локальной transaction в PostgreSQL.</details>

> «Что всё ещё не гарантировано?»

<details><summary>Ответ</summary>Kafka publish всё ещё может завершиться ошибкой после commit. Поэтому relay должен уметь повторять публикацию.</details>

Вот здесь я хочу, чтобы ученик сам произнёс:

> **«Значит outbox не убирает retry. Он переносит надёжность доставки в отдельный контролируемый процесс.»**

---

# 8. 00:50–01:00 — Пишем Outbox

Теперь ученик добавляет таблицу.

Например:

```sql
CREATE TABLE outbox_events (
    id uuid PRIMARY KEY,
    event_type text NOT NULL,
    aggregate_id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamptz NOT NULL,
    published_at timestamptz
);
```

Спрашиваю:

> «Что здесь является identity события?»

<details><summary>Ответ</summary>`id`, который должен быть стабильным и уникальным для конкретного event.</details>

> «Зачем `aggregate_id`?»

<details><summary>Ответ</summary>Чтобы связать event с конкретной бизнес-сущностью, например Job.</details>

Теперь transaction:

```text
BEGIN

INSERT jobs

INSERT outbox_events

COMMIT
```

Вопрос:

> «Что произойдёт, если первый INSERT прошёл, а второй упал?»

<details><summary>Ответ</summary>Transaction rollback отменит оба изменения.</details>

---

# 9. 01:00–01:15 — Плакат 06B: DELIVERY

Открываем `06B`.

Теперь говорю:

> «Теперь сообщение действительно попадёт в Kafka. Следующий вопрос: что именно мы гарантируем worker?»

---

## 06B / Card 1 — AT-MOST-ONCE

> «Что значит at-most-once?»

<details><summary>Ответ</summary>Сообщение обрабатывается не более одного раза, но может быть потеряно.</details>

---

## Card 2 — AT-LEAST-ONCE

> «Что значит at-least-once?»

<details><summary>Ответ</summary>Повторная доставка возможна, зато нормальная стратегия не должна терять сообщение из-за фиксации delivery раньше или некорректно относительно обработки.</details>

Я сразу говорю:

> «Это будет наша модель для JobFlow.»

И спрашиваю:

> «Что нам теперь необходимо?»

<details><summary>Ответ</summary>Idempotent consumer / обработчик, потому что duplicate возможен.</details>

---

## Card 3 — EXACTLY-ONCE

> «А можно просто сказать: Kafka умеет exactly-once, значит задача решена?»

<details><summary>Ответ</summary>Нет. Exactly-once semantics зависит от границы транзакции и внешнего эффекта. Нельзя автоматически считать внешний side effect exactly-once только из-за свойства транспорта.</details>

> «Что важно определить?»

<details><summary>Ответ</summary>Где одновременно фиксируются message progress, state change и effect.</details>

---

## Card 4 — ACK / OFFSET

Я рисую:

```text
process
   ↓
commit offset
```

И намеренно неправильный вариант:

```text
commit offset
   ↓
process
```

Спрашиваю:

> «Что произойдёт, если процесс упадёт после первого шага?»

<details><summary>Ответ</summary>Broker может считать message обработанным, хотя effect не произошёл. Получается потеря обработки.</details>

А если наоборот:

```text
process
 ↓
crash
 ↓
no offset commit
 ↓
redelivery
```

Что получаем?

<details><summary>Ответ</summary>Повторную обработку — то есть at-least-once с возможным duplicate.</details>

И говорю:

> **«Нам проще безопасно пережить duplicate, чем пытаться сделать мир без duplicate.»**

---

# 10. 01:15–01:25 — Partition / Ordering

Показываю:

```text
topic
 ├── P0
 ├── P1
 └── P2
```

И:

```text
job-42 → partition 1
```

Спрашиваю:

> «Почему нам может быть полезно привязать `job_id` к partition key?»

<details><summary>Ответ</summary>Чтобы сообщения для одной Job попадали в одну область упорядочивания и сохраняли порядок относительно друг друга.</details>

Но:

> «Это гарантирует глобальный порядок всех Job?»

<details><summary>Ответ</summary>Нет. Глобального порядка между разными partitions нет.</details>

---

# 11. 01:25–01:35 — Consumer Group

Показываю:

```text
P0 → Worker A
P1 → Worker B
P2 → Worker C
```

> «Что произойдёт, если у нас 100 workers, а partitions всего 3?»

<details><summary>Ответ</summary>В одной consumer group одновременно полезно работать смогут только consumers, которым назначены partitions; лишние consumers не получат partition для чтения.</details>

Связываем с concurrency:

> «Вот вам ещё один bound. На предыдущем уроке мы ограничивали workers. Здесь параллелизм дополнительно ограничивается структурой messaging topology.»

---

# 12. 01:35–01:50 — Пишем Producer / Relay

Теперь ученик пишет маленький outbox relay.

Логика:

```text
poll outbox
   ↓
publish Kafka
   ↓
mark published
```

Я сразу задаю:

> «Можно ли делать `mark published` до Kafka publish?»

<details><summary>Ответ</summary>Нет, иначе crash после mark и до publish может привести к потере event.</details>

> «А можно после publish?»

<details><summary>Ответ</summary>Да, но тогда возможен duplicate: publish прошёл, процесс упал до mark, следующий цикл опубликует event ещё раз.</details>

И вот тут я хочу, чтобы ученик сам сказал:

> **«Значит relay должен быть рассчитан на duplicate.»**

---

# 13. 01:50–02:05 — Плакат 06C: IDEMPOTENCY

## Card 1 — IDEMPOTENCY

> «Что значит idempotent consumer?»

<details><summary>Ответ</summary>Повторное получение того же логического события не приводит к новому нежелательному эффекту.</details>

Для JobFlow:

```text
event_id
   ↓
processed_events
```

Спрашиваю:

> «Какой invariant мы хотим?»

<details><summary>Ответ</summary>Один логический event не должен создавать более одного требуемого effect.</details>

Например:

```sql
CREATE TABLE processed_events (
    event_id uuid PRIMARY KEY,
    processed_at timestamptz NOT NULL
);
```

Почему primary key?

<details><summary>Ответ</summary>Она сама защищает уникальность обработки на уровне persisted state.</details>

---

# 14. 02:05–02:15 — Idempotency через state transition

Теперь связываем с PostgreSQL.

Например:

```sql
UPDATE jobs
SET status = 'running'
WHERE id = $1
  AND status = 'pending';
```

Что хорошего?

<details><summary>Ответ</summary>Повторный вызов после перехода `pending → running` больше не соответствует условию и не создаёт второй переход.</details>

Но спрашиваю:

> «Это полностью решает idempotency?»

<details><summary>Ответ</summary>Только для конкретного state transition. Если рядом есть внешний effect, его семантика тоже должна быть определена.</details>

---

# 15. 02:15–02:25 — RETRY

Теперь сценарий:

```text
consumer
   ↓
downstream request
   ↓
timeout
```

> «Повторяем?»

<details><summary>Ответ</summary>Сначала определяем, retryable ли ошибка и безопасен ли повторяемый effect.</details>

Показываю:

```text
100ms
 ↓
200ms
 ↓
400ms
 ↓
800ms
```

Спрашиваю:

> «Зачем backoff?»

<details><summary>Ответ</summary>Чтобы повторные попытки не создавали новый burst нагрузки на уже деградировавшую систему.</details>

> «Зачем jitter?»

<details><summary>Ответ</summary>Чтобы множество клиентов/consumers не повторяли запрос одновременно по одинаковому расписанию.</details>

---

# 16. 02:25–02:35 — Retry Policy

Я даю таблицу словами.

```text
timeout
connection reset
temporary unavailable
```

> «Все retryable?»

<details><summary>Ответ</summary>Не автоматически. Retryability определяется контрактом операции и типом ошибки.</details>

```text
invalid payload
permission denied
business invariant violation
```

> «Их повторяем?»

<details><summary>Ответ</summary>Обычно нет, если причина не изменится от повторения.</details>

Главная формула:

> **Retry — policy, а не реакция на любое `error != nil`.**

---

# 17. 02:35–02:50 — Main Combat: JobFlow Kafka/Redpanda

Теперь основная практика.

Я говорю:

> «До этого момента мы рассуждали. Теперь строим.»

Целевая схема:

```text
                   ┌────────────┐
                   │ PostgreSQL │
                   └─────┬──────┘
                         │
                    transaction
                         │
                ┌────────┴────────┐
                │                 │
              jobs            outbox
                │                 │
                └────────┬────────┘
                         │
                       relay
                         │
                         ▼
                  Kafka / Redpanda
                         │
               ┌─────────┼─────────┐
               ▼         ▼         ▼
              W1        W2        W3
               │         │         │
               └─────────┼─────────┘
                         ▼
                    PostgreSQL
```

### Задание ученику

Нужно спроектировать:

```text
Create Job
    ↓
persist Job
    ↓
persist event
    ↓
relay
    ↓
Kafka
    ↓
consumer
    ↓
process Job
    ↓
idempotent state transition
```

---

# 18. 02:50–03:05 — Ошибки как часть задания

Я специально предлагаю ученику **самому назвать**, что может пойти не так.

> «Дайте мне минимум пять failure scenarios.»

<details><summary>Помощь</summary>Посмотрите на три плаката: сеть, delivery и effect.</details>

<details><summary>Ответ</summary>Kafka недоступна, relay падает после publish, consumer падает после effect, сообщение приходит повторно, consumer обрабатывает медленно, downstream timeout, порядок событий нарушен, database временно недоступна.</details>

Теперь:

> «Для каждого назовите policy.»

Пример:

```text
Kafka unavailable
→ retry + backoff

consumer crash after effect
→ redelivery + idempotency

DB unavailable
→ retry according to DB error policy

duplicate event
→ deduplication / idempotent transition
```

---

# 19. 03:05–03:15 — Failure Injection

Теперь начинается настоящий combat.

## Сценарий 1

Relay падает после:

```text
Kafka publish ✓
```

но до:

```text
outbox marked published
```

Что произойдёт?

<details><summary>Ответ</summary>Следующий relay может опубликовать тот же event повторно.</details>

Какой механизм должен его пережить?

<details><summary>Ответ</summary>Idempotent consumer / deduplication.</details>

---

## Сценарий 2

Consumer записал state и умер до фиксации offset.

Что произойдёт?

<details><summary>Ответ</summary>Message будет redelivered, consumer повторит обработку.</details>

Нормально?

<details><summary>Ответ</summary>Да, если effect идемпотентен.</details>

---

## Сценарий 3

Kafka недоступна 30 секунд.

Что делает relay?

<details><summary>Ответ</summary>Оставляет outbox event непубликовавшимся и использует retry/backoff policy. Event не должен исчезать только потому, что broker временно недоступен.</details>

---

# 20. 03:15–03:25 — Interview Drill

Теперь закрываю код.

### Вопрос

> «Что такое distributed system?»

<details><summary>Ответ</summary>Система из независимых компонентов, взаимодействующих через коммуникационные границы, где задержки, отказ и частичное знание становятся частью поведения системы.</details>

### «Что такое at-least-once?»

<details><summary>Ответ</summary>Сообщение может быть доставлено повторно, поэтому consumer должен корректно переживать duplicate.</details>

### «Почему duplicate нормален?»

<details><summary>Ответ</summary>Например, consumer может выполнить effect и упасть до фиксации offset, после чего broker доставит message снова.</details>

### «Что такое idempotency?»

<details><summary>Ответ</summary>Повтор логической операции не создаёт новый нежелательный effect.</details>

### «Что такое outbox?»

<details><summary>Ответ</summary>Паттерн, при котором изменение локального состояния и запись намерения отправить событие фиксируются одной transaction, после чего отдельный relay доставляет event в broker.</details>

### «Что такое partial failure?»

<details><summary>Ответ</summary>Часть распределённой операции завершилась, а другая часть нет, поэтому разные компоненты могут иметь разное представление о результате.</details>

### «Почему timeout опасен?»

<details><summary>Ответ</summary>Потому что он сообщает только об окончании ожидания, а не обязательно о том, была ли операция выполнена.</details>

---

# 21. 03:25–03:35 — Code Review

Теперь смотрим JobFlow.

Я спрашиваю:

> «Где здесь проходит граница между local state и distributed effect?»

<details><summary>Ответ</summary>PostgreSQL transaction локально фиксирует state и outbox intent; Kafka delivery происходит за пределами этой transaction boundary.</details>

> «Где возможен duplicate?»

<details><summary>Ответ</summary>После Kafka publish до фиксации outbox published state; после consumer effect до offset commit; при retry.</details>

> «Где он обезвреживается?»

<details><summary>Ответ</summary>На уровне idempotent consumer, deduplication key или безопасного state transition.</details>

> «Где bounded recovery?»

<details><summary>Ответ</summary>В timeout, retry limit, exponential backoff, queue bounds и failure policy.</details>

---

# 22. 03:35–03:45 — Final JobFlow Architecture

К концу урока я хочу увидеть:

```text
                              CLIENT
                                │
                                ▼
                           HTTP API
                                │
                                ▼
                          JobService
                           /       \
                          ▼         ▼
                     PostgreSQL   policy
                          │
                  ┌───────┴───────┐
                  │ transaction   │
                  │               │
                  ▼               ▼
                jobs           outbox
                                  │
                                  ▼
                               RELAY
                                  │
                                  ▼
                           Kafka / Redpanda
                                  │
                      ┌───────────┼───────────┐
                      ▼           ▼           ▼
                    Worker      Worker      Worker
                      │           │           │
                      └───────────┼───────────┘
                                  ▼
                             PostgreSQL
                                  │
                         idempotent state
```

Я спрашиваю:

> «Что из этого является guaranteed exactly-once?»

<details><summary>Ответ</summary>Не надо объявлять всю схему exactly-once. Нужно определить delivery semantics каждого участка и отдельно обеспечить корректность effect.</details>

Это очень важный финальный ответ.

---

# 23. 03:45–03:55 — Финальная проверка

### Q1

> Почему нельзя просто `DB commit → Kafka publish`?

<details><summary>Ответ</summary>Потому что это две разные failure domains. Между ними возможен partial failure.</details>

### Q2

> Что решает outbox?

<details><summary>Ответ</summary>Надёжно фиксирует local state change и message intent в одной DB transaction.</details>

### Q3

> Решает ли outbox duplicate?

<details><summary>Ответ</summary>Нет. Relay может опубликовать event и упасть до пометки published, поэтому duplicate возможен.</details>

### Q4

> Что решает duplicate?

<details><summary>Ответ</summary>Idempotent consumer, deduplication или безопасная state transition semantics.</details>

### Q5

> Что делать с timeout?

<details><summary>Ответ</summary>Считать результат неизвестным и применять заранее определённую recovery/idempotency policy.</details>

### Q6

> Зачем partition key?

<details><summary>Ответ</summary>Чтобы определить область упорядочивания и привязать связанные события к одной partition.</details>

### Q7

> От чего зависит consumer parallelism?

<details><summary>Ответ</summary>От числа partitions, consumer group topology и внутреннего worker concurrency.</details>

### Q8

> Зачем backoff?

<details><summary>Ответ</summary>Чтобы retries не создавали дополнительный load spike и не усиливали отказ downstream.</details>

### Q9

> Что значит at-least-once?

<details><summary>Ответ</summary>Duplicate возможен и должен быть безопасен для consumer.</details>

### Q10

> Что такое idempotency key?

<details><summary>Ответ</summary>Стабильный идентификатор логической операции/event, по которому система может распознать повтор.</details>

---

# 24. 03:55–04:00 — Closing

Я заканчиваю урок так:

> «Сегодня мы сделали последний большой переход.
>
> До этого у нас были Go goroutines, worker pool и PostgreSQL transaction.
>
> Внутри одного процесса мы могли примерно контролировать всё состояние.
>
> Сегодня появилась сеть.
>
> И с ней появилась неприятная правда: результат операции иногда невозможно узнать сразу.
>
> Сообщение может задержаться.
>
> Может потеряться.
>
> Может прийти повторно.
>
> Может прийти в другом порядке.
>
> Одна часть системы может уже завершить работу, а другая ещё ничего об этом не знать.
>
> Поэтому distributed engineering — это не попытка сделать сеть магически надёжной.
>
> Это проектирование такой семантики, при которой задержка, повтор, отказ и неизвестный результат остаются управляемыми.
>
> Наш JobFlow теперь уже не просто backend на Go.
>
> У него есть persistent state, concurrent workers, message broker, delivery semantics, outbox и idempotent processing.
>
> Теперь мы можем начинать говорить не только о коде, но о поведении системы целиком.»

И последняя фраза:

> **«Если локальная программа должна быть корректной, распределённая система должна ещё и правильно переживать свою некорректную реальность.»**

---

# 25. Домашнее задание — JobFlow v4: Kafka/Redpanda

Домашнее задание специально строится по принципу:

> **Я даю лопату и удочку. Искать месторождение и ловить рыбу вы будете сами.**

На созвоне мы построили минимальную схему. Теперь ученик самостоятельно доводит её до работающей версии.

## 1. Outbox

Добавить:

```text
outbox_events
```

Минимально:

```text
id
event_type
aggregate_id
payload
created_at
published_at
```

Реализовать запись `JobCreated` **в той же PostgreSQL transaction**, что и создание Job.

---

# 26. Outbox Relay

Создать отдельный компонент:

```text
outbox relay
```

Он должен:

```text
select unpublished
        ↓
publish Kafka
        ↓
mark published
```

Требования:

```text
context cancellation
bounded polling
retry
backoff
graceful shutdown
```

---

# 27. Kafka / Redpanda

Поднять broker локально.

Создать topic:

```text
job.events
```

Определить key:

```text
job.ID
```

И убедиться, что события одной Job имеют ожидаемую ordering semantics.

---

# 28. Consumer

Создать consumer group:

```text
job-workers
```

Consumer должен:

```text
read event
   ↓
validate
   ↓
idempotency check
   ↓
process
   ↓
state transition
```

---

# 29. Idempotency

Добавить:

```text
processed_events
```

Например:

```text
event_id PK
processed_at
```

Повторная доставка одного `event_id` **не должна приводить к повторному бизнес-effect**.

---

# 30. Failure Injection

Ученик должен самостоятельно реализовать и проверить минимум пять сценариев:

```text
1. Kafka unavailable

2. relay crash after publish

3. consumer crash after DB effect

4. duplicate event

5. downstream timeout
```

Для каждого написать:

```text
observed effect
expected semantics
recovery policy
```

---

# 31. Retry Policy

Добавить:

```text
max attempts
exponential backoff
jitter
```

Например концептуально:

```text
100ms
200ms
400ms
800ms
...
```

Но конкретные значения выбрать самостоятельно и обосновать.

---

# 32. Tests

Обязательно:

```bash
go test ./...
go test -race ./...
```

И тесты на:

```text
outbox atomicity
duplicate event
idempotent consumer
retry
cancellation
shutdown
```

---

# 33. Acceptance Criteria

## PASS

Ученик:

* объясняет partial failure;
* объясняет timeout ≠ failure;
* различает at-most-once / at-least-once;
* аккуратно объясняет exactly-once;
* понимает partition/order;
* понимает consumer group;
* реализует outbox;
* реализует Kafka/Redpanda producer;
* реализует consumer;
* имеет idempotency strategy;
* имеет retry/backoff policy;
* переживает duplicate;
* корректно завершает relay и consumer;
* может объяснить failure semantics каждого участка pipeline.

## Engineering quality

```text
delivery semantics explicit
failure policy explicit
idempotency key explicit
retry bounded
concurrency bounded
shutdown defined
observable effects recorded
```

---

# 34. Итоговая mental model

После урока у ученика должна остаться одна цепочка:

```text
LOCAL STATE
     ↓
TRANSACTION
     ↓
OUTBOX
     ↓
BROKER
     ↓
DELIVERY
     ↓
DUPLICATE
     ↓
IDEMPOTENCY
     ↓
EFFECT
     ↓
VERIFY
```

И ключевая инженерная формула:

> **Distributed correctness = explicit delivery semantics + idempotent effects + bounded recovery.**

Это уже естественный мост к следующему уроку: **BACKEND ARCHITECTURE**, где мы соберём HTTP, application layer, PostgreSQL, Kafka и workers в полноценную систему с явными control/data boundaries.
