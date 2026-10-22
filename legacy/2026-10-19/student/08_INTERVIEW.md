# Go Backend Combat — 08. INTERVIEW

**Version:** 1.0.0
**Topic:** Mock Interview
**Format:** Подготовка ученика к mock interview + 500 вопросов по всему курсу
**Project:** JobFlow
**Duration:** 60 минут mock interview + самостоятельная подготовка
**Prerequisite:** `01_LANGUAGE.md` — `07_BACKEND_ARCHITECTURE.md`

---

# 1. Зачем этот материал

Финальное занятие не вводит новых технологий.

Его задача:

```text
знание
  ↓
быстрый recall
  ↓
точный ответ
  ↓
инженерное обоснование
  ↓
защита решения
```

Сеньор интервью проверяет не только:

> «Знаете ли вы термин?»

Но и:

> «Можете ли вы за 30–60 секунд точно объяснить смысл, выбрать механизм и защитить trade-off?»

---

# 2. Как готовиться

## 2.1. Не учить 500 ответов наизусть

500 вопросов — это карта покрытия.

Цель подготовки:

```text
question
   ↓
mental model
   ↓
answer
   ↓
example
   ↓
trade-off
```

На каждый вопрос старайтесь отвечать без IDE и без поиска.

---

# 3. Формула сильного ответа

Для большинства engineering questions используйте:

```text
DEFINITION
→ WHY
→ EXAMPLE
→ TRADE-OFF
```

Например:

> Что такое idempotency?

Сильный ответ:

> «Idempotency — свойство, при котором повторение одной логической операции не создаёт новый нежелательный effect. В JobFlow это нужно потому, что Kafka работает с at-least-once и duplicate возможен. Мы можем использовать `event_id` как idempotency key и защищать обработку через persisted state или `processed_events`. Цена — дополнительное состояние и cleanup policy.»

Это существенно сильнее:

> «Idempotency — это когда повторный запрос не меняет результат.»

---

# 4. Если не знаете ответ

Не нужно начинать импровизировать.

Используйте структуру:

> «Точную деталь не помню, но я бы рассуждал так…»

Затем:

```text
requirement
→ invariant
→ boundary
→ mechanism
→ evidence
```

Например:

> «Я не помню точную внутреннюю реализацию scheduler, но понимаю семантику: goroutine становится runnable, scheduler выбирает execution resource, runtime может переключать выполнение и будить goroutine после I/O.»

Это лучше, чем уверенно сообщать неверный факт.

---

# 5. Темп

Тренировочный режим:

```text
0–10 sec
→ понять вопрос

10–30 sec
→ дать основную модель

30–45 sec
→ пример

45–60 sec
→ trade-off / failure
```

Для простого определения обычно достаточно 15–30 секунд.

Не нужно начинать с длинной предыстории.

Сначала ответ.

Потом доказательство.

---

# 6. Что делать, когда интервьюер атакует решение

Интервьюер может сказать:

> «А если PostgreSQL упал?»

> «А если Kafka прислала сообщение дважды?»

> «А если timeout, но операция уже прошла?»

> «А если 100 000 запросов?»

Не защищайте первый вариант любой ценой.

Используйте:

```text
new constraint
→ affected invariant
→ changed design
→ new trade-off
```

Например:

> «При duplicate delivery наша текущая модель должна уже быть идемпотентной. Если это не так, я добавлю persisted idempotency key или сделаю state transition conditional.»

---

# 7. Что особенно важно повторить

Перед интервью отдельно повторить:

```text
Go type system
value / pointer
interfaces
errors
slice / map
goroutines
channels
select
context
WaitGroup
mutex / atomic
race detector
GMP
worker pool
backpressure
SQL
indexes
transactions
MVCC
isolation
ACID
pgxpool
migrations
Kafka partitions
consumer groups
delivery semantics
idempotency
Outbox
retry/backoff
HTTP
application architecture
composition root
Redis
cache-aside
graceful shutdown
failure analysis
```

---

# 8. Что интервьюер должен услышать в ответах

Сильный backend engineer естественно говорит словами:

```text
ownership
invariant
contract
boundary
failure mode
resource limit
idempotency
cancellation
observability
trade-off
blast radius
source of truth
```

Не нужно вставлять их искусственно.

Они должны появляться как результат рассуждения.

---

# 9. Mock Interview Protocol

Основной режим:

```text
Question
 ↓
30 seconds
 ↓
Hint
 ↓
30 seconds
 ↓
Answer
 ↓
Follow-up attack
```

Каждый вопрос оценивается по четырём уровням:

```text
0 — не знает
1 — узнаёт термин
2 — объясняет семантику
3 — объясняет + применяет + защищает trade-off
```

Для strong middle+/senior основной target:

> **уровень 3.**

---

# 10. Финальная подготовка перед mock interview

За неделю до интервью не нужно перечитывать весь курс целиком.

Лучше пройти 4 круга.

## Круг 0 - Posters

Повторить все [постеры](../posters/README.md) освежить в памяти ключевые вещи.

## Круг 1 — Recall

На каждый из [500 вопросов](../interview/README.md):

```text
извлечь ответ
без подсказки
```

Три состояния:

```text
GREEN  → отвечаю уверенно
YELLOW → знаю, но медленно
RED    → не знаю
```

## Круг 2 — Compression

Каждый `RED/YELLOW` сжать до одной фразы.

Например:

```text
MVCC
→ versions + visibility

Outbox
→ local state + event intent in one DB transaction

DIP
→ policy depends on stable contract, not detail

Backpressure
→ producer must respect consumer capacity
```

## Круг 3 — Attack

Для каждого сложного вопроса задать себе ещё три:

```text
А что если failure?
А что если retry?
А что если нагрузка ×10?
```

Именно здесь начинается уровень senior.

---

# 11. Последняя формула курса

На интервью практически любой backend-вопрос можно разбирать так:

```text
REQUIREMENT
    ↓
INVARIANT
    ↓
BOUNDARY
    ↓
MECHANISM
    ↓
FAILURE
    ↓
EVIDENCE
```

И ответы становятся намного устойчивее.

Не:

> «Я бы использовал Kafka.»

А:

> «Нам нужна асинхронная доставка между независимыми компонентами. При этом возможен partial failure, поэтому локальное состояние и intent фиксируем через Outbox, доставку делаем at-least-once, а effect защищаем idempotency. Kafka подходит как durable message transport; конкретный partition key определяется требованием ordering.»

Это именно тот стиль ответа, который нужно тренировать на mock interview.
