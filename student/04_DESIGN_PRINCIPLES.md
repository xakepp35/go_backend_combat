# Go Backend Combat — 03. DESIGN

**Version:** 1.0.0
**Topic:** Design Principles
**Format:** Дидактический материал ученика — SOLID, cohesion, coupling и практический refactoring JobFlow
**Project:** JobFlow
**Prerequisite:** `01_LANGUAGE.md`, `02_CONCURRENCY.md`
**Duration:** 120 минут
**Posters:** `04A_SOLID.svg`, `04B_DESIGN_ENGINEERING.svg`

Плакаты:

* [04A — SOLID](../posters/04A_SOLID.svg)
* [04B — Design Engineering](../posters/04B_DESIGN_ENGINEERING.svg)

---

# 1. Цель занятия

После первых двух занятий JobFlow уже имеет:

```text
Job
Repository
Service
WorkerPool
goroutines
channels
context
synchronization
tests
```

Теперь появляется новый класс проблем:

> **Код работает, но изменения становятся дорогими.**

Например:

```text
HTTP handler
 ├── validation
 ├── business rules
 ├── PostgreSQL
 ├── Kafka
 └── logging
```

Цель занятия:

```text
увидеть стоимость изменения
        ↓
найти границу
        ↓
выбрать минимальную abstraction
        ↓
изменить код
        ↓
сохранить behavior
        ↓
доказать tests
```

Главный принцип:

> **Не применяй SOLID потому, что он существует. Сначала найди стоимость изменения, затем выбери минимальную границу, которая её уменьшает.**

---

# 2. SOLID

![04A — SOLID](../posters/04A_SOLID.svg)

SOLID — пять принципов проектирования:

```text
S — Single Responsibility
O — Open / Closed
L — Liskov Substitution
I — Interface Segregation
D — Dependency Inversion
```

Они помогают управлять:

```text
responsibility
change
contracts
dependencies
```

---

# 3. S — Single Responsibility

> **Одна ось ответственности и одна основная причина для изменения.**

SRP не означает:

> «В классе должна быть ровно одна функция.»

Вопрос:

> **Какая причина изменения принадлежит этому компоненту?**

Плохой пример:

```text
HTTP handler
 ├── decode JSON
 ├── validation
 ├── business rules
 ├── SQL
 ├── Kafka
 └── logging
```

У него много независимых причин изменения:

```text
HTTP API
validation
business rules
database
messaging
logging
```

Это большой blast radius.

---

# 4. O — Open / Closed

> **Стабильную policy следует расширять новым поведением, минимизируя переписывание стабильного ядра.**

OCP не означает:

> «Код нельзя менять.»

Важный вопрос:

> **Можно ли добавить новое поведение без постоянного изменения стабильной части системы?**

Пример:

```text
stable policy
      ↑
      │
new implementation
```

OCP особенно полезен там, где:

* есть устойчивый contract;
* появляются новые реализации;
* изменение behavior должно локализоваться.

Не каждую функцию нужно заранее делать extensible.

---

# 5. L — Liskov Substitution

> **Подстановка реализации не должна нарушать ожидания и контракт потребителя.**

Например:

```go
type JobRepository interface {
	Get(id string) (Job, error)
}
```

У нас:

```text
MemoryJobRepository
PostgresJobRepository
```

Потребитель должен иметь возможность работать с обеими реализациями согласно одному контракту.

Плохой вариант:

```text
MemoryRepository
→ missing Job = ErrJobNotFound

PostgresRepository
→ missing Job = panic
```

Одна implementation изменила ожидаемую семантику.

Это нарушение substitutability.

---

# 6. I — Interface Segregation

> **Потребитель зависит только от нужного ему поведения.**

Плохой interface:

```go
type Repository interface {
	Add(...)
	Get(...)
	Delete(...)
	List(...)
	Export(...)
	RebuildIndex(...)
}
```

Если конкретному consumer нужен только `Get`, нет причины заставлять его зависеть от остальных методов.

Лучше:

```go
type JobReader interface {
	Get(id string) (Job, error)
}
```

и отдельно:

```go
type JobWriter interface {
	Add(job Job) error
}
```

Размер interface должен определяться **потребителем**, а не полнотой возможностей реализации.

---

# 7. D — Dependency Inversion

Плохая зависимость:

```text
JobService
    ↓
PostgreSQL
```

Application policy напрямую знает infrastructure detail.

Лучше:

```text
JobService
    ↓
JobRepository
    ↑
PostgresRepository
```

В Go:

```go
type JobRepository interface {
	Add(ctx context.Context, job Job) error
	Get(ctx context.Context, id string) (Job, error)
}
```

```go
type JobService struct {
	repo JobRepository
}
```

```go
func NewJobService(repo JobRepository) *JobService {
	return &JobService{
		repo: repo,
	}
}
```

Dependency создаётся снаружи через constructor injection.

Главная идея:

> **Policy зависит от стабильного поведения, а не от infrastructure detail.**

---

# 8. SOLID — краткая interview-модель

| Принцип | Главный вопрос                                                           |
| ------- | ------------------------------------------------------------------------ |
| SRP     | Какая причина изменения принадлежит компоненту?                          |
| OCP     | Как добавить новое поведение без лишнего переписывания стабильной части? |
| LSP     | Можно ли заменить implementation без нарушения ожиданий consumer?        |
| ISP     | Не зависит ли consumer от ненужного поведения?                           |
| DIP     | Не зависит ли policy от конкретной infrastructure detail?                |

---

# 9. SOLID не является чеклистом

Не нужно делать:

```text
interface
factory
strategy
adapter
abstraction
```

только потому, что они существуют.

Сначала:

```text
требование
 ↓
изменение
 ↓
граница
```

и только потом:

```text
минимальный механизм
```

---

# 10. KISS

> **Минимальная конструкция, достаточная для требования.**

Если нам нужен:

```go
type JobRepository interface {
	Get(...)
}
```

не нужно строить:

```text
RepositoryFactory
RepositoryRegistry
RepositoryStrategy
RepositoryManager
RepositoryProvider
```

если реального требования на это нет.

KISS управляет complexity.

---

# 11. DRY

> **Одно знание должно иметь одного владельца.**

DRY — не:

> «Нельзя написать два одинаковых фрагмента текста.»

Главный вопрос:

> **Дублируется ли одно и то же знание в нескольких местах?**

Например, если правило перехода Job из `pending` в `running` реализовано отдельно в нескольких местах, изменение правила создаёт synchronization risk.

Лучше:

```text
one rule
   ↓
one owner
```

---

# 12. YAGNI

> **Не создавай возможность до появления требования.**

Если сегодня JobFlow использует:

```text
Get
Add
```

не нужно заранее создавать:

```text
Delete
Archive
Export
Snapshot
BulkImport
```

пока нет соответствующих requirements.

YAGNI снижает:

```text
complexity
API surface
maintenance cost
test surface
```

---

# 13. Cohesion

> **Связанное поведение держи вместе.**

Высокая cohesion означает, что компонент содержит тесно связанное по смыслу поведение.

Хорошо:

```text
JobService
├── create
├── start
├── complete
└── fail
```

Плохо:

```text
JobService
├── create
├── renderPrometheus
├── encodeJSON
└── reconnectKafka
```

Во втором случае в одном месте смешаны разные причины изменения.

---

# 14. Coupling

> **Изменение одной части должно минимально затрагивать другие.**

Сильное coupling:

```text
JobService
   ↓
*sql.DB
   ↓
PostgreSQL API
```

При изменении storage нужно менять policy.

Слабее:

```text
JobService
   ↓
JobRepository
   ↑
PostgresRepository
```

Теперь concrete storage находится за boundary.

---

# 15. Cohesion + Coupling

Целевая модель:

> **Высокая cohesion внутри, низкая coupling снаружи.**

```text
component
┌─────────────────────┐
│ closely related     │
│ responsibilities    │
│                     │
└─────────────────────┘
     │   │   │
     ▼   ▼   ▼
   small external
   contracts
```

---

# 16. Cost of Change

![04B — Design Engineering](../posters/04B_DESIGN_ENGINEERING.svg)

Основной инженерный вопрос:

> **Что станет дорого менять?**

Стоимость изменения растёт, когда:

```text
одна причина
     ↓
затрагивает много компонентов
     ↓
через много зависимостей
```

Упрощённая модель:

```text
change
  ↓
blast radius
  ↓
verification cost
  ↓
maintenance cost
```

Хорошая boundary уменьшает blast radius.

---

# 17. Responsibility

Вопрос:

> **Что должно изменяться вместе?**

Связанное поведение с общей причиной изменения логично держать вместе.

Независимые оси изменения лучше разделять.

Пример:

```text
transport
application policy
persistence
messaging
observability
```

могут изменяться независимо.

---

# 18. Dependency Direction

Плохой вариант:

```text
APPLICATION
     ↓
PostgreSQL
```

Лучше:

```text
APPLICATION
     ↓
CONTRACT
     ↑
PostgreSQL ADAPTER
```

Dependency direction должна защищать policy от инфраструктурной детали.

---

# 19. Consumer-side Interface

Для JobFlow:

```go
type JobRepository interface {
	Add(ctx context.Context, job Job) error
	Get(ctx context.Context, id string) (Job, error)
}
```

Потребитель определяет, какое поведение ему нужно.

Concrete adapter удовлетворяет контракт автоматически.

Это особенно естественно для Go благодаря implicit interfaces.

---

# 20. Когда interface не нужен

Interface не является обязательной частью хорошего Go-кода.

Не нужно создавать abstraction только потому, что:

* код «должен быть тестируемым»;
* существует правило SOLID;
* есть возможность заменить реализацию когда-нибудь;
* так выглядит «правильная архитектура».

Нужна реальная boundary.

Хорошие причины:

```text
dependency isolation
alternative implementation
clear consumer contract
independent testing
separate responsibility
```

Если boundary не существует, abstraction может только увеличить сложность.

---

# 21. Практика: плохой JobFlow

Рассмотрим намеренно плохой handler:

```go
func CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateJobRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.Type == "" {
		http.Error(w, "missing type", http.StatusBadRequest)
		return
	}

	db := postgres.Open(...)
	defer db.Close()

	job := Job{
		ID:      uuid.NewString(),
		Type:    req.Type,
		Payload: req.Payload,
		Status:  StatusPending,
	}

	if _, err := db.ExecContext(
		r.Context(),
		"INSERT INTO jobs ...",
	); err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	kafka.Publish("job.created", job)

	log.Printf("created job %s", job.ID)

	json.NewEncoder(w).Encode(job)
}
```

---

# 22. Диагностика

Не начинай с:

> «Какой SOLID principle здесь нарушен?»

Сначала перечисли причины изменения.

```text
HTTP API
validation
business rules
PostgreSQL
Kafka
logging
```

Затем выдели:

```text
transport
policy
storage
messaging
observability
```

---

# 23. Refactoring Boundary

Целевая модель:

```text
HTTP
  ↓
JobService
  ├── JobRepository
  └── EventPublisher
```

Concrete infrastructure:

```text
JobRepository
    ↑
PostgresRepository
```

```text
EventPublisher
    ↑
KafkaPublisher
```

---

# 24. JobRepository

Минимальный contract:

```go
type JobRepository interface {
	Add(ctx context.Context, job Job) error
	Get(ctx context.Context, id string) (Job, error)
}
```

Consumer использует только необходимое поведение.

---

# 25. EventPublisher

Отдельная инфраструктурная ответственность:

```go
type EventPublisher interface {
	Publish(ctx context.Context, event JobCreated) error
}
```

Application service:

```go
type JobService struct {
	repo      JobRepository
	publisher EventPublisher
}
```

Dependencies:

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

Теперь Service не знает:

```text
PostgreSQL driver
Kafka client
connection string
topic implementation
```

---

# 26. Почему не один Infrastructure interface?

Плохой пример:

```go
type Infrastructure interface {
	SaveJob(...)
	PublishKafka(...)
	ReadConfig(...)
	WriteLog(...)
	GetRedis(...)
	SendHTTP(...)
}
```

Проблемы:

```text
large dependency surface
weak cohesion
strong coupling
harder substitution
unrelated responsibilities
```

Лучше несколько минимальных contracts:

```text
JobRepository
EventPublisher
JobCache
```

Каждый принадлежит своему consumer.

---

# 27. Практическое упражнение

Выполни refactoring плохого handler.

Нужно получить:

```text
HTTP transport
        ↓
JobService
       / \
      /   \
Repository Publisher
```

Условия:

* HTTP handler не открывает DB;
* HTTP handler не вызывает Kafka client;
* JobService не знает PostgreSQL;
* JobService не знает Kafka;
* concrete dependencies передаются через constructor;
* behavior должен остаться прежним.

---

# 28. Проверка refactoring

После изменения кода необходимо проверить:

```text
создание Job
validation
repository failure
publisher failure
error semantics
```

Запустить:

```bash
go test ./...
```

Если код содержит concurrent paths:

```bash
go test -race ./...
```

Главный принцип:

> **Архитектурный refactoring должен сохранять observable behavior, пока изменение поведения не является частью требования.**

---

# 29. Change Scenarios

Используй boundary как инструмент анализа изменений.

### PostgreSQL → другая БД

Ожидаемый blast radius:

```text
repository adapter
wiring
```

Policy должна оставаться стабильной.

### HTTP → gRPC

Ожидаемый blast radius:

```text
transport adapter
mapping
wiring
```

Application contract остаётся стабильным.

### Kafka → другая messaging system

Ожидаемый blast radius:

```text
publisher adapter
wiring
```

Business policy по возможности не меняется.

### Изменение business rule

Например:

```text
pending → running
```

получает новое условие.

Ожидаемый blast radius:

```text
domain/application policy
tests
```

---

# 30. OCP на практике

Допустим, есть:

```text
JobProcessor
```

и нам нужно добавить:

```text
EmailProcessor
WebhookProcessor
```

Хорошая boundary позволяет добавлять новую implementation, не переписывая весь стабильный pipeline.

Но сначала должен существовать реальный extension point.

> **OCP — не требование создать abstraction заранее.**

---

# 31. LSP на практике

Если две реализации соответствуют:

```go
type JobRepository interface {
	Get(id string) (Job, error)
}
```

обе должны сохранять semantic expectations consumer.

Например:

```text
missing Job
→ ErrJobNotFound
```

должно иметь одинаковую meaning across implementations.

Нельзя считать implementation заменяемой только потому, что compiler принял interface.

Нужно проверить **behavioral compatibility**.

---

# 32. Финальная mental model

```text
                REQUIREMENT
                     ↓
               CHANGE COST
                     ↓
                 BOUNDARY
                /    |    \
               /     |     \
         POLICY    STATE   EFFECT
            │        │       │
            ▼        ▼       ▼
         CONTRACT  OWNER   ADAPTER
```

SOLID помогает описывать отдельные проблемы:

```text
S → responsibility
O → extension
L → substitution contract
I → consumer surface
D → dependency direction
```

KISS / DRY / YAGNI помогают управлять complexity:

```text
KISS → минимальная конструкция
DRY  → один владелец знания
YAGNI → не строить без требования
```

И всё это служит одной цели:

> **Сделать изменение локальным, зависимости ясными, а сложность оправданной.**

---

# 33. Checkpoint

### Что такое SOLID?

<details><summary>Ответ</summary>Пять принципов проектирования, помогающих управлять ответственностями, изменениями, контрактами и зависимостями.</details>

### Что такое SRP?

<details><summary>Ответ</summary>Одна ось ответственности и одна основная причина изменения для компонента.</details>

### Что такое OCP?

<details><summary>Ответ</summary>Возможность расширять стабильную policy новым поведением с ограниченным переписыванием существующего ядра.</details>

### Что такое LSP?

<details><summary>Ответ</summary>Подстановка implementation не нарушает ожиданий и контракта consumer.</details>

### Что такое ISP?

<details><summary>Ответ</summary>Consumer зависит только от нужного ему поведения.</details>

### Что такое DIP?

<details><summary>Ответ</summary>Policy зависит от стабильного контракта, а не от concrete infrastructure detail.</details>

### Что такое cohesion?

<details><summary>Ответ</summary>Степень связанности поведения внутри компонента.</details>

### Что такое coupling?

<details><summary>Ответ</summary>Степень зависимости компонентов друг от друга и потенциальный blast radius изменений.</details>

### Когда нужен interface?

<details><summary>Ответ</summary>Когда существует реальная boundary, которая выигрывает от contract, substitution, dependency isolation или независимого тестирования.</details>

### Почему большой interface опасен?

<details><summary>Ответ</summary>Он увеличивает dependency surface и заставляет consumers зависеть от ненужного поведения.</details>

---

# 34. Домашнее задание — JobFlow v3

## Часть A — Refactoring

Перестроить JobFlow:

```text
TRANSPORT
    ↓
APPLICATION
    ↓
small contracts
    ↓
INFRASTRUCTURE
```

Добавить минимальные contracts:

```go
type JobRepository interface {
	Add(ctx context.Context, job Job) error
	Get(ctx context.Context, id string) (Job, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event JobCreated) error
}
```

Concrete implementations внедряются через constructors.

---

## Часть B — Новое поведение

Добавить второй processor.

Например:

```text
EmailProcessor
WebhookProcessor
```

Архитектура должна позволять добавить implementation без постоянного переписывания стабильной policy.

---

## Часть C — Alternative implementation

Добавить вторую реализацию repository для тестов.

Например:

```text
MemoryJobRepository
FakeJobRepository
```

Проверить, что `JobService` не зависит от concrete implementation.

---

## Часть D — Design Note

Написать не более 10 строк:

1. Почему `JobRepository` определён на стороне consumer?
2. Почему `EventPublisher` является отдельным interface?
3. Почему не нужен единый `Infrastructure` interface?
4. Какие изменения теперь локализованы лучше?
5. Какой принцип SOLID лучше всего объясняет каждую выбранную boundary?

---

# 35. Verification

Запустить:

```bash
go fmt ./...
go test ./...
go test -race ./...
go vet ./...
```

Проверить:

```text
existing behavior preserved
new processor works
alternative repository works
business policy independent from infrastructure
no race introduced
```

---

# 36. Acceptance Criteria

## PASS

Ученик способен:

* объяснить SOLID без подсказки;
* дать Go-specific пример каждого принципа;
* определить причины изменения компонента;
* назвать coupling и cohesion;
* определить consumer-side interface;
* обосновать dependency direction;
* не создавать abstraction без требования;
* заменить infrastructure implementation;
* добавить новую implementation с ограниченным изменением стабильной части;
* доказать сохранение behavior tests.

## FAIL

Решение не считается освоенным, если аргумент звучит только так:

> «Нам нужен interface, потому что DIP.»

Необходимо уметь ответить:

> **Какую конкретную dependency или стоимость изменения мы этим уменьшаем?**

---

# 37. Следующая тема

После этого урока JobFlow имеет явные boundaries:

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

Следующий вопрос:

> **Как сохранить корректное состояние этой системы, когда данные становятся durable, операции выполняются конкурентно, а несколько workers могут одновременно менять одну Job?**
