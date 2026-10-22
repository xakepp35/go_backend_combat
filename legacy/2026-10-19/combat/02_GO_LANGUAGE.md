<!-- artifact:combat/homework/02_GO_LANGUAGE_BUG_HUNT.md v1.0.0 2026-09-23T17:02:00Z Topic:Delivery Summary:Практика семантики Go через диагностику и исправление типовых ошибок -->

# Go Backend Combat — 02. GO LANGUAGE

# Homework: Language Bug Hunt

## Цель

На уроке вы уже познакомились с основными конструкциями языка Go.

Теперь задача меняется.

Вам не нужно ещё раз демонстрировать, что вы умеете написать:

```go
for
```

или:

```go
append(...)
```

Нужно показать, что вы понимаете, **что именно делает написанный Go-код**.

Во всех заданиях ниже код намеренно содержит ошибку.

Часть программ:

- компилируется;
- проходит очевидный happy path;
- выглядит вполне разумно;

но нарушает заявленный контракт.

Ваша задача:

```text
READ
  ↓
PREDICT
  ↓
RUN
  ↓
EXPLAIN
  ↓
FIX
  ↓
TEST
```

---

# Общие правила

Каждую задачу разместите отдельно:

```text
homework/
├── task01/
├── task02/
├── task03/
├── task04/
├── task05/
├── task06/
├── task07/
├── task08/
├── task09/
└── task10/
```

Для каждой задачи:

1. прочитайте исходный код;
2. до запуска сформулируйте ожидаемое поведение;
3. запустите код;
4. найдите причину ошибки;
5. исправьте минимально необходимую часть;
6. добавьте test, фиксирующий требуемое поведение.

В начале решения оставьте:

```go
// BUG:
// кратко описать причину
//
// FIX:
// кратко описать изменение
//
// GUARANTEE:
// какое свойство теперь гарантируется
```

Не переписывайте задачу полностью, если достаточно исправить несколько строк.

---

# Task 01 — Add без добавления

Есть функция добавления Job:

```go
package task01

type Job struct {
	ID string
}

func AddJob(
	jobs []Job,
	job Job,
) {
	jobs = append(jobs, job)
}
```

Использование:

```go
jobs := []Job{
	{ID: "job-1"},
}

AddJob(
	jobs,
	Job{ID: "job-2"},
)

fmt.Println(len(jobs))
```

Автор ожидает:

```text
2
```

Но получает:

```text
1
```

## Требование

После исправления следующий код должен работать:

```go
jobs := []Job{
	{ID: "job-1"},
}

jobs = AddJob(
	jobs,
	Job{ID: "job-2"},
)
```

После вызова:

```text
len(jobs) == 2
jobs[0].ID == "job-1"
jobs[1].ID == "job-2"
```

## Test

Добавьте test минимум для:

```text
empty slice + one Job
existing slice + one Job
multiple sequential additions
```

## После решения

Объясните своими словами:

> Что именно передаётся в функцию, когда аргумент имеет тип `[]Job`?

---

# Task 02 — Clone, который ничего не клонирует

Разработчик написал функцию копирования:

```go
package task02

type Job struct {
	ID     string
	Status string
}

func CloneJobs(jobs []Job) []Job {
	clone := jobs
	return clone
}
```

Использование:

```go
source := []Job{
	{
		ID:     "job-1",
		Status: "pending",
	},
}

clone := CloneJobs(source)

clone[0].Status = "completed"

fmt.Println(source[0].Status)
```

Автор ожидает:

```text
pending
```

Но исходный slice также изменяется.

## Контракт

`CloneJobs` должна возвращать независимый slice элементов.

После:

```go
clone := CloneJobs(source)
```

изменение:

```go
clone[0].Status = "completed"
```

не должно изменять:

```go
source[0].Status
```

---

## Дополнительные cases

Проверьте:

```text
nil slice
empty slice
slice из нескольких элементов
```

Не требуется сохранять:

```text
nil
```

и:

```text
empty non-nil slice
```

как разные состояния, если вы явно не выбрали такой contract.

Главное требование:

> Изменение элементов clone не изменяет элементы source.

---

# Task 03 — Указатель не на тот Job

Необходимо найти Job внутри slice и вернуть pointer на него.

```go
package task03

type Job struct {
	ID     string
	Status string
}

func FindJob(
	jobs []Job,
	id string,
) *Job {
	for _, job := range jobs {
		if job.ID == id {
			return &job
		}
	}

	return nil
}
```

Использование:

```go
jobs := []Job{
	{
		ID:     "job-1",
		Status: "pending",
	},
}

job := FindJob(
	jobs,
	"job-1",
)

job.Status = "completed"
```

Автор ожидает:

```text
jobs[0].Status == "completed"
```

Но это не выполняется.

## Требование

После исправления pointer, возвращённый `FindJob`, должен указывать именно на элемент исходного slice.

Должно выполняться:

```go
job := FindJob(jobs, "job-1")

job.Status = "completed"
```

и затем:

```text
jobs[0].Status == "completed"
```

Для неизвестного `ID`:

```text
result == nil
```

---

## Вопрос

Объясните разницу между:

```go
job
```

в:

```go
for _, job := range jobs
```

и:

```go
jobs[i]
```

с точки зрения того, **какой объект вы изменяете**.

---

# Task 04 — Repository, которого нет, но error есть

Есть custom error:

```go
package task04

type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return "invalid field: " + e.Field
}
```

Функция validation:

```go
func Validate(id string) error {
	var validationErr *ValidationError

	if id == "" {
		validationErr = &ValidationError{
			Field: "id",
		}
	}

	return validationErr
}
```

Использование:

```go
err := Validate("job-1")

if err != nil {
	fmt.Println("validation failed:", err)
	return
}

fmt.Println("valid")
```

Для корректного `ID` автор ожидает:

```text
valid
```

Но программа ведёт себя иначе.

## Контракт

Для:

```go
Validate("")
```

должна возвращаться ошибка.

Для:

```go
Validate("job-1")
```

должно выполняться:

```text
err == nil
```

Не:

```text
пустой ValidationError
```

не:

```text
interface, внутри которого nil pointer
```

а именно:

```go
nil
```

---

## Tests

Проверьте:

```text
empty ID → error
valid ID → nil
```

Для ошибки используйте `errors.As`.

---

## После решения

Ответьте:

> Может ли interface value быть `!= nil`, если concrete pointer внутри него равен `nil`?

---

# Task 05 — Error потерял identity

Repository использует публичную error identity:

```go
package task05

import (
	"errors"
	"fmt"
)

var ErrJobNotFound = errors.New("job not found")
```

Но service добавляет контекст так:

```go
func LoadJob() error {
	err := ErrJobNotFound

	return fmt.Errorf(
		"load job: %v",
		err,
	)
}
```

Caller:

```go
err := LoadJob()

if errors.Is(err, ErrJobNotFound) {
	fmt.Println("not found")
}
```

Автор ожидает:

```text
not found
```

Но условие не выполняется.

## Требование

Service имеет право добавлять человекочитаемый context:

```text
load job: job not found
```

При этом machine-readable identity должна сохраняться:

```go
errors.Is(
	err,
	ErrJobNotFound,
) == true
```

---

## Расширение

Добавьте:

```go
type ValidationError struct {
	Field string
}
```

и функцию:

```go
func ValidateJob(...) error
```

Caller должен иметь возможность извлечь ошибку через:

```go
errors.As(...)
```

даже если между domain error и caller существует дополнительный error context.

---

# Task 06 — Ошибка исчезла из-за `:=`

Есть service:

```go
package task06

import "errors"

var ErrSave = errors.New("save failed")

type Repository interface {
	Add() error
}

func Save(
	repository Repository,
	enabled bool,
) error {
	var err error

	if enabled {
		err := repository.Add()

		if err != nil {
			// здесь ошибка действительно существует
		}
	}

	return err
}
```

Repository возвращает:

```go
ErrSave
```

Но `Save` неожиданно возвращает:

```text
nil
```

## Требование

Если:

```text
enabled == false
```

должно возвращаться:

```text
nil
```

Если:

```text
enabled == true
```

и repository вернул:

```go
ErrSave
```

caller должен получить ошибку, для которой:

```go
errors.Is(err, ErrSave) == true
```

---

## Ограничение

Не добавляйте package-level state.

Исправление должно находиться внутри функции `Save`.

---

## После решения

Опишите:

> Какие две разные переменные с именем `err` существуют в исходной программе?

---

# Task 07 — Method существует, но interface его не видит

Есть contract:

```go
package task07

type Job struct {
	ID string
}

type Repository interface {
	Add(Job) error
}
```

Implementation:

```go
type MemoryRepository struct{}

func (r *MemoryRepository) Add(
	job Job,
) error {
	return nil
}
```

Composition:

```go
func NewRepository() MemoryRepository {
	return MemoryRepository{}
}
```

Использование:

```go
var repository Repository

repository = NewRepository()
```

Код не компилируется.

## Требование

Добейтесь корректной composition.

При этом:

- `Repository` оставьте interface;
- `Add` должен остаться method конкретного repository;
- не создавайте фиктивный второй `Add`;
- не используйте type assertion ради обхода compiler.

После исправления должно компилироваться:

```go
repository := NewRepository()

err := repository.Add(
	Job{ID: "job-1"},
)
```

---

## Добавьте compile-time assertion

```go
var _ Repository = (*MemoryRepository)(nil)
```

или эквивалентную проверку для выбранной вами конструкции.

---

## После решения

Объясните:

```text
method set T
```

и:

```text
method set *T
```

в контексте этой задачи.

---

# Task 08 — defer запомнил слишком рано

Необходимо вывести финальный status операции.

```go
package task08

import "fmt"

func Run() {
	status := "started"

	defer fmt.Println(
		"final status:",
		status,
	)

	status = "completed"
}
```

Автор ожидает:

```text
final status: completed
```

Но получает:

```text
final status: started
```

## Требование

Сохраните `defer`.

После исправления:

```go
Run()
```

должна печатать:

```text
final status: completed
```

---

## Вопрос

Когда вычисляются arguments функции в:

```go
defer f(x)
```

?

Когда выполняется сама `f`?

Это два разных момента.

---

# Task 09 ★ — Stringer уничтожает программу

Есть тип:

```go
package task09

import "fmt"

type Job struct {
	ID string
}

func (j Job) String() string {
	return fmt.Sprintf(
		"Job{%v}",
		j,
	)
}
```

Использование:

```go
job := Job{
	ID: "job-42",
}

fmt.Println(job)
```

Автор хочет получить:

```text
Job{job-42}
```

Но программа не может корректно сформировать строку.

## Задача

Найдите recursive call chain и исправьте `String()`.

После исправления:

```go
fmt.Sprint(job)
```

должно возвращать:

```text
Job{job-42}
```

---

## Ограничение

Не удаляйте:

```go
String() string
```

Метод должен продолжать реализовывать:

```go
fmt.Stringer
```

---

# Task 10 ★ — SafeCall, который ничего не сообщает

Необходимо безопасно вызвать callback.

Исходный код:

```go
package task10

import "fmt"

func SafeCall(
	fn func(),
) error {
	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Println(
				"recovered:",
				recovered,
			)
		}
	}()

	fn()

	return nil
}
```

Использование:

```go
err := SafeCall(func() {
	panic("boom")
})
```

Panic действительно перехватывается.

Но caller получает:

```text
err == nil
```

и не может программно узнать, что callback завершился аварийно.

## Требование

Объявите:

```go
type PanicError struct {
	Value any
}
```

`SafeCall` должна выполнять:

### Normal callback

```go
SafeCall(func() {})
```

результат:

```text
err == nil
```

### Panic callback

```go
SafeCall(func() {
	panic("boom")
})
```

результат:

```text
err != nil
```

и:

```go
errors.As(
	err,
	&PanicError{},
)
```

или корректный эквивалент должен позволять получить `PanicError`.

Сам panic не должен выйти наружу из `SafeCall`.

---

# Финальная задача — Language Review

После выполнения всех задач закройте исходные материалы.

Для каждого фрагмента ниже ответьте без запуска кода:

## A

```go
a := []int{1, 2, 3}
b := a
b[0] = 100

fmt.Println(a[0])
```

Что будет выведено и почему?

---

## B

```go
var p *int

var value any = p

fmt.Println(value == nil)
```

Что будет выведено и почему?

---

## C

```go
type T struct{}

func (*T) Run() {}

var x interface {
	Run()
}

x = T{}
```

Компилируется ли это?

Почему?

---

## D

```go
x := 10

defer fmt.Println(x)

x = 20
```

Что напечатает `defer`?

---

## E

```go
for _, value := range []int{1, 2, 3} {
	value *= 10
}
```

Изменится ли исходный slice?

Почему?

---

## F

```go
var err error

if true {
	err := errors.New("boom")
	_ = err
}

fmt.Println(err)
```

Что будет выведено?

---

# Acceptance

Домашняя работа считается выполненной, если:

```bash
go fmt ./...
go test ./...
go vet ./...
```

проходят.

Но этого недостаточно.

По каждой задаче ученик должен уметь объяснить:

```text
что именно было сломано
        ↓
почему Go ведёт себя именно так
        ↓
какая language semantic за этим стоит
        ↓
какое минимальное изменение исправило поведение
        ↓
какой test фиксирует гарантию
```

---

# Что должно остаться после домашки

Не набор правил:

```text
append надо так
pointer надо вот так
interface иногда странный
```

а более точная модель:

```text
VALUE
  │
  ├── копируется
  │
  └── имеет конкретный type

SLICE
  │
  ├── header передаётся как value
  └── backing array может быть shared

POINTER
  │
  └── указывает на конкретный объект

METHOD
  │
  └── принадлежит method set

INTERFACE
  │
  ├── dynamic type
  └── dynamic value

ERROR
  │
  └── interface + identity chain

SCOPE
  │
  └── identifier относится к конкретной declaration

DEFER
  │
  ├── регистрируется сейчас
  ├── arguments вычисляются сейчас
  └── call выполняется позже
```

Финальная мысль:

> **В Go мало синтаксиса. Сложность начинается не в том, как написать конструкцию, а в том, какое значение, память, тип и lifetime она реально представляет.**
