# Go Backend Combat — 02. CONCURRENCY

# Часть I — Programming Model

## 00:00–00:05 — Возвращаемся к JobFlow

Я хочу начать с того места, где мы закончили.

На прошлом занятии у нас появился `JobFlow`. У нас есть `Job`, есть repository, есть service. Всё работает.

Но пока наша программа устроена очень просто.

Пришла задача.

Мы её обработали.

Пришла следующая.

Мы обработали следующую.

Представим, что сегодня утром приходит не одна Job, а тысяча.

И каждая задача занимает, допустим, 100 миллисекунд.

Что произойдёт, если мы будем обрабатывать их строго последовательно?

<details style="display: inline;"><summary>Подсказка</summary>Сначала оцените стоимость одной задачи, потом умножьте её на количество задач.</details>

<details style="display: inline;"><summary>Ответ</summary>При последовательной обработке стоимость складывается. 1000 задач по 100 мс — около 100 секунд работы одного последовательного потока обработки.</details>

Да.

Теперь представьте, что эти задачи в основном независимы.

Почему мы вообще должны ждать завершения первой, прежде чем начать вторую?

Не должны.

И сегодня мы хотим прийти примерно к такой системе:

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

Но прежде чем мы начнём писать `go`, я хочу, чтобы вы запомнили одну вещь.

**Concurrency — это не "давайте запустим побольше goroutine".**

Это управление совместным выполнением, передачей данных, состоянием, временем жизни и завершением.

И сегодня мы будем каждый раз спрашивать:

> Кто владеет работой?

> Кто владеет состоянием?

> Кто имеет право остановить эту работу?

> Что ограничивает её количество?

> Как мы докажем, что всё завершилось?

Это будет наш combat-loop.

---

# 00:05–00:15 — Плакат 03A: GO CONCURRENCY MODEL

Открываем плакат `03A`.

Смотрите на него.

Я хочу, чтобы мы не читали карточки сверху вниз как учебник.

Я буду задавать вопросы.

---

## 03A / Card 1 — OWNERSHIP

Что означает ownership в concurrent-программе?

<details style="display: inline;"><summary>Подсказка</summary>Представьте, что у нас есть работа и кто-то должен отвечать за её запуск, состояние и завершение.</details>

<details style="display: inline;"><summary>Ответ</summary>У concurrent operation должен быть понятный владелец, который отвечает за её запуск, используемые ресурсы и lifecycle.</details>

И вот это очень важная идея.

Если я создаю goroutine, но потом никто не знает, зачем она существует и кто её должен остановить, у меня уже есть архитектурная проблема.

Поэтому первая привычка:

> **Не создавай goroutine без понятного владельца и смысла её существования.**

---

## 03A / Card 2 — GOROUTINE

Теперь вторая карточка.

Что делает:

```go
go process(job)
```

?

<details style="display: inline;"><summary>Подсказка</summary>Что происходит с вызовом функции и кто дальше управляет его выполнением?</details>

<details style="display: inline;"><summary>Ответ</summary>Вызов функции выполняется как goroutine — логическая единица concurrent execution, управляемая Go runtime.</details>

Хорошо.

А goroutine — это OS thread?

<details style="display: inline;"><summary>Ответ</summary>Нет. Goroutine — более лёгкая единица выполнения, управляемая runtime-ом. Go runtime распределяет goroutines по OS threads.</details>

Вот поэтому:

```text
1000 goroutines
```

не означает:

```text
1000 OS threads
```

И это мы сейчас ещё увидим на втором плакате.

---

## 03A / Card 3 — CHANNEL

Теперь вопрос.

У нас есть producer и worker.

Как передать Job от одного к другому?

<details style="display: inline;"><summary>Подсказка</summary>Нужен механизм, который одновременно позволяет передать значение между goroutines и согласовать их работу.</details>

<details style="display: inline;"><summary>Ответ</summary>Channel.</details>

Например:

```go
jobs := make(chan Job)
```

и:

```go
jobs <- job
```

а на другой стороне:

```go
job := <-jobs
```

Но channel — это не просто «очередь».

Он ещё задаёт **синхронизацию**.

Что произойдёт с unbuffered channel, если отправитель отправляет Job, а receiver пока не читает?

<details style="display: inline;"><summary>Ответ</summary>Операция отправки будет блокироваться до тех пор, пока другой участник не сможет принять значение.</details>

То есть channel уже говорит нам что-то о времени.

Отправитель и получатель должны встретиться в определённой точке выполнения.

---

## 03A / Card 4 — SELECT

Теперь у worker две возможные ситуации.

Либо появилась Job.

Либо пришла отмена.

Как worker должен ждать оба события?

<details style="display: inline;"><summary>Ответ</summary>Через `select`.</details>

Например:

```go
select {
case job := <-jobs:
	process(job)

case <-ctx.Done():
	return
}
```

Что теперь является частью worker lifecycle?

<details style="display: inline;"><summary>Ответ</summary>Worker может продолжить работу, если есть Job, либо завершиться при cancellation.</details>

---

## 03A / Card 5 — SYNCHRONIZATION

Теперь неприятный вопрос.

Что происходит, если две goroutines одновременно меняют одну переменную?

<details style="display: inline;"><summary>Подсказка</summary>Проблема не в том, что две goroutines существуют. Проблема в общем изменяемом состоянии.</details>

<details style="display: inline;"><summary>Ответ</summary>Возникает shared mutable state, которому нужны правила синхронизации. Иначе может возникнуть data race и нарушение invariant.</details>

Какой примитив мы можем использовать?

<details style="display: inline;"><summary>Ответ</summary>В зависимости от задачи — channel, mutex или atomic.</details>

---

## 03A / Card 6 — LIFECYCLE

Последний вопрос по плакату.

Как concurrent operation должна завершиться?

<details style="display: inline;"><summary>Подсказка</summary>Не «когда-нибудь». Назовите конкретные события.</details>

<details style="display: inline;"><summary>Ответ</summary>Нормальное завершение, cancellation, timeout, error или закрытие входного потока — в зависимости от контракта операции.</details>

Отлично.

Вот теперь мы начинаем писать код.

---

# 00:15–00:25 — Первая goroutine

Я открываю наш `JobFlow`.

У нас есть:

```go
func process(job Job) {
    fmt.Println("processing", job.ID)
}
```

Сначала обычный вызов:

```go
process(job)
```

Теперь:

```go
go process(job)
```

Что изменилось?

<details style="display: inline;"><summary>Ответ</summary>Функция запускается как goroutine и может выполняться concurrently с вызывающей goroutine.</details>

Хорошо.

Теперь я делаю намеренно плохой пример:

```go
func main() {
    go process(job)

    fmt.Println("main finished")
}
```

Что я ожидаю увидеть?

<details style="display: inline;"><summary>Подсказка</summary>Какая goroutine контролирует сам процесс?</details>

<details style="display: inline;"><summary>Ответ</summary>`main` завершается, процесс завершается, и запущенная goroutine может не успеть выполнить работу.</details>

Запускаем несколько раз.

Здесь может появиться очень важное наблюдение:

> **Запустить работу и дождаться её завершения — это разные задачи.**

---

# 00:25–00:35 — WaitGroup

Что нам нужно?

<details style="display: inline;"><summary>Ответ</summary>Нужно дождаться завершения goroutines.</details>

Пишем:

```go
var wg sync.WaitGroup

wg.Add(1)

go func() {
    defer wg.Done()
    process(job)
}()

wg.Wait()
```

Что делает `Add`?

<details style="display: inline;"><summary>Ответ</summary>Регистрирует количество ожидаемых завершений.</details>

А `Done`?

<details style="display: inline;"><summary>Ответ</summary>Уменьшает счётчик завершённых работ на единицу.</details>

А `Wait`?

<details style="display: inline;"><summary>Ответ</summary>Блокирует вызывающую goroutine до тех пор, пока счётчик не станет нулём.</details>

Почему мне нравится конструкция:

```go
defer wg.Done()
```

?

<details style="display: inline;"><summary>Помощь</summary>Что произойдёт, если функция завершится раньше из-за `return`?</details>

<details style="display: inline;"><summary>Ответ</summary>`Done()` всё равно будет вызван при выходе из goroutine, поэтому меньше риск забыть уменьшить счётчик.</details>

---

# 00:35–00:40 — Concurrency vs Parallelism

Теперь о двух словах, которые постоянно путают.

Что такое concurrency?

<details style="display: inline;"><summary>Ответ</summary>Структура программы, в которой несколько независимых задач могут находиться в процессе выполнения и прогрессировать независимо.</details>

А parallelism?

<details style="display: inline;"><summary>Ответ</summary>Фактическое одновременное выполнение нескольких частей работы на разных CPU execution resources.</details>

Очень простая формула:

> **Concurrency — как организована работа. Parallelism — сколько работы реально выполняется одновременно.**

Это будет особенно важно, когда мы посмотрим 03B.

---

# 00:40–00:50 — Один worker на каждую Job

Теперь я даю ученику задачу.

Есть:

```go
jobs []Job
```

Нужно обработать всё concurrently.

Как бы вы сделали?

<details style="display: inline;"><summary>Подсказка</summary>Для каждой Job можно запустить работу отдельно. Как дождаться всех?</details>

<details style="display: inline;"><summary>Ответ</summary>Создавать goroutine на каждую Job и использовать `WaitGroup` для ожидания.</details>

Получаем примерно:

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

А теперь мой следующий вопрос намного важнее.

Сколько goroutines появится, если у нас миллион Job?

<details style="display: inline;"><summary>Ответ</summary>Потенциально миллион goroutines.</details>

И это уже плохая идея.

Не потому, что goroutine «дорогая» сама по себе.

А потому, что **входной объём работы начал напрямую определять объём concurrent resources**.

Вот теперь открываем второй плакат.

---

# 00:50–01:10 — Плакат 03B: GO RUNTIME

Смотрите на плакат `03B`.

Теперь мы можем ответить на вопрос:

> Что реально происходит после `go process(job)`?

---

## 03B / Card 1 — CREATE

После:

```go
go process(job)
```

что создаётся?

<details style="display: inline;"><summary>Ответ</summary>Создаётся goroutine G, которая становится готовой к выполнению.</details>

То есть:

```text
go f()
 ↓
G
```

Но G ещё не означает «сейчас выполняется на CPU».

Правильно?

<details style="display: inline;"><summary>Ответ</summary>Да. Она становится runnable и попадает в область работы scheduler.</details>

---

## 03B / Card 2 — RUNNABLE

Где находится готовая работа?

<details style="display: inline;"><summary>Ответ</summary>В scheduler run queues — локальных очередях, связанных с P, и глобальной очереди runnable work.</details>

Смотрите на карточку.

Мы видим:

```text
P0 [G][G][G]
P1 [G][G]
P2 [G][G][G][G]
```

Это не означает, что четыре G у P2 сейчас одновременно выполняются.

Это означает:

> **у P есть runnable work, которую scheduler может выдавать на выполнение.**

---

## 03B / Card 3 — SCHEDULE

Теперь самый важный вопрос.

У нас есть runnable G.

Что должно произойти, чтобы она реально начала выполняться?

<details style="display: inline;"><summary>Подсказка</summary>Посмотрите на три буквы в карточке: G / P / M.</details>

<details style="display: inline;"><summary>Ответ</summary>Scheduler связывает runnable G с P и M, после чего Go-код может выполняться на соответствующем OS thread.</details>

Что такое G?

<details style="display: inline;"><summary>Ответ</summary>Работа — goroutine.</details>

Что такое P?

<details style="display: inline;"><summary>Ответ</summary>Ресурс runtime scheduler, необходимый для выполнения Go-кода.</details>

Что такое M?

<details style="display: inline;"><summary>Ответ</summary>OS thread, на котором фактически выполняется Go-код.</details>

И теперь важный вывод:

> **G, P и M — не три последовательных шага. Scheduler связывает их для выполнения.**

Это очень важно не перепутать.

---

# 01:10–01:20 — 03B / EXECUTE

Теперь G действительно выполняется.

Что может произойти дальше?

<details style="display: inline;"><summary>Подсказка</summary>На карточке четыре пути. Что может сделать работа, когда CPU ей больше не нужен?</details>

<details style="display: inline;"><summary>Ответ</summary>Она может продолжать вычисление, быть вытеснена scheduler, заблокироваться на syscall или ждать network I/O.</details>

Почему preemption важен?

<details style="display: inline;"><summary>Ответ</summary>Он позволяет scheduler не давать одной goroutine бесконтрольно удерживать execution resource и поддерживать прогресс других runnable goroutines.</details>

---

## 03B / WAIT

Теперь network request.

Представим, что worker делает:

```text
read from socket
```

Но данные пока не пришли.

Что делать goroutine?

<details style="display: inline;"><summary>Ответ</summary>Она должна ждать I/O, а runtime может использовать netpoller, чтобы не удерживать CPU просто ради ожидания сети.</details>

То есть:

```text
G
 ↓
NETWORK WAIT
 ↓
NETPOLL
```

И когда данные готовы?

---

## 03B / WAKE

Что происходит?

<details style="display: inline;"><summary>Ответ</summary>G снова становится runnable и возвращается в scheduler flow.</details>

Вот теперь мы видим настоящий цикл:

```text
G
 ↓
RUNQ
 ↓
SCHEDULE
 ↓
EXECUTE
 ↓
WAIT
 ↓
WAKE
 ↓
RUNQ
```

И это гораздо более полезная модель, чем просто запомнить:

> G = goroutine, P = processor, M = machine.

---

# Часть II — Engineering Model

# 01:20–01:30 — Плакат 03C: CONCURRENCY ENGINEERING

Теперь открываем третий плакат.

И говорим:

> «Мы знаем, как писать concurrency. Мы знаем, как runtime её исполняет. Теперь главный backend-вопрос: как сделать так, чтобы система не развалилась под нагрузкой?»

Смотрим на первый блок.

---

## 03C / Card 1 — LOAD

Сколько работы может прийти в систему?

<details style="display: inline;"><summary>Подсказка</summary>Одного количества запросов недостаточно. Что ещё важно?</details>

<details style="display: inline;"><summary>Ответ</summary>Rate, burst и стоимость одной операции.</details>

Да.

Например:

```text
1000 jobs/sec
```

при дешёвой операции и:

```text
1000 jobs/sec
```

при операции, занимающей 500 ms CPU и 100 MB памяти —

это две совершенно разные системы.

Поэтому concurrency начинается не с:

> «Сколько goroutines поставить?»

А с:

> **Какова наша нагрузка и стоимость работы?**

---

# 01:30–01:40 — BOUND

Допустим:

```text
1 000 000 jobs
```

и:

```text
32 workers
```

Почему я вообще хочу число 32 ограничить?

<details style="display: inline;"><summary>Ответ</summary>Чтобы размер входной нагрузки не превращался автоматически в неограниченное потребление CPU, памяти, соединений и других ресурсов.</details>

Вот наш основной принцип:

> **Concurrency должна быть bounded.**

И мы ограничиваем не только workers.

Что ещё можно ограничивать?

<details style="display: inline;"><summary>Ответ</summary>Количество in-flight operations, queue capacity, memory, database connections, downstream requests и другие ресурсы.</details>

И здесь я бы прямо спросил:

> «Что произойдёт, если producer продолжает создавать работу быстрее, чем workers её обрабатывают?»

<details style="display: inline;"><summary>Ответ</summary>Начнёт расти очередь или количество накопленной работы. Без ограничения это может привести к росту памяти, latency и в итоге к отказу системы.</details>

---

# 01:40–01:50 — Backpressure

Вот теперь мы пришли к очень важному слову.

**Backpressure.**

Что это означает?

<details style="display: inline;"><summary>Подсказка</summary>Представьте, что consumer обрабатывает 100 jobs/sec, а producer присылает 1000.</details>

<details style="display: inline;"><summary>Ответ</summary>Система должна каким-то образом ограничить или замедлить поступление работы: блокировать producer, ограничивать очередь, отклонять работу, снижать нагрузку или применять другой policy.</details>

И принцип:

> **Если consumer не успевает, producer не должен иметь бесконечный кредит памяти.**

Это одна из самых важных production-идей concurrency.

---

# 01:50–02:00 — DISTRIBUTE / Worker Pool

Теперь собираем нормальный worker pool.

Я показываю:

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

И спрашиваю:

> «Кто владеет queue?»

<details style="display: inline;"><summary>Ответ</summary>Producer владеет отправкой, workers потребляют. Конкретная ownership policy должна быть определена явно.</details>

Теперь:

> «Кто закрывает channel?»

<details style="display: inline;"><summary>Ответ</summary>Сторона, которая владеет отправкой и знает, что новых значений больше не будет.</details>

И теперь пусть ученик сам реализует.

```go
type WorkerPool struct {
    workers int
}
```

---

# 02:00–02:10 — Worker Pool: первая реализация

Ученик пишет worker:

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

Если ученик забывает `Done`, спрашиваю:

> «Кто уведомит owner, что worker закончил?»

Если забывает `close(jobs)`:

> «Как worker узнает, что новых задач больше не будет?»

Если закрывает channel из worker:

> «Кто владеет отправкой?»

Это всё не syntax questions.

Это ownership questions.

---

# 02:10–02:20 — BUFFERED vs UNBUFFERED

Теперь:

```go
jobs := make(chan Job)
```

Что означает unbuffered channel?

<details style="display: inline;"><summary>Ответ</summary>Отправитель и receiver синхронизируются непосредственно: send обычно ждёт соответствующего receive.</details>

А:

```go
jobs := make(chan Job, 100)
```

?

<details style="display: inline;"><summary>Ответ</summary>Channel имеет buffer capacity 100. Producer может временно опередить consumer, пока buffer не заполнен.</details>

Очень важно:

> **Buffer не устраняет backpressure. Он только переносит его во времени до заполнения capacity.**

Что произойдёт, когда buffer заполнится?

<details style="display: inline;"><summary>Ответ</summary>Следующая отправка должна будет ждать или, в зависимости от выбранной архитектуры, система должна применить другую policy: reject, drop, throttle и т.п.</details>

---

# 02:20–02:30 — Select + Context + Timeout

Теперь нам надо научить worker не просто ждать Job, а уметь остановиться.

Показываю:

```go
select {
case job := <-jobs:
    process(job)

case <-ctx.Done():
    return
}
```

Что контролирует `ctx.Done()`?

<details style="display: inline;"><summary>Ответ</summary>Сигнал cancellation, который сообщает, что текущая работа больше не должна продолжаться.</details>

Теперь:

```go
ctx, cancel := context.WithTimeout(
    context.Background(),
    time.Second,
)
defer cancel()
```

Что произойдёт через секунду?

<details style="display: inline;"><summary>Ответ</summary>Context будет отменён по deadline, `ctx.Done()` станет готовым к чтению, и worker сможет завершить или прервать соответствующую операцию.</details>

И главное:

> **Cancellation — это сигнал, а не принудительное убийство goroutine.**

Что должен сделать код?

<details style="inline;"><summary>Ответ</summary>Наблюдать context и корректно завершать работу при получении cancellation.</details>

---

# 02:30–02:40 — BREAK THE CODE: DATA RACE

Теперь я намеренно ломаю систему.

Пишем:

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

Мы ждём завершения всех goroutines.

Что ожидаем?

<details style="display: inline;"><summary>Ответ</summary>1000.</details>

А получим ли мы гарантированно 1000?

<details style="display: inline;"><summary>Ответ</summary>Нет. `processed++` — составная операция read-modify-write. Несогласованный concurrent access может привести к data race и потерянным обновлениям.</details>

Теперь запускаем:

```bash
go test -race ./...
```

И читаем diagnostic.

Я здесь ничего не объясняю заранее.

Я спрашиваю:

> «Что runtime/tooling нам сейчас доказал?»

<details style="display: inline;"><summary>Ответ</summary>Что несколько goroutines имеют конфликтующий concurrent access к общей памяти без необходимой синхронизации.</details>

---

# 02:40–02:50 — Mutex

Теперь вопрос.

Какой invariant мы хотим сохранить?

<details style="display: inline;"><summary>Ответ</summary>Каждое увеличение `processed` должно происходить как согласованная операция, а итоговое состояние не должно терять обновления.</details>

Добавляем:

```go
type Stats struct {
    mu        sync.Mutex
    processed int
}
```

И:

```go
stats.mu.Lock()
stats.processed++
stats.mu.Unlock()
```

Почему можно написать:

```go
stats.mu.Lock()
defer stats.mu.Unlock()
```

<details style="display: inline;"><summary>Ответ</summary>Чтобы гарантировать освобождение mutex при выходе из текущей функции, включая ранний return.</details>

И вот теперь я задаю более senior-level вопрос:

> «Что на самом деле защищает mutex? Переменную или invariant?»

<details style="display: inline;"><summary>Ответ</summary>Mutex защищает согласованное изменение некоторого shared state и invariant, который должен сохраняться во время этих операций. Переменная — лишь часть состояния.</details>

Вот это важно.

Не:

> «У меня есть mutex, значит всё безопасно».

А:

> **«Вот какое состояние я защищаю, вот какие операции должны быть атомарными относительно друг друга».**

---

# 02:50–03:00 — Atomic

Теперь тот же простой counter.

Нужно ли здесь обязательно использовать mutex?

<details style="display: inline;"><summary>Ответ</summary>Нет. Для отдельного числового counter можно использовать atomic.</details>

Например:

```go
type Stats struct {
    processed atomic.Int64
}
```

И:

```go
stats.processed.Add(1)
```

Что здесь важно понять?

<details style="display: inline;"><summary>Ответ</summary>Atomic — это примитив для атомарных операций над отдельным значением, а не универсальная замена mutex.</details>

Даю пример:

```go
processed
failed
total
```

И спрашиваю:

> «Допустим, нам нужно сохранять invariant `processed + failed == total`. Что естественнее — три независимых atomic или mutex вокруг составного состояния?»

<details style="display: inline;"><summary>Ответ</summary>Mutex естественнее, потому что invariant относится к нескольким связанным значениям, которые нужно менять согласованно.</details>

Это один из ключевых переходов от «знаю primitives» к «умею проектировать state».

---

# 03:00–03:10 — Concurrent Memory Repository

Теперь возвращаемся к JobFlow.

Наш repository:

```go
type MemoryJobRepository struct {
    jobs map[string]Job
}
```

А теперь у нас несколько workers.

Что произойдёт, если два workers одновременно вызовут `Add` или один читает, а другой пишет?

<details style="display: inline;"><summary>Ответ</summary>Обычная map не предназначена для concurrent writes без соответствующей синхронизации. Можно получить race и некорректное выполнение.</details>

Что добавляем?

<details style="display: inline;"><summary>Ответ</summary>Mutex или RWMutex, если модель доступа это оправдывает.</details>

Например:

```go
type MemoryJobRepository struct {
    mu   sync.RWMutex
    jobs map[string]Job
}
```

Почему `RWMutex`?

<details style="display: inline;"><summary>Ответ</summary>Он позволяет нескольким readers работать одновременно, сохраняя эксклюзивность writer. Но использовать его стоит только если такая модель действительно подходит workload.</details>

И я специально добавляю:

> «Не оптимизируйте mutex только потому, что существует `RWMutex`. Сначала нужна модель доступа, потом выбор примитива.»

---

# 03:10–03:20 — Goroutine Leak

Теперь ещё один intentional failure.

Показываю:

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

Что здесь не так?

<details style="display: inline;"><summary>Подсказка</summary>Представьте, что producer больше никогда не отправит Job. Что остановит worker?</details>

<details style="display: inline;"><summary>Ответ</summary>Ничто. Goroutine может остаться навсегда ожидающей channel. Это потенциальный goroutine leak.</details>

Как исправить?

<details style="display: inline;"><summary>Ответ</summary>Нужен явный stop condition: закрытие channel, context cancellation или другой определённый lifecycle mechanism.</details>

Пишем:

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

И я возвращаю ученика к первому плакату:

> **OWNERSHIP.**

> Кто владеет worker?

> Кто его запускает?

> Кто его останавливает?

> Какое событие означает `DONE`?

---

# 03:20–03:30 — Graceful Shutdown

Теперь представим production.

Процесс получает:

```text
SIGTERM
```

Что мы хотим сделать?

<details style="display: inline;"><summary>Подсказка</summary>Мы не хотим просто убить процесс. Что нужно сделать с новой работой и текущими workers?</details>

<details style="display: inline;"><summary>Ответ</summary>Перестать принимать новую работу, инициировать cancellation, дать допустимой работе завершиться по policy, дождаться выхода workers и затем завершить процесс.</details>

Получаем:

```text
STOP ACCEPTING
      ↓
CANCEL
      ↓
FINISH ALLOWED WORK
      ↓
WAIT
      ↓
EXIT
```

И я спрашиваю:

> «Почему просто `os.Exit()` недостаточно?»

<details style="display: inline;"><summary>Ответ</summary>Потому что можно потерять незавершённую работу, не закрыть ресурсы и не выполнить необходимую cleanup/recovery policy.</details>

---

# 03:30–03:40 — Большой практикум

Теперь я прекращаю объяснять.

Говорю ученику:

> «Сейчас я дам вам требования. Вы проектируете worker pool сами. Я не буду давать структуру кода.»

Требования:

```text
workers = configurable

all jobs are processed

maximum concurrent workers is bounded

context cancellation stops new processing

Run waits for workers to exit

processing errors are observable

repository is safe for concurrent access

no worker goroutine remains after Run
```

Спрашиваю:

> «С чего начнёте?»

<details style="display: inline;"><summary>Подсказка 1</summary>Сначала назовите ресурсы и ownership.</details>

<details style="display: inline;"><summary>Подсказка 2</summary>Кто создаёт workers? Кто владеет queue? Кто закрывает queue?</details>

<details style="display: inline;"><summary>Подсказка 3</summary>Как worker узнает про Job и про cancellation одновременно?</details>

<details style="display: inline;"><summary>Ответ</summary>Producer владеет отправкой в jobs channel. Создаётся фиксированное количество workers. Каждый worker получает jobs через receive-only channel и завершает работу по `ctx.Done()` либо после закрытия channel. `WaitGroup` ждёт завершения всех workers. Shared state repository защищается synchronization primitive. Количество workers задаёт верхнюю границу concurrent work.</details>

---

# 03:40–03:50 — Атакуем решение

Теперь я намеренно меняю условия.

### Attack 1

> «Очередь получила в десять раз больше задач. Что произойдёт?»

<details style="display: inline;"><summary>Ответ</summary>При bounded queue начнёт проявляться backpressure. Нужно явно выбрать policy: блокировать producer, ограничивать rate, отклонять часть работы или применять другой механизм.</details>

### Attack 2

> «Worker выполняет Job две минуты.»

<details style="display: inline;"><summary>Ответ</summary>Должен существовать timeout/deadline/cancellation policy, иначе in-flight work может неограниченно удерживать ресурсы.</details>

### Attack 3

> «Клиент больше не ждёт результат.»

<details><summary>Ответ</summary>Если работа больше не нужна, cancellation должна распространяться по call chain, чтобы освобождать ресурсы.</details>

### Attack 4

> «Один worker упал.»

<details><summary>Ответ</summary>Нужно определить failure policy: recovery, retry, error propagation, replacement worker или иной механизм в зависимости от архитектуры.</details>

### Attack 5

> «Приходит SIGTERM.»

<details><summary>Ответ</summary>System-wide graceful shutdown: stop accepting, cancel according to policy, finish allowed work, wait, exit.</details>

---

# 03:50–04:00 — Финальная проверка

Теперь я выключаю код.

Смотрим на три плаката.

Первый вопрос.

Что такое concurrency?

<details><summary>Ответ</summary>Организация независимого выполнения и координации между несколькими задачами.</details>

Что делает goroutine?

<details><summary>Ответ</summary>Представляет независимую единицу concurrent execution, которой управляет Go runtime.</details>

Зачем channel?

<details><summary>Ответ</summary>Для передачи значений и координации между goroutines.</details>

Зачем `select`?

<details><summary>Ответ</summary>Для ожидания нескольких channel operations или событий.</details>

Зачем `WaitGroup`?

<details><summary>Ответ</summary>Чтобы дождаться завершения зарегистрированной группы goroutines.</details>

Когда mutex?

<details><summary>Ответ</summary>Когда нужно защищать shared mutable state и связанный с ним invariant.</details>

Когда atomic?

<details><summary>Ответ</summary>Для атомарных операций над простым независимым shared state.</details>

Что такое backpressure?

<details><summary>Ответ</summary>Механизм, который не позволяет producer бесконтрольно опережать consumer и заставляет систему применять ограничение нагрузки.</details>

Что такое goroutine leak?

<details><summary>Ответ</summary>Goroutine, которая продолжает жить после потери необходимости в её работе или не имеет достижимого stop condition.</details>

Что делает context?

<details><summary>Ответ</summary>Передаёт cancellation и deadline через цепочку вызовов.</details>

И последний вопрос:

> «Какие четыре вещи ты теперь проверишь, увидев любой concurrent component?»

<details><summary>Подсказка</summary>Посмотрите на три плаката. Особенно 03A и 03C.</details>

<details><summary>Ответ</summary>Ownership, lifecycle/stop condition, synchronization/resource bounds и failure/cancellation behavior.</details>

---

# 04:00 — Закрытие урока

Я бы завершил занятие не списком терминов, а вот так:

> «Сегодня мы начали с очень простой идеи: у нас слишком много Job, и мы хотим обрабатывать их одновременно.
>
> Но посмотрите, до чего мы дошли.
>
> Нам пришлось определить, кто владеет работой.
>
> Нам пришлось отделить concurrency от parallelism.
>
> Нам понадобились channels, потому что concurrent работа должна как-то взаимодействовать.
>
> Нам понадобился worker pool, потому что бесконечное число goroutines — это не управление ресурсами.
>
> Нам понадобились mutex и atomic, потому что shared state должен иметь определённую семантику.
>
> Нам понадобился context, потому что ненужную работу надо уметь отменять.
>
> Нам понадобился graceful shutdown, потому что программа должна не только запускаться, но и корректно заканчиваться.
>
> И наконец race detector, потому что наш intuition о concurrent code недостаточен. Нам нужны evidence и verification.
>
> Вот поэтому concurrency — это не знание команды `go`.
>
> **Concurrency — это управление работой, состоянием, ресурсами и временем жизни.**»

Дальше:

> «И теперь у нас есть хороший worker pool.
>
> Но он пока живёт внутри одной программы.
>
> В следующем уроке я задам другой вопрос:
>
> **А что, если HTTP, business logic и storage начинают жить по разные стороны границ?**
>
> Тогда нам придётся научиться управлять уже не только concurrent execution, но и зависимостями между частями системы.»

---

# Домашнее задание после 02.CONCURRENCY

## 1. Worker Pool v2

Добавить:

```go
pool := NewWorkerPool(4)
```

и:

```go
Run(ctx, jobs, process)
```

Параметры workers должны быть конфигурируемыми.

## 2. Concurrent-safe repository

Переделать `MemoryJobRepository` так, чтобы он был безопасен для concurrent access.

Требование:

```text
concurrent Add
concurrent Get
concurrent Delete
concurrent List
```

не должны создавать data race.

## 3. Statistics

Добавить:

```text
processed
failed
```

Для независимых counters использовать atomic.

## 4. Cancellation tests

Обязательно проверить:

```text
cancel before processing
cancel while waiting
cancel during processing
```

## 5. Shutdown test

Проверить:

> после `Run()` не остаётся worker goroutines, относящихся к pool.

## 6. Verification

Запустить:

```bash
go test ./...
go test -race ./...
go vet ./...
```

PASS означает:

```text
worker count bounded
race-free
cancellation works
shutdown completes
repository is concurrency-safe
tests pass
```

---

# Mentor note

На этом уроке **не надо пытаться за 120 минут досконально объяснить runtime internals**.

03B используется как короткая mental model:

```text
go f()
→ G
→ RUNQ
→ scheduler
→ G + P + M
→ execute
→ wait
→ wake
→ RUNQ
```

Главный coding outcome урока — не знание GMP.

Он должен выйти с пониманием:

> **«Я могу написать concurrent worker system, определить её ownership, ограничить concurrency, обнаружить race, отменить работу и корректно завершить систему».**

И это уже реальный переход от знания Go к **Go backend engineering**.
