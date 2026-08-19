artifact:student/02_CONCURRENCY.md v1.0.0 2026-08-19T07:48:00Z Topic: Delivery Summary: Дидактический материал ученика — Go Concurrency, runtime model, worker pool и безопасный lifecycle

# Go Backend Combat — 02. CONCURRENCY

**Version:** 1.0.0
**Topic:** Go Concurrency
**Format:** учебный материал ученика
**Project:** JobFlow
**Duration:** 120 минут
**Prerequisite:** `01_LANGUAGE.md`

Плакаты:

* [03A — Concurrency Model](../posters/03A_CONCURRENCY_MODEL.svg)
* [03B — Concurrency Runtime](../posters/03B_CONCURRENCY_RUNTIME.svg)
* [03C — Concurrency Engineering](../posters/03C_CONCURRENCY_ENGINEERING.svg)

---

# 1. Цель занятия

Научиться строить конкурентную обработку Job с явными:

```text
ownership
lifecycle
synchronization
cancellation
resource limits
```

Основная модель:

```text
JOBS
  ↓
QUEUE
  ↓
WORKERS
  ↓
RESULTS
```

Но корректный concurrent backend требует большего:

```text
ownership
    ↓
concurrent execution
    ↓
communication
    ↓
synchronization
    ↓
bounded concurrency
    ↓
cancellation
    ↓
graceful shutdown
    ↓
verification
```

Главная идея:

> **Concurrency — это управление совместным выполнением, передачей данных, состоянием, ресурсами и временем жизни.**

---

# 2. JobFlow: зачем нужна concurrency

На предыдущем занятии JobFlow обрабатывал Job последовательно.

Допустим:

```text
1000 jobs
1 job = 100 ms
```

Последовательная обработка:

```text
1000 × 100 ms = 100 000 ms
```

то есть около:

```text
100 секунд
```

Если Job независимы и выполнение можно перекрывать, появляется возможность использовать concurrent processing.

Целевая модель:

```text
                JOBS
                  │
                  ▼
                QUEUE
                  │
        ┌─────────┼─────────┐
        ▼         ▼         ▼
      worker    worker    worker
        │         │         │
        └─────────┼─────────┘
                  ▼
               RESULTS
```

Но количество входных Job не должно автоматически определять количество concurrent resources.

---

# 3. Concurrency Model

![03A — Concurrency Model](../posters/03A_CONCURRENCY_MODEL.svg)

Плакат задаёт базовые вопросы concurrent-программы:

```text
OWNERSHIP
GOROUTINE
CHANNEL
SELECT
SYNCHRONIZATION
LIFECYCLE
```

---

## 3.1 Ownership

Concurrent operation должна иметь понятного владельца.

Owner отвечает за:

```text
start
resources
lifecycle
stop
```

Плохая модель:

```text
go something()
```

и неизвестно:

```text
кто запустил?
зачем?
когда остановить?
кто ждёт завершения?
```

Хорошая модель явно отвечает на эти вопросы.

> **Каждая долгоживущая goroutine должна иметь owner и stop condition.**

---

# 4. Goroutine

Запуск:

```go
go process(job)
```

создаёт goroutine — единицу concurrent execution, управляемую Go runtime.

Это не то же самое, что:

```text
1 goroutine = 1 OS thread
```

Runtime самостоятельно планирует goroutines поверх OS threads.

Поэтому большое количество goroutines возможно, но это не означает бесконтрольное создание такого же количества OS threads.

---

## Проверка понимания

Что происходит после:

```go
go process(job)
```

<details><summary>Ответ</summary>Вызов запускается как goroutine и может выполняться concurrently с вызывающей goroutine.</details>

Является ли goroutine OS thread?

<details><summary>Ответ</summary>Нет. Goroutine управляется Go runtime и выполняется на OS threads, которыми управляет runtime.</details>

---

# 5. Lifecycle первой goroutine

Пример:

```go
func main() {
	go process(job)

	fmt.Println("main finished")
}
```

Проблема:

```text
main
 │
 ├── start goroutine
 │
 └── exit
      ↓
 process exits
```

Запуск работы и ожидание завершения — разные задачи.

Это первый фундаментальный lifecycle rule:

> **Создание concurrent work не означает ожидание её завершения.**

---

# 6. WaitGroup

`sync.WaitGroup` используется, когда необходимо дождаться завершения группы goroutines.

```go
var wg sync.WaitGroup

wg.Add(1)

go func() {
	defer wg.Done()

	process(job)
}()

wg.Wait()
```

Семантика:

```text
Add
 ↓
register work

Done
 ↓
work completed

Wait
 ↓
wait until counter == 0
```

Практический idiom:

```go
wg.Add(1)

go func() {
	defer wg.Done()

	// work
}()
```

`defer wg.Done()` снижает риск забыть обозначить завершение при раннем `return`.

---

# 7. Concurrency vs Parallelism

### Concurrency

Структура выполнения, при которой несколько независимых задач могут прогрессировать независимо.

### Parallelism

Фактическое одновременное выполнение нескольких частей работы на разных execution resources.

Условная модель:

```text
Concurrency

G1 ────────┐
G2 ────────┼── scheduler
G3 ────────┘
```

против:

```text
Parallelism

CPU 1 → G1
CPU 2 → G2
CPU 3 → G3
```

Главная мысль:

> **Concurrency описывает организацию работы. Parallelism описывает фактическое одновременное выполнение.**

---

# 8. Goroutine per Job

Наивный вариант:

```go
func processJobs(jobs []Job) {
	var wg sync.WaitGroup

	for _, job := range jobs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			process(job)
		}()
	}

	wg.Wait()
}
```

При:

```text
100 jobs
```

может быть:

```text
100 goroutines
```

При:

```text
1 000 000 jobs
```

потенциально:

```text
1 000 000 goroutines
```

Проблема не в том, что goroutines существуют.

Проблема:

> **Входной объём работы напрямую определяет количество concurrent resources.**

Нужна bounded concurrency.

---

# 9. Channel

Channel используется для передачи значений и координации между goroutines.

```go
jobs := make(chan Job)
```

Producer:

```go
jobs <- job
```

Consumer:

```go
job := <-jobs
```

Модель:

```text
producer
   │
   │ jobs <- job
   ▼
channel
   │
   │ job := <-jobs
   ▼
consumer
```

Channel — не просто контейнер.

Он также задаёт synchronization semantics.

---

# 10. Unbuffered Channel

```go
jobs := make(chan Job)
```

При send:

```go
jobs <- job
```

sender ждёт receiver.

При receive:

```go
job := <-jobs
```

receiver ждёт sender.

Упрощённо:

```text
sender                  receiver

jobs <- job ──────────→ job := <-jobs
         │                  │
         └ synchronization ┘
```

Вопрос:

> Что произойдёт, если никто не читает из unbuffered channel?

<details><summary>Ответ</summary>Отправка будет блокироваться до появления receiver.</details>

---

# 11. Buffered Channel

```go
jobs := make(chan Job, 100)
```

Теперь существует capacity:

```text
┌───────────────────────────────┐
│ job │ job │ job │ ...         │
└───────────────────────────────┘
             capacity = 100
```

Producer может временно опережать consumer до заполнения buffer.

Но:

> **Buffer не создаёт бесконечную очередь и не устраняет backpressure.**

После заполнения capacity следующий send снова должен ждать либо система применит другую policy.

---

# 12. Channel Direction

Worker, которому нужно только читать:

```go
jobs <-chan Job
```

Producer, которому нужно только отправлять:

```go
chan<- Job
```

Это позволяет выразить направление использования прямо в типе и помогает compiler обнаруживать неправильные операции.

```text
chan Job
 ├── <-chan Job
 │      receive
 │
 └── chan<- Job
        send
```

---

# 13. Closing Channel

Producer может закрыть канал после того, как новых Job больше не будет:

```go
close(jobs)
```

Worker:

```go
for job := range jobs {
	process(job)
}
```

После закрытия channel:

* уже переданные значения могут быть прочитаны;
* новые значения отправлять нельзя;
* `range` завершится после исчерпания значений.

`close` означает:

> **Новых значений больше не будет.**

Он не означает:

> **Немедленно остановить receiver.**

---

# 14. Кто закрывает channel?

Основное правило:

> **Channel закрывает сторона, которая владеет отправкой и знает, что новых значений больше не будет.**

В worker pool:

```text
producer
   │
   │ owns sending
   ▼
jobs channel
   │
   ├── worker
   ├── worker
   └── worker
```

Worker не должен закрывать канал, которым он не владеет.

---

# 15. `select`

Worker может одновременно ожидать Job и cancellation:

```go
select {
case job := <-jobs:
	process(job)

case <-ctx.Done():
	return
}
```

Модель:

```text
             ┌── job available ──→ process
select ──────┤
             └── context done ───→ stop
```

`select` позволяет ожидать несколько channel operations или событий.

---

# 16. Плакат 03B — Go Runtime

![03B — Concurrency Runtime](../posters/03B_CONCURRENCY_RUNTIME.svg)

Плакат показывает, что происходит внутри Go runtime после запуска goroutine.

Основная модель:

```text
go f()
   ↓
G
   ↓
RUNQ
   ↓
scheduler
   ↓
G + P + M
   ↓
execute
   ↓
wait / preempt
   ↓
wake
   ↓
RUNQ
```

---

# 17. G / P / M

### G

Goroutine — единица работы.

### P

Processor runtime scheduler — ресурс, необходимый для выполнения Go-кода.

### M

Machine — OS thread, на котором фактически выполняется Go code.

Упрощённая связь:

```text
G = work
P = runtime execution resource
M = OS thread
```

Важно:

> **G, P и M — не три последовательных шага. Это связанные сущности runtime scheduler.**

---

# 18. RUNQ

Runnable goroutines ожидают выполнения в scheduler queues.

Упрощённо:

```text
P0 → G G G
P1 → G G
P2 → G G G G
```

Это не значит, что все G выполняются одновременно.

Это означает, что scheduler располагает runnable work и выбирает, что выполнять дальше.

---

# 19. Scheduler

Когда существует runnable G, runtime scheduler выбирает execution path.

Упрощённая модель:

```text
G
↓
run queue
↓
scheduler
↓
P + M
↓
execution
```

Scheduler позволяет многим goroutines прогрессировать даже при ограниченном количестве CPU resources.

---

# 20. Execute / Preemption / Wait

Во время выполнения goroutine может:

```text
compute
wait for syscall
wait for network I/O
be preempted
```

При network I/O Go runtime использует netpoller, чтобы эффективно управлять ожиданием сетевых событий.

Упрощённая модель:

```text
G
 ↓
network wait
 ↓
netpoll
 ↓
event ready
 ↓
G runnable again
 ↓
scheduler
```

То есть ожидание network I/O не должно означать бессмысленное удержание CPU только ради ожидания.

---

# 21. Важный runtime вывод

После:

```go
go process(job)
```

нельзя считать, что:

```text
process(job)
```

немедленно выполняется на CPU.

Более точная mental model:

```text
create
→ runnable
→ scheduler
→ execute
```

А во время жизни goroutine выполнение может переходить между состояниями runnable, running и waiting.

---

# 22. Плакат 03C — Concurrency Engineering

![03C — Concurrency Engineering](../posters/03C_CONCURRENCY_ENGINEERING.svg)

Теперь вопрос меняется.

Мы умеем создавать concurrent work.

Но:

> **Как сделать concurrent system управляемой?**

Ключевые элементы:

```text
LOAD
BOUND
BACKPRESSURE
WORKER POOL
CANCELLATION
SHARED STATE
```

---

# 23. Load

Нельзя оценивать concurrency только количеством запросов.

Нужно учитывать:

```text
rate
burst
cost per operation
```

Например:

```text
1000 jobs/sec
```

и:

```text
1 ms / job
```

радикально отличаются от:

```text
1000 jobs/sec
```

и:

```text
500 ms / job
```

Поэтому вопрос:

> **Сколько goroutines создать?**

не является первым инженерным вопросом.

Сначала:

> **Какова нагрузка и стоимость одной операции?**

---

# 24. Bounded Concurrency

Пусть:

```text
1 000 000 jobs
```

и:

```text
32 workers
```

Вместо:

```text
1 000 000 concurrent executions
```

получаем:

```text
maximum concurrent workers = 32
```

Это ограничивает потребление:

```text
CPU
memory
connections
downstream capacity
```

Необходимы явные resource bounds.

---

# 25. Backpressure

Если:

```text
producer = 1000 jobs/sec
consumer = 100 jobs/sec
```

то backlog растёт.

Без ограничения:

```text
queue
 ↓
memory
 ↓
latency
 ↓
resource exhaustion
```

Backpressure заставляет producer и систему учитывать реальную capacity consumer.

Возможные policy:

```text
block
throttle
reject
drop
buffer
shed load
```

Конкретная policy зависит от требований системы.

Главная идея:

> **Producer не должен иметь бесконечный кредит ресурсов.**

---

# 26. Worker Pool

Базовая архитектура:

```text
                 jobs
                  │
                  ▼
                QUEUE
                  │
        ┌─────────┼─────────┐
        ▼         ▼         ▼
       W1        W2        W3
        │         │         │
        └─────────┼─────────┘
                  ▼
               RESULT
```

Основная идея:

```text
input volume
     ≠
worker count
```

Количество workers ограничено конфигурацией.

---

# 27. Базовый worker

```go
func worker(
	id int,
	jobs <-chan Job,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
		process(job)
	}
}
```

Запуск:

```go
const workerCount = 4

var wg sync.WaitGroup

wg.Add(workerCount)

for i := 0; i < workerCount; i++ {
	go worker(i, jobs, &wg)
}
```

Producer:

```go
for _, job := range jobsToProcess {
	jobs <- job
}

close(jobs)
wg.Wait()
```

Lifecycle:

```text
start
 ↓
receive
 ↓
process
 ↓
repeat
 ↓
channel close
 ↓
return
 ↓
WaitGroup Done
```

---

# 28. Context и Cancellation

Context используется для передачи:

```text
cancellation
deadline
request-scoped metadata
```

Типичный backend flow:

```text
HTTP request
     │
     ▼
context.Context
     │
     ├── service
     │     └── repository
     │
     └── downstream call
```

Cancellation — сигнал, а не принудительное убийство goroutine.

Код должен самостоятельно реагировать на:

```go
<-ctx.Done()
```

---

# 29. Timeout

Пример:

```go
ctx, cancel := context.WithTimeout(
	context.Background(),
	time.Second,
)
defer cancel()
```

Через заданный срок context будет отменён.

Worker:

```go
select {
case job := <-jobs:
	process(job)

case <-ctx.Done():
	return
}
```

Главная идея:

> **Длительная операция должна иметь допустимый срок ожидания или другой определённый путь завершения.**

---

# 30. Shared Mutable State

Проблемный код:

```go
type Stats struct {
	processed int
}
```

и:

```go
for i := 0; i < 1000; i++ {
	go func() {
		stats.processed++
	}()
}
```

Наивное ожидание:

```text
processed == 1000
```

не гарантируется.

Причина:

```text
read
 ↓
add
 ↓
write
```

`processed++` — не одна неделимая операция.

---

# 31. Data Race

Concurrent access:

```text
G1: read 10
G2: read 10
G1: write 11
G2: write 11
```

Возможный результат:

```text
expected → 12
actual   → 11
```

Это пример потерянного обновления.

В более общем смысле проблема — data race и отсутствие необходимой synchronization.

---

# 32. Race Detector

Проверка:

```bash
go test -race ./...
```

Race detector помогает обнаруживать конфликтующий concurrent access во время выполнения тестов.

Это важная часть evidence для concurrent code.

Но:

> Race detector не доказывает отсутствие всех concurrency bugs.

Он является инструментом проверки определённого класса ошибок.

---

# 33. Mutex

Для shared state:

```go
type Stats struct {
	mu        sync.Mutex
	processed int
}
```

Использование:

```go
stats.mu.Lock()
defer stats.mu.Unlock()

stats.processed++
```

Mutex защищает критическую секцию и связанное с ней состояние.

Главный вопрос при проектировании:

> **Какой invariant защищает этот mutex?**

Не просто:

> «На эту переменную нужно поставить mutex.»

---

# 34. Atomic

Для независимого простого counter:

```go
type Stats struct {
	processed atomic.Int64
}
```

Increment:

```go
stats.processed.Add(1)
```

Read:

```go
stats.processed.Load()
```

Сравнение:

```text
Mutex
→ critical section / related state

Atomic
→ atomic operation on one value
```

Atomic — не «быстрый mutex».

Это другой primitive и другая модель использования.

---

# 35. Mutex vs Atomic

Для:

```text
processed++
```

подходит atomic counter.

Но если есть:

```go
type Stats struct {
	processed int
	failed    int
}
```

и invariant:

```text
processed + failed == total
```

нужно менять связанное состояние согласованно.

В таком случае mutex или другой composite-state synchronization обычно естественнее.

---

# 36. Concurrent Memory Repository

Исходный repository:

```go
type MemoryJobRepository struct {
	jobs map[string]Job
}
```

При concurrent access появляется новая проблема.

Обычная map не является механизмом синхронизации concurrent writes.

Возможная модель:

```go
type MemoryJobRepository struct {
	mu   sync.RWMutex
	jobs map[string]Job
}
```

Write:

```go
r.mu.Lock()
defer r.mu.Unlock()

r.jobs[job.ID] = job
```

Read:

```go
r.mu.RLock()
defer r.mu.RUnlock()

job, ok := r.jobs[id]
```

`RWMutex` стоит использовать только если выбранная модель нагрузки оправдывает shared-reader semantics.

---

# 37. Goroutine Leak

Проблемный код:

```go
func startWorker(jobs <-chan Job) {
	go func() {
		for {
			job := <-jobs
			process(job)
		}
	}()
}
```

Если producer больше никогда не отправляет Job и не существует другого stop condition, worker может остаться жить бесконечно.

Это goroutine leak.

---

# 38. Lifecycle-safe Worker

Возможный вариант:

```go
func worker(
	ctx context.Context,
	jobs <-chan Job,
) {
	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			process(job)
		}
	}
}
```

Теперь существует два явных пути завершения:

```text
context cancellation
        или
channel close
```

---

# 39. Graceful Shutdown

Production-процесс должен уметь завершаться контролируемо.

Модель:

```text
SIGTERM
   ↓
STOP ACCEPTING
   ↓
CANCEL / DRAIN
   ↓
WORKERS EXIT
   ↓
WAIT
   ↓
CLOSE RESOURCES
   ↓
EXIT
```

Принцип:

> **Graceful shutdown — это политика завершения работы, а не убийство goroutines.**

---

# 40. Практическое задание — WorkerPool

Реализовать:

```go
type WorkerPool struct {
	workers int
}
```

с методом:

```go
func (p *WorkerPool) Run(
	ctx context.Context,
	jobs []Job,
	process func(context.Context, Job) error,
) error
```

Требования:

```text
workers configurable
all allowed jobs processed
maximum concurrency bounded
context cancellation supported
Run waits for workers
processing errors observable
no worker goroutine remains after Run
```

Архитектурные вопросы перед реализацией:

```text
кто владеет queue?
кто создаёт workers?
кто закрывает queue?
кто останавливает workers?
кто ждёт завершения?
```

---

# 41. Support Ladder

Если решение не получается, проверяй его последовательно.

### Level 1

Где находится очередь задач?

### Level 2

Кто отправляет Job?

### Level 3

Кто читает Job?

### Level 4

Как ограничено количество readers/workers?

### Level 5

Как дождаться завершения workers?

### Level 6

Как worker узнает о shutdown?

### Level 7

Что произойдёт, если cancellation приходит во время ожидания Job?

Ожидаемые инструменты:

```text
channel
WaitGroup
select
context
```

---

# 42. Failure Drills

После реализации WorkerPool проверь систему не только в happy path.

## Сценарий 1

Producer быстрее workers.

Что произойдёт?

<details><summary>Ответ</summary>Начнёт расти backlog. Нужна bounded queue или другая backpressure policy.</details>

## Сценарий 2

Worker работает слишком долго.

Что нужно?

<details><summary>Ответ</summary>Timeout/deadline/cancellation policy.</details>

## Сценарий 3

Worker перестал быть нужен.

Что остановит его?

<details><summary>Ответ</summary>Явный stop condition: context cancellation, channel close или другой определённый lifecycle mechanism.</details>

## Сценарий 4

Два workers меняют одну map.

Что проверить?

<details><summary>Ответ</summary>Есть ли shared mutable state и необходимая synchronization. Запустить `go test -race ./...`.</details>

## Сценарий 5

Приходит SIGTERM.

Что должно произойти?

<details><summary>Ответ</summary>Остановить приём новой работы, инициировать cancellation/drain, дождаться завершения workers и закрыть ресурсы.</details>

---

# 43. Финальная mental model

К концу урока должна быть понятна вся цепочка:

```text
                    CONCURRENCY
                         │
       ┌─────────────────┼─────────────────┐
       │                 │                 │
   execution         communication      state
       │                 │                 │
   goroutine           channel          mutex
   scheduler           select           atomic
       │                 │                 │
       └─────────────────┼─────────────────┘
                         │
                  worker pool
                         │
                 bounded resources
                         │
                  cancellation
                         │
                  graceful shutdown
```

Ключевая формула:

> **Goroutine выполняет работу. Channel передаёт работу. Select выбирает между событиями. WaitGroup ждёт завершения. Context сообщает об отмене. Mutex защищает shared state. Atomic выполняет атомарную операцию.**

И более важный инженерный вывод:

> **Concurrent component должен иметь явные ownership, lifecycle, synchronization и resource bounds.**

---

# 44. Checkpoint

Ответь без IDE.

### Что такое goroutine?

<details><summary>Ответ</summary>Единица concurrent execution, управляемая Go runtime.</details>

### Что такое concurrency?

<details><summary>Ответ</summary>Организация независимого выполнения и координации нескольких задач.</details>

### Что такое parallelism?

<details><summary>Ответ</summary>Фактическое одновременное выполнение нескольких частей работы на разных execution resources.</details>

### Зачем `WaitGroup`?

<details><summary>Ответ</summary>Дождаться завершения группы goroutines.</details>

### Зачем channel?

<details><summary>Ответ</summary>Для передачи значений и coordination между goroutines.</details>

### Кто закрывает channel?

<details><summary>Ответ</summary>Обычно сторона, владеющая отправкой и знающая, что новых значений больше не будет.</details>

### Что делает `select`?

<details><summary>Ответ</summary>Позволяет ожидать несколько channel operations или событий.</details>

### Что такое data race?

<details><summary>Ответ</summary>Конкурентный конфликтующий доступ к общей памяти без необходимой synchronization, при котором хотя бы одна операция является записью.</details>

### Зачем `go test -race`?

<details><summary>Ответ</summary>Для обнаружения определённых data races во время выполнения тестов.</details>

### Зачем mutex?

<details><summary>Ответ</summary>Для защиты shared mutable state и связанных invariant.</details>

### Зачем atomic?

<details><summary>Ответ</summary>Для атомарных операций над простым shared state.</details>

### Что такое backpressure?

<details><summary>Ответ</summary>Механизм ограничения producer, когда consumer не успевает обрабатывать поступающую работу.</details>

### Что такое goroutine leak?

<details><summary>Ответ</summary>Goroutine, которая продолжает жить без необходимости или без достижимого stop condition.</details>

### Зачем context?

<details><summary>Ответ</summary>Для передачи cancellation и deadline через call chain.</details>

### Что такое graceful shutdown?

<details><summary>Ответ</summary>Контролируемое прекращение приёма работы, cancellation/drain, ожидание завершения owned operations и закрытие ресурсов.</details>

---

# 45. Домашнее задание — JobFlow v2

Продолжить тот же JobFlow.

## 1. Worker Pool

Добавить:

```go
pool := NewWorkerPool(4)
```

и:

```go
Run(ctx, jobs, process)
```

Количество workers должно быть configurable.

## 2. Concurrent-safe Repository

Сделать `MemoryJobRepository` безопасным для concurrent access.

Проверить:

```text
concurrent Add
concurrent Get
concurrent Delete
concurrent List
```

## 3. Statistics

Добавить:

```text
processed
failed
```

Для независимых counters использовать atomic.

## 4. Cancellation

Проверить:

```text
cancel before processing
cancel while waiting for jobs
cancel during processing
```

## 5. Shutdown

Доказать, что после завершения `Run()` worker goroutines, принадлежащие pool, не продолжают работу.

## 6. Tests

Добавить concurrent tests.

Обязательно:

```bash
go test ./...
go test -race ./...
go vet ./...
```

---

# 46. PASS / FAIL

## PASS

JobFlow:

```text
worker count bounded
race-free under tested scenarios
cancellation works
shutdown completes
repository is concurrency-safe
tests pass
```

Ученик способен объяснить:

```text
owner
lifecycle
queue
worker
backpressure
shared state
synchronization
cancellation
shutdown
```

## FAIL

Решение не считается готовым, если:

* количество workers зависит от количества Job;
* worker не имеет stop condition;
* shared state изменяется без определённой synchronization;
* `go test -race ./...` обнаруживает race;
* `Run()` завершается, пока owned workers продолжают жить;
* ученик не может объяснить ownership и lifecycle concurrent component.

---

# 47. Следующая тема

После этого урока JobFlow уже умеет выполнять работу конкурентно.

Следующая задача:

```text
HTTP
   ↓
Application
   ↓
JobService
   ↓
Worker / Repository
```

Новый вопрос:

> **Как отделить HTTP transport от business policy и infrastructure, не превратив приложение в набор связанных деталей?**
