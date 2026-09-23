# 🔥 Senior Go Backend Combat

> **Боевой курс по GO для выхода на middle+/senior. 1 production проект. 1 инженерная модель мышления. 1 финальное mock interview.**

![LOGO](./LOGO.png)

**Senior Go Backend Combat** — интенсив по Go backend engineering для разработчика, который уже умеет программировать, но хочет быстро перейти к уровню **strong middle / senior interview mindset**.

---

# 🧭 Навигация по курсу

| Урок   | Тема                                                                   | Теория                                       | Практика                                      | Ментор                                      |
| ------ | ---------------------------------------------------------------------- | -------------------------------------------- | --------------------------------------------- | ------------------------------------------- |
| **01** | Модель мышления Engineering Mindset                                    | [Теория](student/01_ENGINEERING_MINDSET.md)  | [Практика](combat/01_ENGINEERING_MINDSET.md)  | [Ментор](mentor/01_ENGINEERING_MINDSET.md)  |
| **01** | Go Language + JobFlow v0                                               | [Теория](student/02_GO_LANGUAGE.md)          | [Практика](combat/02_GO_LANGUAGE.md)          | [Ментор](mentor/02_GO_LANGUAGE.md)          |
| **02** | Go Concurrency, Worker Pool, cancellation, synchronization             | [Теория](student/03_CONCURRENCY.md)          | [Практика](combat/03_CONCURRENCY.md)          | [Ментор](mentor/03_CONCURRENCY.md)          |
| **03** | Design Principles, SOLID                                               | [Теория](student/04_DESIGN_PRINCIPLES.md)    | [Практика](combat/04_DESIGN_PRINCIPLES.md)    | [Ментор](mentor/04_DESIGN_PRINCIPLES.md)    |
| **04** | Data & Persistence,  PostgreSQL, SQL, Goose, pgxpool, transactions     | [Теория](student/05_DATA_PERSISTENCE.md)     | [Практика](combat/05_DATA_PERSISTENCE.md)     | [Ментор](mentor/05_DATA_PERSISTENCE.md)     |
| **05** | Distributed Systemsб Kafka/Redpanda, Outbox, idempotency, retry        | [Теория](student/06_DISTRIBUTED_SYSTEMS.md)  | [Практика](combat/06_DISTRIBUTED_SYSTEMS.md)  | [Ментор](mentor/06_DISTRIBUTED_SYSTEMS.md)  |
| **06** | Backend Architecture, HTTP, application layer, Redis, composition root | [Теория](student/07_BACKEND_ARCHITECTURE.md) | [Практика](combat/07_BACKEND_ARCHITECTURE.md) | [Ментор](mentor/07_BACKEND_ARCHITECTURE.md) |
| **07** | Project Completion, Финальная сборка, smoke tests, failure review      | [Теория](student/07_BACKEND.md)              | [Практика](combat/07_BACKEND.md)              | [Ментор](mentor/07_BACKEND_ARCHITECTURE.md) |
| **08** | Mock Interview 60-минутное интервью по всему курсу                     | [Теория](student/08_INTERVIEW.md)            | [Практика](combat/08_INTERVIEW.md)            | [Ментор](mentor/08_INTERVIEW.md)            |

---


# 🎯 Что должно измениться после курса

Мы не изучаем Go как набор синтаксических конструкций. Мы строим один backend-проект — **JobFlow** — и постепенно превращаем его из небольшой Go-программы в распределённую production-like систему.

Главная задача курса:

> **Научиться не угадывать «правильный паттерн», а принимать, объяснять и защищать инженерное решение.**

К концу курса ученик должен уметь не только написать Go-код, но и объяснить:

```text
Что система должна делать?
Какой invariant сохраняем?
Где boundary?
Кто владеет state?
Какой механизм нужен?
Что произойдёт при concurrency?
Что произойдёт при failure?
Как ограничены ресурсы?
Как доказать correctness?
```

То есть на интервью мы ожидаем не ответ:

> «Я бы поставил Kafka и Redis».

А ответ:

> «PostgreSQL является authoritative state. Kafka используется как asynchronous delivery boundary. Между DB state и event intent стоит Outbox transaction. Consumer работает с at-least-once semantics, поэтому effect идемпотентен. Redis — derived cache и может деградировать без потери correctness. Worker concurrency ограничена, cancellation распространяется через context, а recovery bounded retry policy».

---

# 🥊 ПЕТ-Проект курса — JobFlow

Весь курс строится вокруг одного пет-проекта, который пишет ученик в процессе освоения материалов курса.

JobFlow — backend для создания и обработки фоновых задач.

К концу практической части он выглядит концептуально так:

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
                    │             │
                    ▼             ▼
               PostgreSQL      Kafka
                    │             │
                    │          Workers
                    │             │
                    └──────┬──────┘
                           │
                        Redis
                        cache
```

Проект последовательно проходит через:

```text
structs / methods
→ interfaces / errors / tests
→ goroutines / channels / worker pool
→ SOLID / cohesion / coupling
→ PostgreSQL / SQL / migrations
→ transactions / MVCC / isolation
→ Kafka / Redpanda
→ Outbox / idempotency / retry
→ HTTP API
→ Redis cache
→ composition root
→ graceful shutdown
→ failure analysis
```

---

# 📚 Структура репозитория

```text
.
├── README.md
│
├── mentor/
│   ├── 01_ENGINEERING_MINDSET.md
│   ├── 02_GO_LANGUAGE.md
│   ├── 03_CONCURRENCY.md
│   ├── 04_DESIGN_PRINCIPLES.md
│   ├── 05_DATA_PERSISTENCE.md
│   ├── 06_DISTRIBUTED_SYSTEMS.md
│   ├── 07_BACKEND_ARCHITECTURE.md
│   └── 08_INTERVIEW.md
│
├── student/
│   └── ...
│
├── posters/
│   ├── 01_ENGINEERING_MINDSET.svg
│   ├── 02_LANGUAGE.svg
│   ├── 03A_CONCURRENCY_MODEL.svg
│   ├── 03B_CONCURRENCY_RUNTIME.svg
│   ├── 03C_CONCURRENCY_ENGINEERING.svg
│   ├── 04A_SOLID.svg
│   ├── 04B_DESIGN_ENGINEERING.svg
│   ├── 05A_DATA_MODEL.svg
│   ├── 05B_SQL_SCHEMA.svg
│   ├── 05C_TRANSACTIONS_CONCURRENCY.svg
│   ├── 06A_DISTRIBUTED_MODEL.svg
│   ├── 06B_DELIVERY_MESSAGING.svg
│   ├── 06C_FAILURE_EFFECTS.svg
│   ├── 07A_HTTP_API.svg
│   ├── 07B_APPLICATION_ARCHITECTURE.svg
│   ├── 07C_RUNTIME_CACHE.svg
│   └── README.md
│
└── interview/
    └── README.md
```

---

# 👨‍🏫 Материалы ментора

Папка [`mentor/`](mentor/) содержит сценарии проведения занятий.

Это не обычные конспекты. Каждый урок построен как interactive mentor-led session:

```text
EXPLAIN
 ↓
QUESTION
 ↓
PREDICT
 ↓
CODE
 ↓
RUN
 ↓
BREAK
 ↓
FIX
 ↓
VERIFY
```

Ментор не выдаёт ответ заранее.

Ученик сначала формулирует решение, затем пишет код, после чего решение специально атакуется:

```text
What if it races?
What if it times out?
What if DB fails?
What if Kafka duplicates?
What if Redis disappears?
What if the process dies here?
```

Каждый урок заканчивается:

```text
review
↓
checkpoint
↓
homework
↓
PASS / FAIL
↓
next engineering problem
```

---

# 👨‍💻 Материалы ученика

Папка [`student/`](student/) содержит материалы, которыми ученик пользуется непосредственно во время обучения и после него.

Это не текст речи ментора.

Материал ученика должен позволять:

* быстро восстановить модель после занятия;
* повторить ключевые понятия;
* самостоятельно продолжить JobFlow;
* выполнить домашнее задание;
* проверить себя перед следующим уроком;
* использовать материалы как reference перед интервью.

Главный принцип:

> **Ментор ведёт. Ученик строит.**
---
<!-- reply:"resp.md" v1.0.0 2026-08-19T07:10:00Z Topic: Delivery Summary: Краткая секция CI/CD для корневого README -->

# 🐳 CI/CD и локальная инфраструктура

JobFlow сразу работает как полноценный backend-проект: локальная инфраструктура поднимается через **Docker Compose**, команды разработки унифицированы через **Makefile**, конфигурация задаётся через `.env`, а приложение собирается multi-stage Docker build и проверяется в CI.

В `cicd/` лежат готовые файлы для копирования в корень собственного `jobflow`:

* Docker Compose: PostgreSQL, Redis, Redpanda
* Makefile: запуск, migrations, tests, race, build
* `.env` / `.env.example`
* `Dockerfile`

**[→ Полная инструкция по работе с CI/CD](cicd/README.md)**

Минимальный старт:

```bash
cp .env.example .env
make up
make migrate
make run
```

Перед commit:

```bash
make verify
```

> **Умение воспроизводимо поднять, проверить и собрать backend — часть senior backend engineering.**

---

# 🖼️ Плакаты

Весь курс имеет отдельную визуальную систему из A3-плакатов.

Плакаты работают как постоянная spatial memory курса: во время занятия ментор может вернуться к одной и той же модели, а ученик после курса — восстановить её без полного перечитывания урока.

Полное содержание и единое повествование:

**[→ posters/README.md](posters/README.md)**

Основная последовательность:

```text
01  ENGINEERING MINDSET
02  GO LANGUAGE

03A CONCURRENCY MODEL
03B CONCURRENCY RUNTIME
03C CONCURRENCY ENGINEERING

04A SOLID
04B DESIGN ENGINEERING

05A DATA MODEL
05B SQL SCHEMA
05C TRANSACTIONS CONCURRENCY

06A DISTRIBUTED MODEL
06B DELIVERY MESSAGING
06C FAILURE EFFECTS

07A HTTP API
07B APPLICATION ARCHITECTURE
07C RUNTIME CACHE
```

Плакат — это не отдельный урок.

Это **визуальная модель, к которой возвращается урок**.

---

# 🎯 Финальная подготовка к интервью

Последний этап курса — не изучение новых технологий.

Это перевод полученных знаний в режим интервью.

Ученик должен уметь быстро:

```text
понять вопрос
→ определить область
→ сформулировать invariant
→ назвать trade-off
→ предложить решение
→ атаковать своё решение
→ объяснить failure mode
```

Перед mock interview необходимо пройти два reference-набора:

### Плакаты

**[→ posters/README.md](posters/README.md)**

Это компактная карта всей инженерной модели курса.

### 500 interview questions

**[→ interview/README.md](interview/README.md)**

500 вопросов покрывают:

```text
Go
Concurrency
GMP / runtime
Channels
Synchronization
Worker pools
SQL
PostgreSQL
Transactions
MVCC
Isolation
ACID
pgxpool
Migrations
Kafka
Partitions
Consumer groups
Delivery semantics
Outbox
Idempotency
Retry
HTTP
Architecture
Redis
Caching
Observability
Production
System Design
```

Вопросы построены в формате:

```text
QUESTION

HELP
→ если ученик застрял

ANSWER
→ эталон проверки
```

Работать с ними нужно не как с флеш-картами, а как с тренировкой reasoning.

---

# 🎤 Как готовиться к mock interview

Не нужно пытаться выучить 500 ответов дословно.

Нужно натренировать **скорость восстановления модели**.

Для каждого вопроса полезно пройти четыре шага:

**1. Ответить сразу.**
Без поиска и без IDE.

**2. Если не получается — открыть «Помощь».**
Она должна дать направление мысли, но не готовый ответ.

**3. Сформулировать ответ самостоятельно.**

**4. Сравнить с эталоном.**

Особенно важно отслеживать не количество выученных терминов, а тип ошибки.

Например:

```text
не знаю термин
        ≠
не понимаю механизм
        ≠
понимаю механизм, но не могу объяснить trade-off
        ≠
понимаю trade-off, но не могу применить к системе
```

Для уровня strong middle / senior последний переход наиболее важен.

---

# 🧠 Главный режим подготовки

На интервью почти любой вопрос можно вернуть к одной инженерной последовательности:

> **Requirement → Invariant → Boundary → Mechanism → Failure → Evidence**

Например, вместо ответа:

> «Поставим Redis.»

Нужно уметь сказать:

> «Чтений намного больше записей. PostgreSQL является authoritative state. Нужен derived cache для снижения latency и DB load. Поэтому используем cache-aside с TTL. Redis outage не должен ломать correctness, поэтому read path деградирует к PostgreSQL. Проверяем hit rate, DB load и tail latency.»

Это уже не знание технологии.
Это инженерное решение.

---

# 🥊 Финальный стандарт курса

К завершению курса ученик должен уметь открыть JobFlow и пройти по нему от начала до конца:

```text
HTTP request
    ↓
transport
    ↓
application policy
    ↓
domain state
    ↓
PostgreSQL transaction
    ↓
Outbox
    ↓
Kafka
    ↓
consumer
    ↓
idempotency
    ↓
worker
    ↓
PostgreSQL
    ↓
Redis
    ↓
HTTP response
```

И на любом участке ответить:

```text
Кто владеет этим состоянием?
Какой здесь контракт?
Какие ресурсы ограничены?
Что произойдёт при timeout?
Что произойдёт при duplicate?
Что произойдёт при crash?
Как система восстановится?
Чем это доказано?
```

---

# 🏁 Definition of Done

Курс считается пройденным, когда выполнены одновременно четыре условия.

### Code

JobFlow собирается, запускается и проходит verification:

```bash
go test ./...
go test -race ./...
go vet ./...
```

### System

Работает полный pipeline:

```text
HTTP
→ PostgreSQL
→ Outbox
→ Kafka
→ Worker
→ PostgreSQL
→ Redis
```

### Engineering

Ученик может объяснить:

```text
ownership
boundaries
invariants
failure modes
resource limits
delivery semantics
trade-offs
```

### Interview

Ученик способен не только назвать технологию, но и защитить решение:

> **почему именно так, какую гарантию это даёт, что может сломаться и как это проверить.**

---

# 🔥 Final Principle

> **Не доказывай, что знаешь Go.
> Доказывай, что умеешь строить backend, управлять состоянием, ограничивать ресурсы, переживать отказ и защищать свои инженерные решения.**

**Go Backend Combat** заканчивается не тогда, когда вы запомнили Go.

Он заканчивается тогда, когда вы можете открыть пустой backend-задачу, за несколько минут определить границы проблемы и начать строить решение, которое можно объяснить, проверить и защитить.
