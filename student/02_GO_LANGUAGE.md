reply:"student/02_GO_LANGUAGE.md" v1.0.0 2026-08-19T04:45:00Z Topic: Delivery Summary: Дидактический материал ученика — Go Language и первая версия JobFlow

# Go Backend Combat — 01. LANGUAGE

**Version:** 1.0.0
**Topic:** Go Language
**Format:** учебный материал ученика
**Project:** JobFlow
**Poster:** [02 — Go Language](../posters/02_LANGUAGE.svg)

---

# 1. Цель

Собрать рабочую mental model Go и сразу применить её в JobFlow.

Не нужно начинать с запоминания синтаксиса.

Основная модель:

```text
TYPES
  ↓
FUNCTIONS
  ↓
METHODS
  ↓
INTERFACES
  ↓
ERRORS
  ↓
TESTING
```

![02 — Go Language](../posters/02_LANGUAGE.svg)

К концу занятия должна существовать первая рабочая версия JobFlow:

```text
Job
 ↓
JobStatus
 ↓
methods
 ↓
validation
 ↓
repository
 ↓
service
 ↓
tests
```

---

# 2. Go как модель языка

Go особенно хорошо раскрывается через несколько простых идей:

```text
explicit state
explicit behavior
small contracts
explicit errors
composition
built-in testing
```

Вместо:

```text
изучить много syntax
```

используем:

```text
задача
 ↓
модель
 ↓
Go construction
 ↓
рабочий код
 ↓
test
```

---

# 3. Types

## Struct

`struct` объединяет связанные данные в один пользовательский тип.

Для Job:

```go
type Job struct {
	ID      string
	Type    string
	Payload []byte
}
```

Смысл:

> **Состояние сущности выражается явным типом.**

### Почему не отдельные переменные?

Вместо:

```go
jobID string
jobType string
payload []byte
```

получаем:

```go
job := Job{
	ID:      "job-1",
	Type:    "email",
	Payload: nil,
}
```

Теперь Job является единым значением и его можно передавать между функциями, хранить в коллекциях и расширять поведением.

---

# 4. Static Typing и Type Inference

Go имеет static typing.

Обе записи корректны:

```go
var count int = 10
```

```go
count := 10
```

Во втором случае compiler выводит тип:

```text
10
↓
int
```

Type inference уменьшает количество текста, но не отменяет статический тип.

Например:

```go
count := 10
count = "ten"
```

не компилируется.

Mental model:

```text
explicit type
      │
      └── compiler knows type

type inference
      │
      └── compiler derives type
```

---

# 5. Zero Value

Каждый Go type имеет zero value.

```go
var name string
var count int
var active bool
var payload []byte
var job *Job
```

Получаем:

```text
string  → ""
int     → 0
bool    → false
slice   → nil
pointer → nil
```

Zero value — часть дизайна Go API.

Хороший type имеет понятное состояние по умолчанию.

---

# 6. Domain Types

Не всегда полезно использовать primitive напрямую.

Вместо:

```go
Status string
```

можно определить:

```go
type JobStatus string
```

и допустимые значения:

```go
const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)
```

Теперь API выражает смысл:

```go
func SetStatus(status JobStatus)
```

вместо:

```go
func SetStatus(status string)
```

Принцип:

> **Хороший тип делает неправильное состояние труднее выразить.**

---

# 7. Job Model

Итоговая модель:

```go
type Job struct {
	ID        string
	Type      string
	Payload   []byte
	Status    JobStatus
	CreatedAt time.Time
}
```

Здесь:

```text
ID        → identity
Type      → тип задачи
Payload   → данные задачи
Status    → текущее состояние
CreatedAt → время создания
```

---

# 8. Functions

Function описывает действие и его контракт.

Пример:

```go
func NewJob(
	id string,
	jobType string,
	payload []byte,
) Job {
	return Job{
		ID:        id,
		Type:      jobType,
		Payload:   payload,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}
```

Её контракт:

```text
input
  ↓
Job
```

Go также позволяет возвращать несколько значений:

```go
func FindJob(
	jobs []Job,
	id string,
) (Job, bool)
```

Например:

```go
job, ok := FindJob(jobs, "job-1")
```

---

# 9. Methods

Method — function с receiver.

```go
func (j Job) IsFinished() bool {
	return j.Status == StatusCompleted ||
		j.Status == StatusFailed
}
```

Использование:

```go
if job.IsFinished() {
	// ...
}
```

Смысл:

> **Поведение находится рядом с типом, которому оно принадлежит.**

---

# 10. Value Receiver и Pointer Receiver

## Value receiver

```go
func (j Job) Complete() {
	j.Status = StatusCompleted
}
```

`j` — значение receiver.

Это не меняет исходный объект `job`.

```text
job
 │
 └── copy → j
              ↓
           изменяется
```

## Pointer receiver

```go
func (j *Job) Complete() {
	j.Status = StatusCompleted
}
```

Теперь method работает с объектом через pointer и может менять его состояние.

```text
job
 │
 └────────→ j
              ↓
         изменяет job
```

Практическое правило:

> **Если method должен менять state объекта, pointer receiver обычно естественнее.**

Pointer в Go не означает ручное управление памятью. Память управляется runtime.

---

# 11. Slice

Slice представляет последовательность элементов.

```go
jobs := []Job{
	NewJob("job-1", "email", nil),
	NewJob("job-2", "sms", nil),
	NewJob("job-3", "email", nil),
}
```

Основные свойства:

```go
len(jobs)
cap(jobs)
```

`len` — количество доступных элементов.

`cap` — capacity backing array от текущего начала slice.

---

## Append

```go
jobs = append(
	jobs,
	NewJob("job-4", "webhook", nil),
)
```

Почему результат нужно присваивать обратно?

Потому что `append` при нехватке capacity может создать новый backing array и вернуть новый slice.

---

# 12. Slice Practice

Реализовать:

```go
func FindJob(
	jobs []Job,
	id string,
) *Job
```

Ожидаемая форма:

```go
for i := range jobs {
	if jobs[i].ID == id {
		return &jobs[i]
	}
}

return nil
```

Почему индекс?

Потому что `jobs[i]` — элемент самого slice.

---

# 13. Map

Когда основной access pattern:

```text
JobID → Job
```

естественным выбором становится map:

```go
jobsByID := map[string]Job{}
```

Добавление:

```go
jobsByID[job.ID] = job
```

Получение:

```go
job, ok := jobsByID["job-1"]
```

---

# 14. Map Lookup

Форма:

```go
value, ok := m[key]
```

`ok` отвечает на вопрос:

> Существует ли ключ?

Это важно, потому что отсутствующий key возвращает zero value значения.

Например:

```go
jobs := map[string]Job{}

job, ok := jobs["missing"]
```

Результат:

```text
job = zero value Job
ok  = false
```

---

# 15. Nil Map

Zero value map:

```go
var jobs map[string]Job
```

равен:

```go
nil
```

Чтение допустимо:

```go
job, ok := jobs["job-1"]
```

Запись вызовет panic:

```go
jobs["job-1"] = job
```

Для записи используется:

```go
jobs := make(map[string]Job)
```

---

# 16. Memory Repository

Создать:

```go
type MemoryJobRepository struct {
	jobs map[string]Job
}
```

Constructor:

```go
func NewMemoryJobRepository() *MemoryJobRepository {
	return &MemoryJobRepository{
		jobs: make(map[string]Job),
	}
}
```

---

# 17. Errors

В Go failure обычно является частью возвращаемого значения.

Например:

```go
job, err := repo.Get(id)

if err != nil {
	return err
}
```

Создадим стабильную error identity:

```go
var ErrJobNotFound = errors.New("job not found")
```

Теперь:

```go
func (r *MemoryJobRepository) Get(
	id string,
) (Job, error) {
	job, ok := r.jobs[id]
	if !ok {
		return Job{}, ErrJobNotFound
	}

	return job, nil
}
```

---

# 18. Error Wrapping

Можно добавить контекст:

```go
return Job{}, fmt.Errorf(
	"get job %s: %w",
	id,
	ErrJobNotFound,
)
```

Проверка:

```go
if errors.Is(err, ErrJobNotFound) {
	// not found
}
```

Разделение:

```text
error message
    ↓
читается человеком

error identity
    ↓
используется программой
```

Принцип:

> **Текст ошибки объясняет проблему. Error identity позволяет принять решение.**

---

# 19. Validation

Создать:

```go
var ErrInvalidJob = errors.New("invalid job")
```

Например:

```go
func ValidateJob(job Job) error {
	if job.ID == "" {
		return fmt.Errorf("job id: %w", ErrInvalidJob)
	}

	if job.Type == "" {
		return fmt.Errorf("job type: %w", ErrInvalidJob)
	}

	return nil
}
```

Использование:

```go
if err := ValidateJob(job); err != nil {
	return err
}
```

В Go error checking обычно находится рядом с операцией, которая могла завершиться ошибкой.

---

# 20. Interface

Когда storage должен быть заменяемым, application не должен зависеть от конкретной реализации.

Определяем минимальный contract:

```go
type JobRepository interface {
	Add(Job) error
	Get(string) (Job, error)
}
```

Memory repository удовлетворяет этому interface автоматически.

Явного:

```text
implements
```

в Go нет.

Тип соответствует interface, если его method set соответствует контракту.

---

# 21. Consumer-side Interface

Не нужно описывать все возможности repository.

Если потребителю нужны только:

```go
Add(Job) error
Get(string) (Job, error)
```

не следует без требования добавлять:

```go
Delete(...)
List(...)
Update(...)
FindByType(...)
```

Принцип:

> **Interface описывает необходимое поведение потребителя.**

---

# 22. Dependency Injection

Application component:

```go
type JobService struct {
	repo JobRepository
}
```

Constructor:

```go
func NewJobService(
	repo JobRepository,
) *JobService {
	return &JobService{
		repo: repo,
	}
}
```

Зависимость передаётся снаружи.

```text
MemoryJobRepository
        │
        ▼
    JobService
```

Позже:

```text
PostgresJobRepository
        │
        ▼
    JobService
```

При этом application contract остаётся тем же.

---

# 23. Composition

Для разных видов Job можно использовать composition.

```go
type EmailJob struct {
	Job
	To string
}
```

```go
type WebhookJob struct {
	Job
	URL string
}
```

Модель:

```text
EmailJob
 ├── Job
 └── To

WebhookJob
 ├── Job
 └── URL
```

Принцип:

> **Композиция связывает небольшие типы и поведение без необходимости строить глубокую иерархию наследования.**

---

# 24. Defer

`defer` откладывает выполнение действия до выхода из текущей функции.

Пример:

```go
func demo() {
	defer fmt.Println("third")
	defer fmt.Println("second")

	fmt.Println("first")
}
```

Результат:

```text
first
second
third
```

Несколько `defer` выполняются LIFO.

На практике `defer` часто используется для cleanup:

```text
unlock
close
cancel
release
```

---

# 25. Testing

Go имеет встроенный testing toolchain.

Пример:

```go
func TestJobComplete(t *testing.T) {
	job := NewJob("job-1", "email", nil)

	job.Complete()

	if job.Status != StatusCompleted {
		t.Fatalf(
			"expected %q, got %q",
			StatusCompleted,
			job.Status,
		)
	}
}
```

Запуск:

```bash
go test ./...
```

Главный объект тестирования:

> **observable behavior**

а не внутренняя реализация.

---

# 26. Table-driven Tests

Для похожих cases:

```go
tests := []struct {
	name    string
	job     Job
	wantErr bool
}{
	{
		name: "valid",
		job: Job{
			ID:   "job-1",
			Type: "email",
		},
		wantErr: false,
	},
	{
		name: "missing id",
		job: Job{
			Type: "email",
		},
		wantErr: true,
	},
	{
		name: "missing type",
		job: Job{
			ID: "job-1",
		},
		wantErr: true,
	},
}
```

Запуск:

```go
for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		err := ValidateJob(tt.job)

		if (err != nil) != tt.wantErr {
			t.Fatalf(
				"error = %v, wantErr = %v",
				err,
				tt.wantErr,
			)
		}
	})
}
```

Преимущество:

```text
cases = data
logic = one test body
```

---

# 27. JobFlow v0

К концу первой версии модель должна выглядеть примерно так:

```text
Job
├── JobStatus
├── NewJob
├── Start
├── Complete
├── Fail
└── IsFinished

Validation
└── ValidateJob

Repository
├── JobRepository
└── MemoryJobRepository

Service
└── JobService

Tests
├── Job behavior
├── validation
└── repository
```

Lifecycle:

```text
pending
   │
   ▼
running
  /   \
 ▼     ▼
done  failed
```

---

# 28. Практическое задание

Самостоятельно реализовать первую рабочую версию JobFlow.

Необходимо иметь:

```go
func NewJob(
	id string,
	jobType string,
	payload []byte,
) Job
```

```go
func (j *Job) Start()
```

```go
func (j *Job) Complete()
```

```go
func (j *Job) Fail()
```

```go
func (j Job) IsFinished() bool
```

```go
func ValidateJob(job Job) error
```

```go
type JobRepository interface {
	Add(Job) error
	Get(string) (Job, error)
}
```

```go
type MemoryJobRepository struct {
	jobs map[string]Job
}
```

```go
type JobService struct {
	repo JobRepository
}
```

---

# 29. Проверка понимания

Ответь без IDE.

### Что такое `struct`?

<details><summary>Ответ</summary>Пользовательский тип, объединяющий связанные данные.</details>

### Что такое method?

<details><summary>Ответ</summary>Функция с receiver, выражающая поведение типа.</details>

### Когда нужен pointer receiver?

<details><summary>Ответ</summary>Когда метод должен изменять состояние объекта через указатель либо pointer semantics оправданы другими причинами.</details>

### Что такое slice?

<details><summary>Ответ</summary>Представление последовательности элементов с length и capacity.</details>

### Что такое map?

<details><summary>Ответ</summary>Структура доступа по ключу, например `JobID → Job`.</details>

### Что означает `value, ok := m[key]`?

<details><summary>Ответ</summary>`value` — найденное значение или zero value; `ok` показывает наличие ключа.</details>

### Что такое zero value?

<details><summary>Ответ</summary>Определённое Go начальное значение типа.</details>

### Что такое interface?

<details><summary>Ответ</summary>Контракт поведения, которому тип соответствует через method set.</details>

### Что такое error?

<details><summary>Ответ</summary>Значение, описывающее неуспешный результат операции и являющееся частью её контракта.</details>

### Зачем `%w`?

<details><summary>Ответ</summary>Для wrapping error с сохранением исходной identity.</details>

### Зачем tests?

<details><summary>Ответ</summary>Для воспроизводимого evidence ожидаемого поведения.</details>

---

# 30. Домашнее задание — JobFlow v1

Расширить repository:

```go
type JobRepository interface {
	Add(Job) error
	Get(string) (Job, error)
	Delete(string) error
	List() []Job
}
```

Добавить:

```text
Delete
List
FindByType
```

## Требуемое поведение

### Get

```text
existing → Job + nil
missing  → zero Job + ErrJobNotFound
```

### Delete

```text
existing → nil
missing  → ErrJobNotFound
```

### List

> Вернуть все Job.

### FindByType

> Вернуть все Job указанного типа.

---

# 31. Tests

Добавить tests для:

```text
create
get
not found
delete
delete missing
list
find by type
start
complete
fail
validation
invalid state transition
```

Минимально выполнить:

```bash
go fmt ./...
go test ./...
go vet ./...
```

---

# 32. PASS / FAIL

## PASS

Домашнее задание считается выполненным, если:

* новые методы добавлены в interface осознанно;
* MemoryRepository реализует их;
* tests проверяют новое поведение;
* существующие tests продолжают проходить;
* ученик может объяснить value/pointer semantics;
* ученик понимает `slice`, `map`, `error`;
* ученик понимает implicit interface implementation;
* `go test ./...` проходит;
* `go vet ./...` проходит.

## FAIL

Работа не считается законченной, если решение работает только после копирования готового кода и его семантика не может быть объяснена.

---

# 33. Следующая проблема

Сейчас JobFlow предполагает последовательное выполнение.

В следующем занятии появляется:

```text
1000 jobs
   │
   ├── worker
   ├── worker
   ├── worker
   └── ...
```

Новый вопрос:

> **Что произойдёт, если несколько goroutines одновременно начнут менять одно состояние?**
