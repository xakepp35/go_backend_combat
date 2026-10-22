
# Go Backend Combat — 03. CONCURRENCY

# Homework: Concurrency Bug Hunt

## Цель

В этом домашнем задании **все программы сломаны намеренно**.

Каждый пример:

- компилируется или почти компилируется;
- выглядит правдоподобно;
- содержит один или несколько concurrency defects;
- может зависать, паниковать, терять данные, оставлять goroutines, нарушать invariant или некорректно завершаться.

Ваша задача — не переписать всё с нуля, а пройти инженерный цикл:

```text
OBSERVE
   ↓
FIND THE BROKEN GUARANTEE
   ↓
IDENTIFY OWNERSHIP
   ↓
FIX SYNCHRONIZATION / LIFECYCLE
   ↓
VERIFY
```

Главный вопрос каждой задачи:

> **Какую гарантию исходный код не способен обеспечить?**

---

# Общие правила

Для каждой задачи:

1. Запустите исходный код.
2. Если применимо — запустите его с race detector.
3. Определите дефект.
4. Исправьте код.
5. Сохраните исходный смысл программы.
6. Докажите, что исправление работает.

В начале каждого решения оставьте короткий комментарий:

```go
// BUG:
// ...
//
// FIX:
// ...
//
// GUARANTEE:
// ...
```

Не требуется длинное эссе.

Достаточно нескольких предложений.

---

## Запрещённый способ исправления

Нельзя использовать:

```go
time.Sleep(...)
```

как средство синхронизации.

Например:

```go
go work()

time.Sleep(time.Second) // "ну вроде должна успеть"
```

не является корректным ожиданием завершения работы.

`time.Sleep` разрешён только внутри примеров, где он имитирует длительную операцию.

---

## Verification

После выполнения домашней работы:

```bash
go fmt ./...
go test ./...
go test -race ./...
go vet ./...
```

Для concurrency-кода:

```text
"у меня один раз сработало"
```

не является доказательством корректности.

---

# CORE

Следующие задачи являются основной частью домашнего задания.

---

# Task 01 — Lost Counter

## Тема

```text
shared mutable state
data race
Mutex / atomic
WaitGroup
```

Программа должна запустить 1000 goroutines.

Каждая goroutine увеличивает общий счётчик 1000 раз.

Ожидаемый результат:

```text
1_000_000
```

Исходный код:

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	const goroutines = 1000
	const increments = 1000

	var wg sync.WaitGroup
	var counter int64

	for i := 0; i < goroutines; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < increments; j++ {
				counter++
			}
		}()
	}

	wg.Wait()

	fmt.Println("counter:", counter)
}
```

## Требуемое поведение

После исправления программа всегда должна печатать:

```text
counter: 1000000
```

и проходить:

```bash
go run -race .
```

без сообщения о data race.

## Вопрос после решения

Почему:

```go
counter++
```

не является одной неделимой операцией?

Когда для такой задачи разумнее использовать:

```text
sync.Mutex
```

а когда:

```text
sync/atomic
```

?

---

# Task 02 — WaitGroup Never Ends

## Тема

```text
goroutine lifecycle
WaitGroup
ownership
copying synchronization primitives
```

Необходимо параллельно выполнить несколько условных HTTP-запросов.

Сетевой вызов здесь имитируется через `time.Sleep`.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Result struct {
	URL    string
	Status int
}

func fetch(
	url string,
	wg sync.WaitGroup,
	results chan<- Result,
) {
	defer wg.Done()

	time.Sleep(100 * time.Millisecond)

	results <- Result{
		URL:    url,
		Status: 200,
	}
}

func main() {
	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://golang.org",
	}

	results := make(chan Result, len(urls))

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go fetch(url, wg, results)
	}

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Printf(
			"%s -> %d\n",
			result.URL,
			result.Status,
		)
	}

	fmt.Println("done")
}
```

## Требуемое поведение

Программа должна:

1. запустить три операции concurrently;
2. дождаться всех трёх;
3. получить три результата;
4. завершиться самостоятельно;
5. вывести `done`.

Количество операций может изменяться.

То есть решение не должно зависеть от магического числа `3`.

## Дополнительная проверка

Запустите:

```bash
go vet ./...
```

Объясните, почему synchronization primitives обычно не следует копировать после начала использования.

---

# Task 03 — Worker Pool Deadlock

## Тема

```text
worker pool
channels
WaitGroup
channel ownership
close
bounded concurrency
```

Есть 100 тяжёлых задач.

Одновременно разрешено выполнять не более пяти.

```go
package main

import (
	"fmt"
	"time"
)

type Task struct {
	ID int
}

func worker(
	id int,
	tasks <-chan Task,
	results chan<- string,
) {
	for task := range tasks {
		fmt.Printf(
			"worker %d started task %d\n",
			id,
			task.ID,
		)

		time.Sleep(50 * time.Millisecond)

		results <- fmt.Sprintf(
			"worker %d completed task %d",
			id,
			task.ID,
		)
	}
}

func main() {
	const numTasks = 100
	const numWorkers = 5

	tasks := make(chan Task)
	results := make(chan string)

	for i := 1; i <= numWorkers; i++ {
		go worker(i, tasks, results)
	}

	for i := 1; i <= numTasks; i++ {
		tasks <- Task{ID: i}
	}

	for result := range results {
		fmt.Println(result)
	}

	fmt.Println("all tasks completed")
}
```

## Требуемое поведение

После исправления:

```text
tasks processed = 100
maximum workers = 5
program terminates
```

Каждая задача должна быть обработана ровно один раз.

После завершения работы workers должны завершиться.

`results` также должен иметь корректный lifecycle.

---

## Вопросы

Для каждого канала определите owner.

### `tasks`

Кто:

```text
создаёт?
пишет?
читает?
закрывает?
```

### `results`

Кто:

```text
создаёт?
пишет?
читает?
понимает, что writers закончились?
закрывает?
```

Не закрывайте канал просто:

> «потому что где-то надо сделать `close`».

Сначала определите ownership.

---

# Task 04 — Fake Backpressure

## Тема

```text
bounded queue
backpressure
unbounded goroutines
resource bounds
```

Разработчик хотел ограничить очередь восьмью элементами.

Workers всего четыре.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

func worker(
	jobs <-chan Job,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("processed:", job.ID)
	}
}

func main() {
	const workers = 4
	const queueSize = 8
	const jobsCount = 100_000

	jobs := make(chan Job, queueSize)

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go worker(jobs, &wg)
	}

	for i := 0; i < jobsCount; i++ {
		job := Job{ID: i}

		go func() {
			jobs <- job
		}()
	}

	close(jobs)

	wg.Wait()
}
```

## Идея автора

Автор считает:

```text
queue size = 8
workers    = 4
```

следовательно:

```text
система bounded
```

## Требуемое поведение

Исправьте программу так, чтобы количество входных `Job` **не превращалось в количество созданных goroutines**.

При:

```text
100_000 jobs
```

система по-прежнему должна иметь явно ограниченный объём concurrent work.

Если workers не успевают обрабатывать jobs, producer должен столкнуться с реальной capacity системы.

---

## Вопрос

Объясните разницу:

```text
BUFFER
```

и:

```text
BACKPRESSURE
```

Почему:

```go
go func() {
	jobs <- job
}()
```

может фактически уничтожить backpressure, даже если сам channel bounded?

---

# Task 05 — Pipeline That Never Finishes

## Тема

```text
pipeline
channel lifecycle
close propagation
producer ownership
```

Необходимо построить конвейер:

```text
GENERATE
   ↓
EVEN
   ↓
SQUARE
   ↓
PRINT
```

Первая стадия генерирует числа `1..10`.

Вторая пропускает только чётные.

Третья возводит их в квадрат.

```go
package main

import "fmt"

func generate(n int) <-chan int {
	out := make(chan int)

	go func() {
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()

	return out
}

func even(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for value := range in {
			if value%2 == 0 {
				out <- value
			}
		}
	}()

	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for value := range in {
			out <- value * value
		}
	}()

	return out
}

func main() {
	numbers := generate(10)
	evens := even(numbers)
	squares := square(evens)

	for value := range squares {
		fmt.Println(value)
	}

	fmt.Println("pipeline completed")
}
```

## Ожидаемый результат

Программа должна вывести:

```text
4
16
36
64
100
pipeline completed
```

и самостоятельно завершиться.

---

## Инвариант pipeline stage

Каждая стадия должна иметь понятную модель:

```text
receive until input is finished
        ↓
process
        ↓
send output
        ↓
finish own output
        ↓
exit
```

Определите, какая goroutine является owner каждого output channel.

---

# Task 06 — Stop All Workers

## Тема

```text
fan-out
broadcast signal
channel close
lifecycle
```

Запущено пять долгоживущих workers.

Необходимо одной операцией остановить их всех.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(
	id int,
	stop <-chan struct{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-stop:
			fmt.Println("worker stopped:", id)
			return

		case <-time.After(100 * time.Millisecond):
			fmt.Println("worker working:", id)
		}
	}
}

func main() {
	const workers = 5

	stop := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go worker(i, stop, &wg)
	}

	time.Sleep(300 * time.Millisecond)

	stop <- struct{}{}

	wg.Wait()

	fmt.Println("all workers stopped")
}
```

## Требуемое поведение

Одна операция остановки должна стать наблюдаемой **всеми workers**.

После сигнала:

```text
W1 exits
W2 exits
W3 exits
W4 exits
W5 exits
```

`main` должен дождаться их завершения.

---

## Вопрос

Чем семантически отличаются:

```go
ch <- struct{}{}
```

и:

```go
close(ch)
```

когда channel используется как сигнал?

Кто имеет право закрывать этот channel?

---

# Task 07 — Broken Fan-In

## Тема

```text
fan-in
multiple producers
single consumer
channel ownership
WaitGroup
```

Есть два независимых источника данных.

Их необходимо объединить в один channel.

```go
package main

import (
	"fmt"
)

func source(start, end int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for i := start; i <= end; i++ {
			out <- i
		}
	}()

	return out
}

func merge(
	a <-chan int,
	b <-chan int,
) <-chan int {
	out := make(chan int)

	copyValues := func(in <-chan int) {
		defer close(out)

		for value := range in {
			out <- value
		}
	}

	go copyValues(a)
	go copyValues(b)

	return out
}

func main() {
	a := source(1, 5)
	b := source(100, 105)

	for value := range merge(a, b) {
		fmt.Println(value)
	}
}
```

## Требуемое поведение

Результирующий channel должен получить все значения:

```text
1..5
100..105
```

Порядок между источниками не определён.

Например:

```text
1
100
2
3
101
...
```

является корректным.

После завершения обоих источников результирующий channel должен закрыться ровно один раз.

---

## Главное требование

Ни один producer не должен закрыть общий output channel, пока другие producers всё ещё способны отправлять в него данные.

Ответьте:

> Кто знает, что **все** writers завершились?

Именно эта сторона должна определять lifecycle результирующего channel.

---

# Task 08 — Timeout With a Hidden Leak

## Тема

```text
timeout
select
abandoned receiver
goroutine leak
```

Функция выполняет условный запрос в базу данных.

Если ответ не пришёл за `200 ms`, вызывающая сторона должна получить timeout.

```go
package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrTimeout = errors.New("query timeout")

func queryDB() string {
	time.Sleep(500 * time.Millisecond)
	return "result"
}

func load() (string, error) {
	result := make(chan string)

	go func() {
		value := queryDB()
		result <- value
	}()

	select {
	case value := <-result:
		return value, nil

	case <-time.After(200 * time.Millisecond):
		return "", ErrTimeout
	}
}

func main() {
	value, err := load()

	fmt.Println("value:", value)
	fmt.Println("error:", err)
}
```

## На первый взгляд

Функция действительно возвращает:

```text
query timeout
```

примерно через `200 ms`.

Но проанализируйте lifecycle запущенной goroutine **после того, как `load()` уже вернулась**.

---

## Требуемое поведение

После timeout concurrent operation должна иметь достижимый путь завершения.

Нельзя оставлять goroutine навсегда заблокированной только потому, что вызывающая сторона больше не ждёт её результат.

Исправление должно сохранить:

```text
timeout ≈ 200 ms
```

при медленной операции.

---

## Вопрос

Что произошло с ownership результата после того, как caller решил:

```text
"результат мне больше не нужен"
```

?

Как producer узнаёт об этом?

---

# Task 09 — Cancellation Doesn't Cancel

## Тема

```text
context.Context
select
cancellation
blocking receive
blocking send
goroutine lifecycle
```

Worker должен завершаться при отмене `ctx`.

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

func worker(
	ctx context.Context,
	jobs <-chan Job,
	results chan<- int,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
		time.Sleep(200 * time.Millisecond)

		results <- job.ID
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	jobs := make(chan Job)
	results := make(chan int)

	var wg sync.WaitGroup
	wg.Add(1)

	go worker(
		ctx,
		jobs,
		results,
		&wg,
	)

	go func() {
		jobs <- Job{ID: 1}
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	wg.Wait()

	fmt.Println("worker stopped")
}
```

## Требуемый контракт

Worker должен уметь завершиться при cancellation:

```text
while waiting for a job
while processing / waiting
while attempting to publish result
```

Если `ctx` отменён, worker не должен зависеть от того:

```text
закроет ли кто-то jobs когда-нибудь
появится ли receiver результатов когда-нибудь
```

`main` должен дождаться worker.

---

## Важная мысль

`context.CancelFunc`:

```text
не убивает goroutine
```

Он только делает наблюдаемым сигнал:

```go
ctx.Done()
```

Goroutine должна сама иметь execution path, который этот сигнал обрабатывает.

---

# Task 10 — Concurrent Repository

## Тема

```text
map
shared mutable state
Mutex
RWMutex
race detector
```

Есть простейший in-memory repository.

```go
package main

import (
	"fmt"
	"sync"
)

type Repository struct {
	values map[int]int
}

func NewRepository() *Repository {
	return &Repository{
		values: make(map[int]int),
	}
}

func (r *Repository) Set(
	key int,
	value int,
) {
	r.values[key] = value
}

func (r *Repository) Get(
	key int,
) (int, bool) {
	value, ok := r.values[key]
	return value, ok
}

func main() {
	repo := NewRepository()

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(2)

		go func(id int) {
			defer wg.Done()

			repo.Set(id, id*10)
		}(i)

		go func(id int) {
			defer wg.Done()

			repo.Get(id)
		}(i)
	}

	wg.Wait()

	fmt.Println("completed")
}
```

## Требуемое поведение

Repository должен поддерживать concurrent:

```text
Get
Set
```

без:

```text
data race
concurrent map read and map write
concurrent map writes
```

После исправления:

```bash
go run -race .
```

должен проходить без race report.

---

## Вопрос

Если выбран `RWMutex`, объясните:

```text
какой invariant он защищает
```

и почему:

```text
RLock
```

может использоваться для `Get`, но не для `Set`.

Сам факт существования `RWMutex` не означает, что он всегда лучше обычного `Mutex`.

---

# Task 11 — Shutdown Means Finished

## Тема

```text
ownership
worker lifecycle
graceful shutdown
channel close
WaitGroup
concurrent shutdown
```

Есть небольшой worker service.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

type Service struct {
	jobs chan Job
	wg   sync.WaitGroup
}

func NewService(workers int) *Service {
	s := &Service{
		jobs: make(chan Job, 16),
	}

	s.wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer s.wg.Done()

			for job := range s.jobs {
				time.Sleep(100 * time.Millisecond)

				fmt.Printf(
					"worker %d processed %d\n",
					id,
					job.ID,
				)
			}
		}(i)
	}

	return s
}

func (s *Service) Submit(job Job) {
	s.jobs <- job
}

func (s *Service) Shutdown() {
	close(s.jobs)

	fmt.Println("service stopped")
}

func main() {
	service := NewService(4)

	for i := 0; i < 20; i++ {
		service.Submit(Job{ID: i})
	}

	service.Shutdown()
}
```

## Проблема

Метод называется:

```go
Shutdown()
```

Но что именно гарантировано в момент его возврата?

Действительно ли service уже завершён?

---

## Требуемый контракт

После исправления:

```go
service.Shutdown()
```

должно означать:

```text
новая работа больше не принимается
        ↓
уже принятая работа завершается
        ↓
workers завершаются
        ↓
Shutdown возвращает управление
```

После возврата из `Shutdown` никакой worker, принадлежащий этому service, не должен продолжать работу.

---

## Дополнительный сценарий

Исправьте API так, чтобы попытка:

```go
service.Submit(job)
```

после начала shutdown:

```text
не приводила к panic: send on closed channel
```

Она должна завершаться контролируемым способом.

Например:

```go
var ErrStopped = errors.New("service stopped")
```

и:

```go
Submit(job) error
```

---

## Combat requirement

Повторный вызов:

```go
service.Shutdown()
```

также не должен приводить к:

```text
panic: close of closed channel
```

Определите явно, кому принадлежит lifecycle service.

---

# CORE PASS / FAIL

Основная часть считается выполненной, если ученик способен починить задачи 01–11 и объяснить:

```text
goroutine ownership
WaitGroup lifecycle
channel ownership
channel close
fan-out
fan-in
pipeline
bounded concurrency
backpressure
context cancellation
timeout
shared mutable state
mutex / atomic
goroutine leak
graceful shutdown
```

И проходят:

```bash
go test ./...
go test -race ./...
go vet ./...
```

---

# COMBAT+

Следующие задачи выходят немного за минимальное содержание урока, но относятся к паттернам, которые полезно узнавать с первого взгляда.

---

# Task 12 ★ — Initialize Exactly Once

## Тема

```text
lazy initialization
sync.Once
publication
data race
```

Создание клиента — дорогая операция.

Все goroutines должны получить один и тот же экземпляр.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Client struct {
	ID int
}

var client *Client

func getClient() *Client {
	if client == nil {
		time.Sleep(100 * time.Millisecond)

		client = &Client{
			ID: 42,
		}
	}

	return client
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			c := getClient()
			fmt.Printf("%p\n", c)
		}()
	}

	wg.Wait()
}
```

## Требования

Heavy initialization должна происходить:

```text
ровно один раз
```

Все callers должны получать один опубликованный экземпляр.

Программа должна быть race-free.

---

# Task 13 ★ — Deadlock Transfer

## Тема

```text
mutex
lock ordering
deadlock
deterministic acquisition order
```

Есть банковские счета:

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Account struct {
	ID      int
	mu      sync.Mutex
	Balance int
}

func transfer(
	from *Account,
	to *Account,
	amount int,
) {
	from.mu.Lock()
	defer from.mu.Unlock()

	time.Sleep(10 * time.Millisecond)

	to.mu.Lock()
	defer to.mu.Unlock()

	from.Balance -= amount
	to.Balance += amount
}

func main() {
	a := &Account{
		ID:      1,
		Balance: 1000,
	}

	b := &Account{
		ID:      2,
		Balance: 1000,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		transfer(a, b, 100)
	}()

	go func() {
		defer wg.Done()
		transfer(b, a, 50)
	}()

	wg.Wait()

	fmt.Println(
		a.Balance,
		b.Balance,
	)
}
```

## Инвариант

Суммарный баланс должен оставаться:

```text
2000
```

Программа должна гарантированно завершаться.

---

## Требование

Исправление не должно основываться на:

```text
"давайте добавим Sleep другой длины"
```

или:

```text
"на моей машине порядок обычно такой"
```

Необходимо определить единый deterministic ordering захвата ресурсов.

---

# Task 14 ★ — Race-Free, But Still Wrong

## Тема

```text
atomic
composite invariant
linearizable snapshot
race-free != correct
```

Разработчик знает про atomics и считает, что полностью решил проблему synchronization.

```go
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

type Stats struct {
	pending atomic.Int64
	done    atomic.Int64
}

func NewStats(total int64) *Stats {
	s := &Stats{}
	s.pending.Store(total)

	return s
}

func (s *Stats) Complete() {
	s.pending.Add(-1)

	runtime.Gosched()

	s.done.Add(1)
}

func (s *Stats) Snapshot() (
	pending int64,
	done int64,
) {
	return s.pending.Load(), s.done.Load()
}

func main() {
	const total = 1000

	stats := NewStats(total)

	var wg sync.WaitGroup
	wg.Add(total)

	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()

			stats.Complete()
		}()
	}

	for {
		pending, done := stats.Snapshot()

		if pending+done != total {
			fmt.Printf(
				"BROKEN SNAPSHOT: pending=%d done=%d\n",
				pending,
				done,
			)

			break
		}

		if done == total {
			break
		}
	}

	wg.Wait()
}
```

## Интересная особенность

Программа может успешно проходить:

```bash
go run -race .
```

Но всё равно выдавать логически невозможное состояние.

---

## Инвариант

Любой `Snapshot()` должен удовлетворять:

```text
pending >= 0
done >= 0
pending + done == total
```

Переход:

```text
pending -= 1
done    += 1
```

является одной логической операцией.

---

## Задача

Исправьте модель synchronization.

Объясните:

> Почему отсутствие data race ещё не доказывает корректность concurrent algorithm?

---

# Task 15 ★ — Rate Limiter With a Timer Leak

## Тема

```text
time.Ticker
rate limiting
timer lifecycle
context
select
```

Необходимо обрабатывать не более пяти запросов в секунду.

```go
package main

import (
	"context"
	"fmt"
	"time"
)

type Request struct {
	ID int
}

func run(
	ctx context.Context,
	requests <-chan Request,
	rate int,
) {
	for request := range requests {
		<-time.Tick(
			time.Second / time.Duration(rate),
		)

		fmt.Println(
			"processed:",
			request.ID,
		)
	}
}

func main() {
	ctx := context.Background()

	requests := make(chan Request)

	go run(
		ctx,
		requests,
		5,
	)

	for i := 0; i < 20; i++ {
		requests <- Request{ID: i}
	}

	close(requests)

	time.Sleep(time.Second)
}
```

## Требуемое поведение

Один limiter должен иметь собственный lifecycle.

Необходимо поддержать:

```text
rate = N operations / second
input channel close
context cancellation
timer cleanup
```

После окончания `run` не должно оставаться принадлежащего ему ticker resource.

---

## Дополнительный вопрос

Чем отличается:

```text
RATE
```

от:

```text
CONCURRENCY
```

?

Например:

```text
workers = 5
```

не означает:

```text
5 requests / second
```

И наоборот.

---

# Legacy Review — Loop Variable Capture

Это **не основная executable-задача для современного Go**.

До Go 1.22 классическим источником ошибок был код вида:

```go
for _, value := range values {
	go func() {
		fmt.Println(value)
	}()
}
```

В старой semantics range-переменная могла переиспользоваться между итерациями, из-за чего goroutines наблюдали неожиданное значение.

Типичным исправлением было:

```go
for _, value := range values {
	go func(v int) {
		fmt.Println(v)
	}(value)
}
```

или создание iteration-local variable.

Начиная с Go 1.22 semantics loop variables изменились, поэтому не следует преподавать старый пример как гарантированно сломанный код современного Go.

Но паттерн необходимо узнавать при:

```text
code review старого Go-кода
поддержке legacy repositories
чтении старых статей
собеседованиях
```

Главная идея остаётся полезной:

> **Всегда понимайте, какое именно значение захватывает closure и каков lifecycle этой переменной.**

---

# Финальный Combat Checklist

После выполнения домашнего задания ученик должен уметь посмотреть на concurrent component и последовательно задать следующие вопросы.

## 1. Ownership

```text
Кто создал concurrent operation?
Кто владеет goroutine?
Кто владеет channel?
Кто имеет право его закрыть?
```

## 2. Lifecycle

```text
Как работа начинается?
Как завершается?
Кто ждёт завершения?
Существует ли stop condition?
```

## 3. Bounds

```text
Сколько concurrent operations может существовать?
Что ограничивает это число?
Что произойдёт под burst load?
```

## 4. Backpressure

```text
Что произойдёт,
если producer быстрее consumer?
```

Необходимо иметь явную policy:

```text
block
buffer
reject
drop
throttle
```

а не случайное поведение.

## 5. Shared State

```text
Какие данные mutable?
Какие goroutines имеют к ним доступ?
Какой invariant защищает synchronization?
```

## 6. Channels

Для каждого channel:

```text
producer?
consumer?
owner?
close condition?
```

## 7. Cancellation

```text
Что произойдёт,
если caller больше не заинтересован в результате?
```

## 8. Blocking Points

Нужно увидеть каждый потенциальный:

```text
channel send
channel receive
mutex lock
Wait
network operation
timer wait
```

и спросить:

> Как этот код выйдет отсюда при shutdown?

## 9. Errors

Concurrency не должна превращать ошибки в невидимые события.

```text
goroutine failed
```

не имеет смысла, если owner никогда не узнает об этом.

## 10. Shutdown

Правильная модель:

```text
STOP ACCEPTING
      ↓
CANCEL / DRAIN
      ↓
FINISH OWNED WORK
      ↓
WAIT
      ↓
RELEASE RESOURCES
      ↓
RETURN
```

---

# Verification Contract

Финальная проверка:

```bash
go fmt ./...
go test ./...
go test -race ./...
go vet ./...
```

Дополнительно каждая задача должна отвечать на три вопроса:

```text
1. Почему исходный код был некорректен?

2. Какую гарантию обеспечивает исправление?

3. Чем вы можете доказать,
   что исправление действительно работает?
```

---

# PASS

Домашняя работа считается пройденной, если ученик не просто знает названия:

```text
goroutine
channel
WaitGroup
Mutex
atomic
context
```

а способен связать их с инженерными задачами:

```text
goroutine → execution

channel → communication / coordination

WaitGroup → completion

Mutex → composite shared-state invariant

atomic → indivisible operation on simple state

context → cancellation / deadline signal

close(channel) → lifecycle / broadcast depending on contract

worker pool → bounded concurrency

backpressure → bounded producer credit

graceful shutdown → controlled lifecycle termination
```

Финальная mental model:

```text
                 CONCURRENT COMPONENT
                         │
       ┌─────────────────┼─────────────────┐
       │                 │                 │
   OWNERSHIP          BOUNDS            STATE
       │                 │                 │
       ▼                 ▼                 ▼
   lifecycle        worker pool      mutex / atomic
   cancellation     backpressure       invariant
       │                 │                 │
       └─────────────────┼─────────────────┘
                         ▼
                    VERIFICATION
                         │
            ┌────────────┼────────────┐
            ▼            ▼            ▼
          tests      race detector   review
```

> **Не спрашивайте сначала: «какой primitive здесь использовать?»  
> Сначала определите: «какую гарантию должна обеспечивать система?»**
> 